package types

type Slice struct {
	Vals TypeDef
}

func (t *Slice) Type(uctx UsageCtx) string {
	return "[]" + t.Vals.Type(uctx)
}
