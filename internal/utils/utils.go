package utils

// BoolAddr returns a pointer to a boolean, useful if you need a pointer to true.
func BoolAddr(b bool) *bool {
	return &b
}
