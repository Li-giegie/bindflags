package bindflags

type FlagTag struct {
	Name  string
	Value string
	Usage string
}

func (f *FlagTag) GetName() string {
	return f.Name
}
func (f *FlagTag) GetValue() string {
	return f.Value
}
func (f *FlagTag) GetUsage() string {
	return f.Usage
}
