package main


// Find returns the smallest index i at which x == a[i],
// or len(a) if there is no such index.
func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}


func Remove(a []string, value string) []string {
	i := Find(a, value)

	if i >= 0 && i < len(a) {
		a[i] = a[len(a)-1] // Copy last element to index i
		a[len(a)-1] = ""   // Erase last element (write zero value)
		a = a[:len(a)-1]   // Truncate slice
	}

	return a
}