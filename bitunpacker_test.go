package itc

import "fmt"

func ExamplePopSimple() {
	bup := NewBitUnPacker([]byte{0x44, 0x00, 0x00, 0x01, 0x80, 0x00, 0x00, 0x00})
	fmt.Printf("Base = 010001\n")
	fmt.Printf("Pop(3) = %d\n", bup.Pop(3))
	fmt.Printf("Pop(1) = %d\n", bup.Pop(1))
	fmt.Printf("Pop(2) = %d\n", bup.Pop(2))
	fmt.Printf("Pop(25) = %d\n", bup.Pop(25))
	fmt.Printf("Pop(2) = %d\n", bup.Pop(2))
	// Output:
	// Base = 010001
	// Pop(3) = 2
	// Pop(1) = 0
	// Pop(2) = 1
	// Pop(25) = 0
	// Pop(2) = 3
}
