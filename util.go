package main

func Max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func Abs(a int) int {
	if a < 0 {
		return -1 * a
	}
	return a
}