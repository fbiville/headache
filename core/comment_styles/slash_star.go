package styles

type SlashStar struct{}

func (SlashStar) GetName() string {
	return "SlashStar"
}
func (SlashStar) GetOpeningString() string {
	return "/*"
}
func (SlashStar) GetString() string {
	return " * "
}
func (SlashStar) GetClosingString() string {
	return " */"
}
