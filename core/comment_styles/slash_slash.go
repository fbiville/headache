package styles

type SlashSlash struct{}

func (SlashSlash) GetName() string {
	return "SlashSlash"
}
func (SlashSlash) GetOpeningString() string {
	return ""
}
func (SlashSlash) GetString() string {
	return "// "
}
func (SlashSlash) GetClosingString() string {
	return ""
}
