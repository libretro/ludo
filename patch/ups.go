package patch

import (
	"errors"
	"hash"
	"hash/crc32"
)

type file struct {
	Data     []byte
	Offset   int
	Checksum uint32
	Hash     hash.Hash32
}

type upsData struct {
	Patch  *file // ups patch
	Source *file // game to be patched
	Target *file // result of the patching process
}

func upsPatchRead(data *upsData) (n byte) {
	if data != nil && data.Patch.Offset < len(data.Patch.Data) {
		n = data.Patch.Data[data.Patch.Offset]
		data.Patch.Offset++
		data.Patch.Hash.Write([]byte{n})
		data.Patch.Checksum = ^data.Patch.Hash.Sum32()
		return
	}
	return
}

func upsSourceRead(data *upsData) (n byte) {
	if data != nil && data.Source.Offset < len(data.Source.Data) {
		n = data.Source.Data[data.Source.Offset]
		data.Source.Offset++
		data.Source.Hash.Write([]byte{n})
		data.Source.Checksum = ^data.Source.Hash.Sum32()
		return
	}
	return
}

func upsTargetWrite(data *upsData, n byte) {
	if data != nil && data.Target.Offset < len(data.Target.Data) {
		data.Target.Data[data.Target.Offset] = n
		data.Target.Hash.Write([]byte{n})
		data.Target.Checksum = ^data.Target.Hash.Sum32()
	}

	if data != nil {
		data.Target.Offset++
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
		Patch: &file{
			Data: patchData,
			Hash: crc32.NewIEEE(),
		},
		Source: &file{
			Data: sourceData,
			Hash: crc32.NewIEEE(),
		},
		Target: &file{
			Hash: crc32.NewIEEE(),
		},
	}

	if len(data.Patch.Data) < 18 {
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

	if len(data.Source.Data) != sourceReadLength &&
		len(data.Source.Data) != targetReadLength {
		return nil, errors.New("invalid source")
	}

	targetLength := sourceReadLength
	if len(data.Source.Data) == sourceReadLength {
		targetLength = targetReadLength
	}

	prov := make([]byte, targetLength)
	data.Target.Data = prov

	for data.Patch.Offset < len(data.Patch.Data)-12 {
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

	for data.Source.Offset < len(data.Source.Data) {
		upsTargetWrite(&data, upsSourceRead(&data))
	}
	for data.Target.Offset < len(data.Target.Data) {
		upsTargetWrite(&data, upsSourceRead(&data))
	}

	if err := checks(&data, sourceReadLength, targetReadLength); err != nil {
		return nil, err
	}
	return &data.Target.Data, nil
}

// checks verifies that the patching process went well by comparing checksums
func checks(data *upsData, sourceReadLength, targetReadLength int) error {
	var sourceReadChecksum uint32
	for i := 0; i < 4; i++ {
		sourceReadChecksum |= uint32(upsPatchRead(data)) << (i * 8)
	}
	var targetReadChecksum uint32
	for i := 0; i < 4; i++ {
		targetReadChecksum |= uint32(upsPatchRead(data)) << (i * 8)
	}

	patchResultChecksum := ^data.Patch.Checksum
	data.Source.Checksum = ^data.Source.Checksum
	data.Target.Checksum = ^data.Target.Checksum

	var patchReadChecksum uint32
	for i := 0; i < 4; i++ {
		patchReadChecksum |= uint32(upsPatchRead(data)) << (i * 8)
	}

	if patchResultChecksum != patchReadChecksum {
		return errors.New("invalid patch")
	}

	if data.Source.Checksum == sourceReadChecksum && len(data.Source.Data) == sourceReadLength {
		if data.Target.Checksum == targetReadChecksum && len(data.Target.Data) == targetReadLength {
			return nil
		}
		return errors.New("invalid target")
	} else if data.Source.Checksum == targetReadChecksum && len(data.Source.Data) == targetReadLength {
		if data.Target.Checksum == sourceReadChecksum && len(data.Target.Data) == sourceReadLength {
			return nil
		}
		return errors.New("invalid target")
	}

	return errors.New("invalid source")
}
