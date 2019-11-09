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
	patch := &file{
		Data: patchData,
		Hash: crc32.NewIEEE(),
	}
	source := &file{
		Data: sourceData,
		Hash: crc32.NewIEEE(),
	}
	target := &file{
		Hash: crc32.NewIEEE(),
	}

	if len(patch.Data) < 18 {
		return nil, errors.New("patch too small")
	}

	if upsRead(patch) != 'U' ||
		upsRead(patch) != 'P' ||
		upsRead(patch) != 'S' ||
		upsRead(patch) != '1' {
		return nil, errors.New("invalid patch header")
	}

	sourceReadLength := upsDecode(patch)
	targetReadLength := upsDecode(patch)

	if len(source.Data) != sourceReadLength &&
		len(source.Data) != targetReadLength {
		return nil, errors.New("invalid source")
	}

	targetLength := sourceReadLength
	if len(source.Data) == sourceReadLength {
		targetLength = targetReadLength
	}

	prov := make([]byte, targetLength)
	target.Data = prov

	for patch.Offset < len(patch.Data)-12 {
		for length := upsDecode(patch); length > 0; length-- {
			upsWrite(target, upsRead(source))
		}
		for {
			patchXOR := upsRead(patch)
			upsWrite(target, patchXOR^upsRead(source))
			if patchXOR == 0 {
				break
			}
		}
	}

	for source.Offset < len(source.Data) {
		upsWrite(target, upsRead(source))
	}
	for target.Offset < len(target.Data) {
		upsWrite(target, upsRead(source))
	}

	if err := checks(patch, source, target, sourceReadLength, targetReadLength); err != nil {
		return nil, err
	}
	return &target.Data, nil
}

// checks verifies that the patching process went well by comparing checksums
func checks(patch, source, target *file, sourceReadLength, targetReadLength int) error {
	var sourceReadChecksum uint32
	for i := 0; i < 4; i++ {
		sourceReadChecksum |= uint32(upsRead(patch)) << (i * 8)
	}
	var targetReadChecksum uint32
	for i := 0; i < 4; i++ {
		targetReadChecksum |= uint32(upsRead(patch)) << (i * 8)
	}

	patchResultChecksum := ^patch.Checksum
	source.Checksum = ^source.Checksum
	target.Checksum = ^target.Checksum

	var patchReadChecksum uint32
	for i := 0; i < 4; i++ {
		patchReadChecksum |= uint32(upsRead(patch)) << (i * 8)
	}

	if patchResultChecksum != patchReadChecksum {
		return errors.New("invalid patch")
	}

	if source.Checksum == sourceReadChecksum && len(source.Data) == sourceReadLength {
		if target.Checksum == targetReadChecksum && len(target.Data) == targetReadLength {
			return nil
		}
		return errors.New("invalid target")
	} else if source.Checksum == targetReadChecksum && len(source.Data) == targetReadLength {
		if target.Checksum == sourceReadChecksum && len(target.Data) == sourceReadLength {
			return nil
		}
		return errors.New("invalid target")
	}

	return errors.New("invalid source")
}
