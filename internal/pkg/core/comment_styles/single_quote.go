package styles

type SingleQuote struct {
}

func (SingleQuote) GetName() string {
	return "SingleQuote"
}
func (SingleQuote) GetOpeningString() string {
	return ""
}
func (SingleQuote) GetString() string {
	return "' "
}
func (SingleQuote) GetClosingString() string {
	return ""
}
