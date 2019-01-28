package core

import "github.com/fbiville/headache/helper"

type VersionedHeaderTemplate struct {
	Current  *HeaderTemplate
	Previous *HeaderTemplate
	Revision string
}

func (t VersionedHeaderTemplate) RequiresFullScan() bool {
	return t.Revision == "" ||
		!helper.SliceEqual(t.Current.Lines, t.Previous.Lines) ||
		!helper.SliceEqual(helper.Keys(t.Current.Data), helper.Keys(t.Previous.Data))
}
