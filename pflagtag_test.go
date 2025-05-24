package bindflags

import (
	"fmt"
	"testing"
)

func TestPFlagTag(t *testing.T) {
	testCases := []string{
		"name:123;Shorthand:asd;Value:as;Usage:123",
		"123;asd;;123",
		"123;asd;'a;asd;asdl;s';123",
		"name:123;asd;as;123",
		"123;Shorthand:asd;as;123",
		"name:123;Shorthand:asd;as;123",
		"name:foo;shorthand:f;value:'hello;world';usage:测试",
		"name:\"fo;o';shorthand:\"f;value:hello;usage:'测试;用例'",
		"name:123;Shorthand:asd;Value:\"a;asd\"s;Usage:123",
		"name:123;Shorthand:asd;Value:'a;asd's;Usage:123",
		"name:123;Shorthand:asd;Value:a;s;Usage:123",
		"name:123;Value:\"'a;123'\";Shorthand:asd;asd",
	}

	for _, tc := range testCases {
		fmt.Printf("in: %s\n", tc)
		result, err := scanPFlagTag(tc)
		if err != nil {
			println(err.Error())
		} else {
			fmt.Printf("out:  %v\n", result)
		}
	}
}
