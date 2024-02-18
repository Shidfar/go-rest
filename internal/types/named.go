package types

type Named struct {
	Name  string
	Pkg   string
	Iface bool
}

func (s *Named) Dynamic() bool {
	return s.Pkg == "model" && s.Iface
}

func (s *Named) Type(uctx UsageCtx) string {
	if uctx == DynamicField && s.Dynamic() {
		return "x" + s.Name
	}
	if s.Pkg == "" {
		return s.Name
	}
	return s.Pkg + "." + s.Name
}
