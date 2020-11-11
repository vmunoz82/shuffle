# shuffle
This golang module should be used to produce a shuffled set using a Feistel network.

The algorithm has the following features:
* Performs the shuffle in a complete stateless form (ie.: there is no need for an input, even neither an out set), as you can traverse a channel with each of the output elements.
* It can run it safely in parallel which is not possible with algorithms that need to perform swap to elements of an array.
* It is keyed so you can provide a seed (in form of a Feistel function and round keys) to recall the shuffle order.
* It is random accessible so if you need very few shuffles elements of a large set you only ask those.
* It is reversible so given a shuffled value you can obtain the original unshuffled value.

## Installation
```
go get github.com/vmunoz82/shuffle
```

## Examples

To shuffle elements of the range [1000, 1020)

```
package main

import (
        "github.com/vmunoz82/shuffle"
        "fmt"
)

var keys = []shuffle.FeistelWord{
        0xA45CF355C3B1CD88,
        0x8B9271CC2FC9365A,
        0x33CD458F23C816B1,
        0xC026F9D152DE23A9}

/* Function for each round, could be anything don't need to be reversible */
func roundFunction(v, key shuffle.FeistelWord) shuffle.FeistelWord {
        return (v * 941083987) ^ (key >> (v & 7) * 104729)
}

func main() {
        s, _ := shuffle.Shuffle(1000, 1020, shuffle.NewFeistel(keys, roundFunction))
        for v := range s {
                fmt.Println(v)
        }
}
```
This will output
```
1003
1005
1006
1009
1004
1011
1008
1016
1010
1001
1017
1000
1013
1012
1015
1014
1007
1002
1018
1019
```

In this example it will generate the first shuffled value for the range [0, 1000), that is the value 846, and next we made the inverse operation, given the value 846 asks for which element of the original set produce it.

```
package main

import (
        "fmt"
        "github.com/vmunoz82/shuffle"
)

var keys = []shuffle.FeistelWord{
        0xA45CF355C3B1CD88,
        0x8B9271CC2FC9365A,
        0x33CD458F23C816B1,
        0xC026F9D152DE23A9}

/* Function for each round, could be anything don't need to be reversible */
func roundFunction(v, key shuffle.FeistelWord) shuffle.FeistelWord {
        return (v * 941083987) ^ (key >> (v & 7) * 104729)
}

func main() {
        fn := shuffle.NewFeistel(keys, roundFunction)
        v, _ := shuffle.RandomIndex(0, 1000, fn)
        fmt.Println(v)
        v, _ = shuffle.GetIndex(846, 1000, fn)
        fmt.Println(v)
}
```
This example will output
```
846
0
```
