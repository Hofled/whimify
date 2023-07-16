package sort

import "fmt"

func C() {
}

func D() {
	B()
}

// A documentation
func A() {
	B()
	C()
}

func B() {
	fmt.Println("B")
	C()
}

// R recursive function
func R(b bool) {
	if b {
		R(false)
	}
}
