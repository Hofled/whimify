package sort

import "fmt"

func C() {
}

func D() {
	B()
}

func B() {
	fmt.Println("B")
	C()
}

// A documentation
func A() {
	B()
	C()
}
