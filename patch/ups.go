package patch

import (
	"errors"
	"hash"
	"hash/crc32"
)

type upsData struct {
	PatchData      []byte
	SourceData     []byte
	TargetData     []byte
	PatchOffset    int
	SourceOffset   int
	TargetOffset   int
	PatchChecksum  uint32
	SourceChecksum uint32
	TargetChecksum uint32
	PatchHash      hash.Hash32
	SourceHash     hash.Hash32
	TargetHash     hash.Hash32
}

func upsPatchRead(data *upsData) (n byte) {
	if data != nil && data.PatchOffset < len(data.PatchData) {
		n = data.PatchData[data.PatchOffset]
		data.PatchOffset++
		data.PatchHash.Write([]byte{n})
		data.PatchChecksum = ^data.PatchHash.Sum32()
		return
	}
	return
}

func upsSourceRead(data *upsData) (n byte) {
	if data != nil && data.SourceOffset < len(data.SourceData) {
		n = data.SourceData[data.SourceOffset]
		data.SourceOffset++
		data.SourceHash.Write([]byte{n})
		data.SourceChecksum = ^data.SourceHash.Sum32()
		return
	}
	return
}

func upsTargetWrite(data *upsData, n byte) {
	if data != nil && data.TargetOffset < len(data.TargetData) {
		data.TargetData[data.TargetOffset] = n
		data.TargetHash.Write([]byte{n})
		data.TargetChecksum = ^data.TargetHash.Sum32()
	}

	if data != nil {
		data.TargetOffset++
	}
}

func upsDecode(data *upsData) int {
	var offset = 0
	var shift = 1
	for {
		x := upsPatchRead(data)
		offset += int(x&0x7f) * shift

		if x&0x80 != 0 {
			break
		}
		shift <<= 7
		offset += shift
	}
	return offset
}

func applyUPS(patchData, sourceData []byte) (*[]byte, error) {
	data := upsData{
		PatchData:  patchData,
		SourceData: sourceData,
		PatchHash:  crc32.NewIEEE(),
		SourceHash: crc32.NewIEEE(),
		TargetHash: crc32.NewIEEE(),
	}

	if len(data.PatchData) < 18 {
		return nil, errors.New("patch too small")
	}

	if upsPatchRead(&data) != 'U' ||
		upsPatchRead(&data) != 'P' ||
		upsPatchRead(&data) != 'S' ||
		upsPatchRead(&data) != '1' {
		return nil, errors.New("invalid patch header")
	}

	sourceReadLength := upsDecode(&data)
	targetReadLength := upsDecode(&data)

	if len(data.SourceData) != sourceReadLength &&
		len(data.SourceData) != targetReadLength {
		return nil, errors.New("invalid source")
	}

	targetLength := sourceReadLength
	if len(data.SourceData) == sourceReadLength {
		targetLength = targetReadLength
	}

	if len(data.TargetData) < targetLength {
		prov := make([]byte, targetLength)
		data.TargetData = prov
	}

	for data.PatchOffset < len(data.PatchData)-12 {
		for length := upsDecode(&data); length > 0; length-- {
			upsTargetWrite(&data, upsSourceRead(&data))
		}
		for {
			patchXOR := upsPatchRead(&data)
			upsTargetWrite(&data, patchXOR^upsSourceRead(&data))
			if patchXOR == 0 {
				break
			}
		}
	}

	for data.SourceOffset < len(data.SourceData) {
		upsTargetWrite(&data, upsSourceRead(&data))
	}
	for data.TargetOffset < len(data.TargetData) {
		upsTargetWrite(&data, upsSourceRead(&data))
	}

	var sourceReadChecksum uint32
	for i := 0; i < 4; i++ {
		sourceReadChecksum |= uint32(upsPatchRead(&data)) << (i * 8)
	}
	var targetReadChecksum uint32
	for i := 0; i < 4; i++ {
		targetReadChecksum |= uint32(upsPatchRead(&data)) << (i * 8)
	}

	patchResultChecksum := ^data.PatchChecksum
	data.SourceChecksum = ^data.SourceChecksum
	data.TargetChecksum = ^data.TargetChecksum

	var patchReadChecksum uint32
	for i := 0; i < 4; i++ {
		patchReadChecksum |= uint32(upsPatchRead(&data)) << (i * 8)
	}

	if patchResultChecksum != patchReadChecksum {
		return nil, errors.New("invalid patch")
	}

	if data.SourceChecksum == sourceReadChecksum &&
		len(data.SourceData) == sourceReadLength {
		if data.TargetChecksum == targetReadChecksum &&
			len(data.TargetData) == targetReadLength {
			return &data.TargetData, nil
		}
		return nil, errors.New("invalid target")
	} else if data.SourceChecksum == targetReadChecksum &&
		len(data.SourceData) == targetReadLength {
		if data.TargetChecksum == sourceReadChecksum &&
			len(data.TargetData) == sourceReadLength {
			return &data.TargetData, nil
		}
		return nil, errors.New("invalid target")
	}

	return nil, errors.New("invalid source")
}
