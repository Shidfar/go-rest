package types

type Ptr struct {
	To TypeDef
}

func (p *Ptr) Type(uctx UsageCtx) string {
	return "*" + p.To.Type(uctx)
}
