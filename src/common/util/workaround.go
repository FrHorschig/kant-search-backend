package util

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
