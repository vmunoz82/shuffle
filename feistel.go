package shuffle

const constA = 654188429
const constB = 104729

// A FeistelWord abstract the underlying unsigned type to represent a word in the Feistel context.
type FeistelWord uint64

// MaxFeistelWord is the maximum possible value of the FesitelWord type.
const MaxFeistelWord = ^FeistelWord(0)

// A FeistelFunc is the type that define a function to be used in each round of the Feistel application.
type FeistelFunc func(FeistelWord, FeistelWord) FeistelWord

// Feistel is where the state of Feistel network is stored, contains the keys to be applied on each round and a function of type FeistelFunc.
type Feistel struct {
	keys []FeistelWord
	f    FeistelFunc
}

func defaultRoundFunction(v, key FeistelWord) FeistelWord {
	return (v * constA) ^ (key * constB)
}

// NewFeistel constructs a new Feistel context with the given keys and FeistelFunc.
func NewFeistel(keys []FeistelWord, f FeistelFunc) *Feistel {
	return &Feistel{keys, f}
}

// NewFeistelDefault constructs a new Feistel context with the given keys and using a default FeistelFunc.
func NewFeistelDefault(keys []FeistelWord) *Feistel {
	return &Feistel{keys, defaultRoundFunction}
}

func (f *Feistel) core(a, b, amask, bmask FeistelWord) (FeistelWord, FeistelWord) {
	var a2 FeistelWord

	mask := (amask | bmask)
	rounds := len(f.keys)
	/* passtrhu when there are not rounds */
	if rounds == 0 {
		return a & mask, b & mask
	}

	a1 := b
	b1 := a ^ f.f(b, (f.keys[0]&amask)|(f.keys[rounds-1]&bmask))&mask
	for round := 1; round < rounds; round++ {
		a2 = a1
		a1 = b1
		b1 = (a2 ^ f.f(b1, (f.keys[round]&amask)|(f.keys[rounds-1-round]&bmask))) & mask
	}
	return a1, b1
}

// Cipher applys the Feistel cipher step to a word described in its left and right parts.
func (f *Feistel) Cipher(left, right, mask FeistelWord) (FeistelWord, FeistelWord) {
	return f.core(left, right, mask, 0)
}

// Decipher applys the Fesitel decipher step to a word described in its left and right parts.
func (f *Feistel) Decipher(left, right, mask FeistelWord) (FeistelWord, FeistelWord) {
	right, left = f.core(right, left, 0, mask)
	return left, right
}
