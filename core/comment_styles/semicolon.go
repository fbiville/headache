package styles

type SemiColon struct {
}

func (SemiColon) GetName() string {
	return "SemiColon"
}
func (SemiColon) GetOpeningString() string {
	return ""
}
func (SemiColon) GetString() string {
	return "; "
}
func (SemiColon) GetClosingString() string {
	return ""
}
