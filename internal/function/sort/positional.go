package sort

func positional(t1, t2 *dependencyTree) bool {
	return t1.StartPos() < t2.StartPos()
}
