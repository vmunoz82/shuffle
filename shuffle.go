package shuffle

import (
	"errors"
	"math/bits"
)

func masks(max FeistelWord) (int, FeistelWord) {
	/* special case when you want to traverse all the inclusive range of FeistelWord type space */
	if max == 0 {
		max = MaxFeistelWord
	}
	/* bit offset and bit mask */
	bitOffset := (bits.Len64(uint64(max-1)) + 1) >> 1
	bitMask := FeistelWord(1)<<bitOffset - 1

	return bitOffset, bitMask
}

// RandomIndex returns the shuffled index of a given index and the higher bound of the set to be shuffled.
// If max is 0 then max act as the FesitelWord max value +1.
func RandomIndex(idx, max FeistelWord, cipher *Feistel) (FeistelWord, error) {
	if len(cipher.keys) == 0 {
		return MaxFeistelWord, errors.New("RandomIndex: Feistel context have a zero rounds key")
	}
	if max != 0 && idx >= max {
		return MaxFeistelWord, errors.New("RandomIndex: requested index surpass higher bound")
	}

	bitOffset, bitMask := masks(max)

	permutation := idx
	for {
		left, right := cipher.Cipher(permutation>>bitOffset, permutation&bitMask, bitMask)
		permutation = left<<bitOffset | right
		if permutation < max || max == 0 {
			return permutation, nil
		}
	}
}

// GetIndex returns the unshuffled index of a given shuffle index and the higher bound of the shuffled set.
// If max is 0 then max act as the FesitelWord max value +1.
func GetIndex(permutation, max FeistelWord, cipher *Feistel) (FeistelWord, error) {
	if len(cipher.keys) == 0 {
		return MaxFeistelWord, errors.New("GetIndex: Feistel context have a zero rounds key")
	}
	if max != 0 && permutation >= max {
		return MaxFeistelWord, errors.New("GetIndex: requested permutation surpass higher bound")
	}

	bitOffset, bitMask := masks(max)
	idx := permutation
	//fmt.Println(idx)
	for {
		left, right := cipher.Decipher(idx>>bitOffset, idx&bitMask, bitMask)
		idx = left<<bitOffset | right
		if idx < max || max == 0 {
			return idx, nil
		}
	}

}

// Shuffle creates a channel that stream shuffled values from a given minimum and maximum.
// If max is 0 then max act as the FesitelWord max value +1.
func Shuffle(min, max FeistelWord, cipher *Feistel) (<-chan FeistelWord, error) {
	if len(cipher.keys) == 0 {
		return nil, errors.New("Shuffle: Feistel context have a zero rounds key")
	}
	shuffle := make(chan FeistelWord)

	if min > max {
		min, max = max, min
	}

	diff := max - min
	go func(max, offset FeistelWord) {
		var permutation FeistelWord

		bitOffset, bitMask := masks(max)
		for i := FeistelWord(0); i < max || max == 0; i++ {
			permutation = i
			for {
				left, right := cipher.Cipher(permutation>>bitOffset, permutation&bitMask, bitMask)
				permutation = left<<bitOffset | right
				if permutation < max || max == 0 {
					shuffle <- permutation + offset
					break
				}
			}
			if i == MaxFeistelWord {
				break
			}
		}
		close(shuffle)
	}(diff, min)
	return shuffle, nil
}
