package shuffle

import (
	"fmt"
	"reflect"
	"testing"
)

func boolXor(a, b bool) bool {
	return (a && !b) || (!a && b)
}

var keys = []FeistelWord{
	0xA45CF355C3B1CD88,
	0x8B9271CC2FC9365A,
	0x33CD458F23C816B1,
	0xC026F9D152DE23A9}

/* Function for each round, could be anything don't need to be reversible */
func roundFunction(v, key FeistelWord) FeistelWord {
	return (v * 941083987) ^ (key >> (v & 7) * 104729)
}

/* This is usefull for coverage */
func TestFeistelConstructor(t *testing.T) {
	NewFeistel([]FeistelWord{}, roundFunction)
}

func TestFeistelNoRounds(t *testing.T) {
	fn := NewFeistelDefault([]FeistelWord{})
	GotLeft, GotRight := fn.Cipher(0x1234, 0x5678, 0xFFFF)
	if GotLeft != 0x1234 || GotRight != 0x5678 {
		t.Errorf("TestFeistelNoRounds: Expected values 0x1234 and 0x5678 and got %x, %x\n", GotLeft, GotRight)
	}
}

func TestShuffleNoRounds(t *testing.T) {
	fn := NewFeistelDefault([]FeistelWord{})

	for _, b := range []bool{false, true} {
		GotShuffle, err := Shuffle(0, 255, fn)
		if boolXor(b, GotShuffle != nil) || boolXor(b, err == nil) {
			t.Errorf("TestShuffleNoRounds: Unexpected values, got %x, %x\n", GotShuffle, err)
		}
		GotRandomIndex, err2 := RandomIndex(0, 255, fn)
		if boolXor(b, GotRandomIndex != MaxFeistelWord) || boolXor(b, err2 == nil) {
			t.Errorf("TestShuffleNoRounds: Unexpected values, got %x, %x\n", GotRandomIndex, err2)
		}
		GotGetIndex, err3 := GetIndex(0, 255, fn)
		if boolXor(b, GotGetIndex != MaxFeistelWord) || boolXor(b, err3 == nil) {
			t.Errorf("TestShuffleNoRounds: Unexpected values, got %x, %x\n", GotGetIndex, err3)
		}
		/* switch the condition checks */
		fn = NewFeistelDefault(keys)
	}
}

func TestShuffleInverseMinMax(t *testing.T) {
	normal, _ := Shuffle(0, 10, NewFeistelDefault(keys))
	inverse, _ := Shuffle(10, 0, NewFeistelDefault(keys))

	for n := range normal {
		i := <-inverse
		if n != i {
			t.Errorf("TestShuffleInverseMinMax: Dispair values: %d != %d\n", n, i)
		}
	}
}

/* This is usefull for coverage */
func TestShuffleWholeIntegerSpace(t *testing.T) {
	if testing.Short() || reflect.TypeOf(FeistelWord(0)).Size() > 2 {
		t.Skip("TestShuffleWholeIntegerSpace: Skipped")
	}
	/* max equals zero means traverse the whole integer space */
	s, _ := Shuffle(0, 0, NewFeistelDefault(keys))

	acc := uint64(0)
	acc2 := uint64(0)
	i := uint64(0)
	for v := range s {
		acc += uint64(v)
		acc2 += i
		i++
	}
	if acc != acc2 {
		t.Errorf("TestShuffleWholeIntegerSpace: Dispair values: %d != %d\n", acc, acc2)
	}
}

func TestShuffleRandomIndexGeneration(t *testing.T) {
	fn := NewFeistelDefault(keys)

	for i := FeistelWord(0); i < 1000; i++ {
		index, err := RandomIndex(i, 1000, fn)
		if err != nil {
			break
		}
		rindex, rerr := GetIndex(index, 1000, fn)
		if rerr != nil {
			break
		}
		if i != rindex {
			t.Errorf("TestShuffleRandomIndexGeneration: Could not recover original index: %d != %d\n", i, rindex)
		}
	}
}

func TestShuffleRandomIndexOutOfBounds(t *testing.T) {
	permutation, err := RandomIndex(10, 9, NewFeistelDefault(keys))
	if permutation != MaxFeistelWord || err == nil {
		t.Errorf("TestShuffleRandomIndexOutOfBounds: Unexpected results: %d != %s\n", permutation, err)
	}
}

func TestShuffleGetIndexOutOfBounds(t *testing.T) {
	index, err := GetIndex(10, 9, NewFeistelDefault(keys))
	if index != MaxFeistelWord || err == nil {
		t.Errorf("TestShuffleGetIndexOutOfBounds: Unexpected results: %d != %s\n", index, err)
	}
}

/* simplest possible benchmark, there is nothing that prevent this algorithm to run in parallel */
func BenchmarkRandomIndexParallel(b *testing.B) {
	fn := NewFeistelDefault(keys)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			RandomIndex(0, 12345, fn)
		}
	})
}

func ExampleShuffle() {
	s, _ := Shuffle(1000, 1020, NewFeistel(keys, roundFunction))
	for v := range s {
		fmt.Println(v)
	}
	// Output:
	// 1003
	// 1005
	// 1006
	// 1009
	// 1004
	// 1011
	// 1008
	// 1016
	// 1010
	// 1001
	// 1017
	// 1000
	// 1013
	// 1012
	// 1015
	// 1014
	// 1007
	// 1002
	// 1018
	// 1019
}
