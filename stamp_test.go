package itc

import (
	"fmt"
	"testing"
)

func Example1() {
	a := NewStamp()
	fmt.Printf("stage 1: a := NewStamp()\n")
	fmt.Printf("a: %s\n", a)
	// Output:
	// stage 1: a := NewStamp()
	// a: (1, 0)
}

func Example2() {
	a := NewStamp()
	b := a.Fork()
	fmt.Printf("stage 2: b := a.Fork()\n")
	fmt.Printf("a: %s\n", a)
	fmt.Printf("b: %s\n", b)

	// Output:
	// stage 2: b := a.Fork()
	// a: ((1, 0), 0)
	// b: ((0, 1), 0)
}

func Example3() {
	a := NewStamp()
	b := a.Fork()
	a.Event()
	b.Event()
	fmt.Printf("stage 3: a.Event(), b.Event()\n")
	fmt.Printf("a: %s\n", a)
	fmt.Printf("b: %s\n", b)

	// Output:
	// stage 3: a.Event(), b.Event()
	// a: ((1, 0), (0, 1, 0))
	// b: ((0, 1), (0, 0, 1))
}

func Example4() {
	a := NewStamp()
	b := a.Fork()
	a.Event()
	b.Event()
	c := a.Fork()
	b.Event()
	fmt.Printf("stage 4: c := a.Fork(), b.Event()\n")
	fmt.Printf("a: %s\n", a)
	fmt.Printf("b: %s\n", b)
	fmt.Printf("c: %s\n", c)

	// Output:
	// stage 4: c := a.Fork(), b.Event()
	// a: (((1, 0), 0), (0, 1, 0))
	// b: ((0, 1), (0, 0, 2))
	// c: (((0, 1), 0), (0, 1, 0))
}

func Example5() {
	a := NewStamp() // 1
	b := a.Fork()   // 2
	a.Event()       // 3
	b.Event()       // 3
	c := a.Fork()   // 4
	b.Event()       // 4
	a.Event()       // 5
	b.Join(c)       // 5
	fmt.Printf("stage 5: a.Event(), b.Join(c)\n")
	fmt.Printf("a: %s\n", a)
	fmt.Printf("b: %s\n", b)

	// Output:
	// stage 5: a.Event(), b.Join(c)
	// a: (((1, 0), 0), (0, (1, 1, 0), 0))
	// b: (((0, 1), 1), (1, 0, 1))
}

func Example6() {
	a := NewStamp() // 1
	b := a.Fork()   // 2
	a.Event()       // 3
	b.Event()       // 3
	c := a.Fork()   // 4
	b.Event()       // 4
	a.Event()       // 5
	b.Join(c)       // 5
	c = b.Fork()    // 6

	fmt.Printf("stage 6: c = b.Fork()\n")
	fmt.Printf("a: %s\n", a)
	fmt.Printf("b: %s\n", b)
	fmt.Printf("c: %s\n", c)

	// Output:
	// stage 6: c = b.Fork()
	// a: (((1, 0), 0), (0, (1, 1, 0), 0))
	// b: (((0, 1), 0), (1, 0, 1))
	// c: ((0, 1), (1, 0, 1))
}

func Example7() {
	a := NewStamp() // 1
	b := a.Fork()   // 2
	a.Event()       // 3
	b.Event()       // 3
	c := a.Fork()   // 4
	b.Event()       // 4
	a.Event()       // 5
	b.Join(c)       // 5
	c = b.Fork()    // 6
	a.Join(b)       // 7

	fmt.Printf("stage 7: a.Join(b)\n")
	fmt.Printf("a: %s\n", a)
	fmt.Printf("c: %s\n", c)

	// Output:
	// stage 7: a.Join(b)
	// a: ((1, 0), (1, (0, 1, 0), 1))
	// c: ((0, 1), (1, 0, 1))
}

func Example8() {
	a := NewStamp() // 1
	b := a.Fork()   // 2
	a.Event()       // 3
	b.Event()       // 3
	c := a.Fork()   // 4
	b.Event()       // 4
	a.Event()       // 5
	b.Join(c)       // 5
	c = b.Fork()    // 6
	a.Join(b)       // 7
	a.Event()       // 8

	fmt.Printf("stage 8: a.Event()\n")
	fmt.Printf("a: %s\n", a)
	fmt.Printf("c: %s\n", c)

	// Output:
	// stage 8: a.Event()
	// a: ((1, 0), 2)
	// c: ((0, 1), (1, 0, 1))
}

func TestStampEmptyStringer(t *testing.T) {
	stampString := NewStamp().String()
	if stampString != "(1, 0)" {
		t.Errorf("stamp did not serialize as expected %q", stampString)
	}
}

func ExampleStamp_MarshalBinary_seed() {
	seed := NewStamp()
	seedData, _ := seed.MarshalBinary()
	fmt.Printf("%s = % x\n", seed, seedData)
	// Output:
	// (1, 0) = 30 00 00 00
}

func ExampleStamp_MarshalBinary_seedAfterFork() {
	seed := NewStamp()
	seed.Fork()
	seedData, _ := seed.MarshalBinary()
	fmt.Printf("%s = % x\n", seed, seedData)
	// Output:
	// ((1, 0), 0) = 8c 00 00 00
}

func ExampleStamp_UnmarshalBinary() {
	data := []byte{0x8c, 00, 00, 00}
	stamp := NewStamp()
	stamp.UnmarshalBinary(data)
	fmt.Printf("% x = %s\n", data, stamp)
	// Output:
	// 8c 00 00 00 = ((1, 0), 0)
}
