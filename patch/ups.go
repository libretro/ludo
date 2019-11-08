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
	PatchLength    uint64
	SourceLength   uint64
	TargetLength   uint64
	PatchOffset    uint64
	SourceOffset   uint64
	TargetOffset   uint64
	PatchChecksum  uint32
	SourceChecksum uint32
	TargetChecksum uint32
	PatchHash      hash.Hash32
	SourceHash     hash.Hash32
	TargetHash     hash.Hash32
}

func upsPatchRead(data *upsData) (n byte) {
	if data != nil && data.PatchOffset < data.PatchLength {
		n = data.PatchData[data.PatchOffset]
		data.PatchOffset++
		data.PatchHash.Write([]byte{n})
		data.PatchChecksum = ^data.PatchHash.Sum32()
		return
	}
	return
}

func upsSourceRead(data *upsData) (n byte) {
	if data != nil && data.SourceOffset < data.SourceLength {
		n = data.SourceData[data.SourceOffset]
		data.SourceOffset++
		data.SourceHash.Write([]byte{n})
		data.SourceChecksum = ^data.SourceHash.Sum32()
		return
	}
	return
}

func upsTargetWrite(data *upsData, n byte) {
	if data != nil && data.TargetOffset < data.TargetLength {
		data.TargetData[data.TargetOffset] = n
		data.TargetHash.Write([]byte{n})
		data.TargetChecksum = ^data.TargetHash.Sum32()
	}

	if data != nil {
		data.TargetOffset++
	}
}

func upsDecode(data *upsData) uint64 {
	var offset uint64 = 0
	var shift uint64 = 1
	for {
		x := upsPatchRead(data)
		offset += uint64(x&0x7f) * shift

		if x&0x80 != 0 {
			break
		}
		shift <<= 7
		offset += shift
	}
	return offset
}

// UPSApplyPatch applies the UPS patch on the target data
func UPSApplyPatch(
	patchData []byte, patchLength uint64,
	sourceData []byte, sourceLength uint64,
	targetData *[]byte, targetLength *uint64) error {

	var patchReadChecksum uint32
	var sourceReadChecksum uint32
	var targetReadChecksum uint32

	data := upsData{
		PatchData:    patchData,
		SourceData:   sourceData,
		TargetData:   *targetData,
		PatchLength:  patchLength,
		SourceLength: sourceLength,
		TargetLength: *targetLength,
		PatchHash:    crc32.NewIEEE(),
		SourceHash:   crc32.NewIEEE(),
		TargetHash:   crc32.NewIEEE(),
	}

	if data.PatchLength < 18 {
		return errors.New("patch too small")
	}

	if upsPatchRead(&data) != 'U' ||
		upsPatchRead(&data) != 'P' ||
		upsPatchRead(&data) != 'S' ||
		upsPatchRead(&data) != '1' {
		return errors.New("invalid patch header")
	}

	sourceReadLength := upsDecode(&data)
	targetReadLength := upsDecode(&data)

	if data.SourceLength != sourceReadLength &&
		data.SourceLength != targetReadLength {
		return errors.New("invalid source")
	}

	*targetLength = sourceReadLength
	if data.SourceLength == sourceReadLength {
		*targetLength = targetReadLength
	}

	if data.TargetLength < *targetLength {
		prov := make([]byte, *targetLength)
		*targetData = prov
		data.TargetData = prov
	}

	data.TargetLength = *targetLength

	for data.PatchOffset < data.PatchLength-12 {
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

	for data.SourceOffset < data.SourceLength {
		upsTargetWrite(&data, upsSourceRead(&data))
	}
	for data.TargetOffset < data.TargetLength {
		upsTargetWrite(&data, upsSourceRead(&data))
	}

	for i := 0; i < 4; i++ {
		sourceReadChecksum |= uint32(upsPatchRead(&data)) << (i * 8)

	}
	for i := 0; i < 4; i++ {
		targetReadChecksum |= uint32(upsPatchRead(&data)) << (i * 8)
	}

	patchResultChecksum := ^data.PatchChecksum
	data.SourceChecksum = ^data.SourceChecksum
	data.TargetChecksum = ^data.TargetChecksum

	for i := 0; i < 4; i++ {
		patchReadChecksum |= uint32(upsPatchRead(&data)) << (i * 8)
	}

	if patchResultChecksum != patchReadChecksum {
		return errors.New("invalid patch")
	}

	if data.SourceChecksum == sourceReadChecksum &&
		data.SourceLength == sourceReadLength {
		if data.TargetChecksum == targetReadChecksum &&
			data.TargetLength == targetReadLength {
			return nil
		}
		return errors.New("invalid target")
	} else if data.SourceChecksum == targetReadChecksum &&
		data.SourceLength == targetReadLength {
		if data.TargetChecksum == sourceReadChecksum &&
			data.TargetLength == sourceReadLength {
			return nil
		}
		return errors.New("invalid target")
	}

	return errors.New("invalid source")
}
