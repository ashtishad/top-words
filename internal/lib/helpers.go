package lib

import "runtime"

// SetGoMaxProcs sets the maximum number of CPUs that can be executing concurrently.
// returns the number, default is 1.
func SetGoMaxProcs() int {
	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)
	return n
}
