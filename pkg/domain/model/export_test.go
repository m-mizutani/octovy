package model

func (x *Report) Sources() map[string]*SourceChanges {
	return x.sources
}
