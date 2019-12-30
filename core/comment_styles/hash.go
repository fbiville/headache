package styles

type Hash struct{}

func (Hash) GetName() string {
	return "Hash"
}
func (Hash) GetOpeningString() string {
	return ""
}
func (Hash) GetString() string {
	return "# "
}
func (Hash) GetClosingString() string {
	return ""
}
