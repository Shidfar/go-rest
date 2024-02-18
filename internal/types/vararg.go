package types

type Vararg struct {
	Vals TypeDef
}

func (s *Vararg) Type(uctx UsageCtx) string {
	switch uctx {
	case Signature:
		return "..." + s.Vals.Type(uctx)
	case Field, DynamicField:
		return "[]" + s.Vals.Type(uctx)
	default:
		panic("bad usage context")
	}
}
