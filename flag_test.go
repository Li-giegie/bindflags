package bindflags

import (
	"flag"
	"fmt"
	"testing"
)

type tf int

func (tf) GetFlag() IFlagTag {
	return &FlagTag{
		Name:  "f",
		Value: "100",
		Usage: "case f usage",
	}
}

type F struct {
	A string  `flag:"name:a;value:abc;usage:case a usage"`
	B int     `flag:"b;123;case b usage"`
	C bool    `flag:"name:c;usage:case c usage"`
	D float64 `flag:"3.1415926;usage:case d usage;name:d"`
	E uint    `flag:"value:1"`
	d string
	F tf
}

func TestBindFlags(t *testing.T) {
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	e := &F{}
	err := BindFlags(f, e)
	if err != nil {
		t.Fatal(err)
		return
	}
	f.PrintDefaults()
	fmt.Printf("%#v\n", e)
}
