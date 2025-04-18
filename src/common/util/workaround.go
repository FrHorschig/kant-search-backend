package util

// TODO adjust name
func ToStrPtr(s string) *string {
	return &s
}

func ToStrVal(s *string) string {
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
