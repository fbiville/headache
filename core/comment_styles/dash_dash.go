package styles

type DashDash struct {
}

func (DashDash) GetName() string {
	return "DashDash"
}
func (DashDash) GetOpeningString() string {
	return ""
}
func (DashDash) GetString() string {
	return "-- "
}
func (DashDash) GetClosingString() string {
	return ""
}
