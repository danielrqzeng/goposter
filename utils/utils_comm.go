package utils

func MinInt(v1, v2 int) (min int) {
	if v1 > v2 {
		return v2
	}
	return v1
}
