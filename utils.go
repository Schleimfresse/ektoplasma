package main

func ConvertBoolToInt(expr bool) interface{} {
	if expr {
		return 1
	} else {
		return 0
	}
}
