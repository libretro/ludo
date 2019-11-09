package patch

import (
	"errors"
)

func ipsAllocTargetData(patch, source []byte) ([]byte, error) {
	offset := 5
	targetLength := len(source)

	for {
		if offset > len(patch)-3 {
			break
		}

		address := int(patch[offset]) << 16
		offset++
		address |= int(patch[offset]) << 8
		offset++
		address |= int(patch[offset]) << 0
		offset++

		if address == 0x454f46 /* EOF */ {
			if offset == len(patch) {
				prov := make([]byte, targetLength)
				return prov, nil
			} else if offset == len(patch)-3 {
				size := int(patch[offset]) << 16
				offset++
				size |= int(patch[offset]) << 8
				offset++
				size |= int(patch[offset]) << 0
				offset++
				targetLength = size
				prov := make([]byte, targetLength)
				return prov, nil
			}
		}

		if offset > len(patch)-2 {
			break
		}

		length := int(patch[offset]) << 8
		offset++
		length |= int(patch[offset]) << 0
		offset++

		if length > 0 /* Copy */ {
			if offset > len(patch)-int(length) {
				break
			}

			for length > 0 {
				address++
				offset++
				length--
			}
		} else /* RLE */ {
			if offset > len(patch)-3 {
				break
			}

			length := int(patch[offset]) << 8
			offset++
			length |= int(patch[offset]) << 0
			offset++

			if length == 0 /* Illegal */ {
				break
			}

			for length > 0 {
				address++
				length--
			}

			offset++
		}

		if address > targetLength {
			targetLength = address
		}
	}

	return nil, errors.New("invalid patch")
}

func applyIPS(patch, source []byte) (*[]byte, error) {
	if len(patch) < 8 {
		return nil, errors.New("patch too small")
	}

	if string(patch[0:5]) != "PATCH" {
		return nil, errors.New("invalid patch header")
	}

	targetData, err := ipsAllocTargetData(patch, source)
	if err != nil {
		return nil, err
	}

	copy(targetData, source)

	offset := 5
	for {
		if offset > len(patch)-3 {
			break
		}

		address := int(patch[offset]) << 16
		offset++
		address |= int(patch[offset]) << 8
		offset++
		address |= int(patch[offset]) << 0
		offset++

		if address == 0x454f46 /* EOF */ {
			if offset == len(patch) {
				return &targetData, nil
			} else if offset == len(patch)-3 {
				size := int(patch[offset]) << 16
				offset++
				size |= int(patch[offset]) << 8
				offset++
				size |= int(patch[offset]) << 0
				offset++
				return &targetData, nil
			}
		}

		if offset > len(patch)-2 {
			break
		}

		length := int(patch[offset]) << 8
		offset++
		length |= int(patch[offset]) << 0
		offset++

		if length > 0 /* Copy */ {
			if offset > len(patch)-length {
				break
			}

			for ; length > 0; length-- {
				targetData[address] = patch[offset]
				address++
				offset++
			}
		} else /* RLE */ {
			if offset > len(patch)-3 {
				break
			}

			length = int(patch[offset]) << 8
			offset++
			length |= int(patch[offset]) << 0
			offset++

			if length == 0 /* Illegal */ {
				break
			}

			for ; length > 0; length-- {
				targetData[address] = patch[offset]
				address++
			}

			offset++
		}
	}

	return nil, errors.New("invalid patch")
}
