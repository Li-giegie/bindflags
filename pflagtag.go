package bindflags

type PFlagTag struct {
	Name      string
	Shorthand string
	Value     string
	Usage     string
}

func (f *PFlagTag) GetName() string {
	return f.Name
}
func (f *PFlagTag) GetShorthand() string {
	return f.Shorthand
}
func (f *PFlagTag) GetValue() string {
	return f.Value
}
func (f *PFlagTag) GetUsage() string {
	return f.Usage
}
