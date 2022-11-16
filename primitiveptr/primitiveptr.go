package primitiveptr

func String(val string) *string {
	return &val
}

func StringDeref(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}
