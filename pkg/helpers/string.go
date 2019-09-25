package helpers

func StringInSlice(str string, list []string) bool {
	for _, strCmp := range list {
		if strCmp == str {
			return true
		}
	}
	return false
}