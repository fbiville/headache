package styles

type Rem struct {
}

func (Rem) GetName() string {
	return "REM"
}
func (Rem) GetOpeningString() string {
	return ""
}
func (Rem) GetString() string {
	return "REM "
}
func (Rem) GetClosingString() string {
	return ""
}
