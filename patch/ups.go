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

func upsRead(f *file) (n byte) {
	if f.Offset < len(f.Data) {
		n = f.Data[f.Offset]
		f.Offset++
		f.Hash.Write([]byte{n})
		f.Checksum = ^f.Hash.Sum32()
	}
	return
}

func upsWrite(f *file, n byte) {
	if f.Offset < len(f.Data) {
		f.Data[f.Offset] = n
		f.Hash.Write([]byte{n})
		f.Checksum = ^f.Hash.Sum32()
	}
	f.Offset++
}

func upsDecode(f *file) int {
	var offset = 0
	var shift = 1
	for {
		x := upsRead(f)
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

	if upsRead(data.Patch) != 'U' ||
		upsRead(data.Patch) != 'P' ||
		upsRead(data.Patch) != 'S' ||
		upsRead(data.Patch) != '1' {
		return nil, errors.New("invalid patch header")
	}

	sourceReadLength := upsDecode(data.Patch)
	targetReadLength := upsDecode(data.Patch)

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
		for length := upsDecode(data.Patch); length > 0; length-- {
			upsWrite(data.Target, upsRead(data.Source))
		}
		for {
			patchXOR := upsRead(data.Patch)
			upsWrite(data.Target, patchXOR^upsRead(data.Source))
			if patchXOR == 0 {
				break
			}
		}
	}

	for data.Source.Offset < len(data.Source.Data) {
		upsWrite(data.Target, upsRead(data.Source))
	}
	for data.Target.Offset < len(data.Target.Data) {
		upsWrite(data.Target, upsRead(data.Source))
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
		sourceReadChecksum |= uint32(upsRead(data.Patch)) << (i * 8)
	}
	var targetReadChecksum uint32
	for i := 0; i < 4; i++ {
		targetReadChecksum |= uint32(upsRead(data.Patch)) << (i * 8)
	}

	patchResultChecksum := ^data.Patch.Checksum
	data.Source.Checksum = ^data.Source.Checksum
	data.Target.Checksum = ^data.Target.Checksum

	var patchReadChecksum uint32
	for i := 0; i < 4; i++ {
		patchReadChecksum |= uint32(upsRead(data.Patch)) << (i * 8)
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
