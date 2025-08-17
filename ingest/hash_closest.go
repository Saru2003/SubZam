package main

func HammingDistance(a, b uint64) int {
	x := a ^ b 
	count := 0
	for x != 0 {
		count++
		x &= x - 1
	}
	return count
}
