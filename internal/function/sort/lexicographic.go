package sort

func lexicographic(t1, t2 *dependencyTree) bool {
	return t1.Name() < t2.Name()
}
