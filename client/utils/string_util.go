package utils

// return string is empty or not
func IsEmpty(s string) bool {
	return s == ""
}

// return any string is empty or not
func IsAnyStringEmpty(ss ...string) bool {
	for _, s := range ss {
		if IsEmpty(s) {
			return true
		}
	}
	return false
}

// panic if any string is empty
func PanicIfAnyStringEmpty(ss ...string) {
	if IsAnyStringEmpty(ss...) {
		panic("any string is empty")
	}
}
