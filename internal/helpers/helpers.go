package helpers

func DefaultString(s, fallback string) string {
	if s == "" {
		return fallback
	}

	return s
}
