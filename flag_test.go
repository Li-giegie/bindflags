package bindflags

import (
	"github.com/spf13/pflag"
	"testing"
)

func TestBindFlags(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	s := new(student)
	MustBindFlags(f, s)
	f.PrintDefaults()
	//  -a, --age int       age
	//  -d, --desc string   desc
	//  -n, --name string   name (default "ss")
	//  -s, --sex           sex (default true)
	//PASS
}

type student struct {
	Name string `flag:"Name:name;shorthand:n;value:ss;usage:name of student"`
	Age  int    `flag:"age;a;0;usage:age of student"`
	Sex  bool   `flag:"sex;s;true;sex"`
	Desc studentDesc
}

type studentDesc string

func (s studentDesc) GetFlagTag() IFlagTag {
	return &FlagTag{
		Name:      "desc",
		Shorthand: "d",
		Value:     "",
		Usage:     "student description",
	}
}
