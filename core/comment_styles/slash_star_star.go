package styles

type SlashStarStar struct{}

func (SlashStarStar) GetName() string {
	return "SlashStarStar"
}
func (SlashStarStar) GetOpeningString() string {
	return "/**"
}
func (SlashStarStar) GetString() string {
	return " * "
}
func (SlashStarStar) GetClosingString() string {
	return " */"
}
