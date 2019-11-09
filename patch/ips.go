package patch

import (
	"errors"
)

func ipsAllocTargetData(patchData, sourceData []byte) ([]byte, error) {
	offset := 5
	targetLength := len(sourceData)

	for {
		if offset > len(patchData)-3 {
			break
		}

		address := int(patchData[offset]) << 16
		offset++
		address |= int(patchData[offset]) << 8
		offset++
		address |= int(patchData[offset]) << 0
		offset++

		if address == 0x454f46 /* EOF */ {
			if offset == len(patchData) {
				prov := make([]byte, targetLength)
				return prov, nil
			} else if offset == len(patchData)-3 {
				size := int(patchData[offset]) << 16
				offset++
				size |= int(patchData[offset]) << 8
				offset++
				size |= int(patchData[offset]) << 0
				offset++
				targetLength = size
				prov := make([]byte, targetLength)
				return prov, nil
			}
		}

		if offset > len(patchData)-2 {
			break
		}

		length := int(patchData[offset]) << 8
		offset++
		length |= int(patchData[offset]) << 0
		offset++

		if length > 0 /* Copy */ {
			if offset > len(patchData)-int(length) {
				break
			}

			for length > 0 {
				address++
				offset++
				length--
			}
		} else /* RLE */ {
			if offset > len(patchData)-3 {
				break
			}

			length := int(patchData[offset]) << 8
			offset++
			length |= int(patchData[offset]) << 0
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

func applyIPS(patchData, sourceData []byte) (*[]byte, error) {
	if len(patchData) < 8 ||
		patchData[0] != 'P' ||
		patchData[1] != 'A' ||
		patchData[2] != 'T' ||
		patchData[3] != 'C' ||
		patchData[4] != 'H' {
		return nil, errors.New("invalid patch header")
	}

	targetData, err := ipsAllocTargetData(patchData, sourceData)
	if err != nil {
		return nil, err
	}

	copy(targetData, sourceData)

	offset := 5
	for {
		if offset > len(patchData)-3 {
			break
		}

		address := int(patchData[offset]) << 16
		offset++
		address |= int(patchData[offset]) << 8
		offset++
		address |= int(patchData[offset]) << 0
		offset++

		if address == 0x454f46 /* EOF */ {
			if offset == len(patchData) {
				return &targetData, nil
			} else if offset == len(patchData)-3 {
				size := int(patchData[offset]) << 16
				offset++
				size |= int(patchData[offset]) << 8
				offset++
				size |= int(patchData[offset]) << 0
				offset++
				return &targetData, nil
			}
		}

		if offset > len(patchData)-2 {
			break
		}

		length := int(patchData[offset]) << 8
		offset++
		length |= int(patchData[offset]) << 0
		offset++

		if length > 0 /* Copy */ {
			if offset > len(patchData)-length {
				break
			}

			for ; length > 0; length-- {
				targetData[address] = patchData[offset]
				address++
				offset++
			}
		} else /* RLE */ {
			if offset > len(patchData)-3 {
				break
			}

			length = int(patchData[offset]) << 8
			offset++
			length |= int(patchData[offset]) << 0
			offset++

			if length == 0 /* Illegal */ {
				break
			}

			for ; length > 0; length-- {
				targetData[address] = patchData[offset]
				address++
			}

			offset++
		}
	}

	return nil, errors.New("invalid patch")
}
