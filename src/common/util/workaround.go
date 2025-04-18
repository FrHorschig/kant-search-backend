package util

func StrPtr(s string) *string {
	return &s
}

func StrVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func FalsePtr() *bool {
	f := false
	return &f
}

func IntPtr(i int) *int {
	return &i
}
