package arith_test

import (
	"testing"

	A "github.com/danwakefield/gosh/arith"
	"github.com/danwakefield/gosh/variables"
)

var EmptyScope = variables.NewScope()

func TestParserBinops(t *testing.T) {
	type TestCase struct {
		in   string
		want int64
	}
	cases := []TestCase{
		{"5 <= 4", A.ShellFalse},
		{"4 <= 4", A.ShellTrue},
		{"3 <= 4", A.ShellTrue},
		{"3 >= 4", A.ShellFalse},
		{"4 >= 4", A.ShellTrue},
		{"5 >= 4", A.ShellTrue},
		{"5 <  4", A.ShellFalse},
		{"3 <  4", A.ShellTrue},
		{"3 >  4", A.ShellFalse},
		{"5 >  4", A.ShellTrue},
		{"5 == 4", A.ShellFalse},
		{"4 == 4", A.ShellTrue},
		{"4 != 4", A.ShellFalse},
		{"5 != 4", A.ShellTrue},
		{"5 & 4", 4},
		{"3 & 4", 0},
		{"3 | 4", 7},
		{"4 | 4", 4},
		{"3 ^ 4", 7},
		{"4 ^ 4", 0},
		{"1 << 4", 16},
		{"16 >> 4", 1},
		{"10 % 4", 2},
		{"3 * 4", 12},
		{"12 / 4", 3},
		{"10 - 4", 6},
		{"10 + 4", 14},
	}

	for _, c := range cases {
		got, err := A.Parse(c.in, EmptyScope)
		if err != nil {
			t.Errorf("Parse returned an error: %s", err.Error())
		}
		if got != c.want {
			t.Errorf("Parse(%s) should return %d not %d", c.in, c.want, got)
		}
	}
}

func TestArithPrefix(t *testing.T) {
	type TestCase struct {
		in   string
		want int64
	}
	cases := []TestCase{
		{"~4", -5},
		{"~~4", 4},
		{"!1", A.ShellTrue},
		{"!4", A.ShellTrue},
		{"!0", A.ShellFalse},
		{"!!1", A.ShellFalse},
		{"1+2*3", 7},
		{"1+(2*3)", 7},
		{"(1+2)*3", 9},
	}
	for _, c := range cases {
		got, err := A.Parse(c.in, EmptyScope)
		if err != nil {
			t.Errorf("Parse returned an error: %s", err.Error())
		}
		if got != c.want {
			t.Errorf("Parse(%s) should return %d not %d", c.in, c.want, got)
		}
	}
}

func TestParserTernary(t *testing.T) {
	type TestCase struct {
		in   string
		want int64
	}
	cases := []TestCase{
		{"1 ? 3 : 4", 3},
		{"0 ? 3 : 4", 4},
	}
	for _, c := range cases {
		got, err := A.Parse(c.in, EmptyScope)
		if err != nil {
			t.Errorf("Parse returned an error: %s", err.Error())
		}
		if got != c.want {
			t.Errorf("Parse(%s) should return %d not %d", c.in, c.want, got)
		}
	}
}

func TestParserAssignment(t *testing.T) {
	type TestCase struct {
		inString string
		inVars   map[string]string
		want     int64
		wantVars map[string]string
	}
	cases := []TestCase{
		{
			"x=2",
			map[string]string{},
			2,
			map[string]string{"x": "2"},
		},
		{
			"x+=2",
			map[string]string{},
			2,
			map[string]string{"x": "2"},
		},
		{
			"x+=2",
			map[string]string{"x": "2"},
			4,
			map[string]string{"x": "4"},
		},
		{
			"x*=4",
			map[string]string{"x": "2"},
			8,
			map[string]string{"x": "8"},
		},
		{
			"0 ? x += 2 : x += 2",
			map[string]string{},
			2,
			map[string]string{"x": "2"},
		},
	}

	for _, c := range cases {
		scp := variables.NewScope()
		for k, v := range c.inVars {
			scp.Set(k, v)
		}
		got, err := A.Parse(c.inString, scp)
		if err != nil {
			t.Errorf("Parse returned an error: %s", err.Error())
		}
		if got != c.want {
			t.Errorf("Variable assignment '%s' should evaluate to '%d'", c.inString, c.want)
		}
		for varName, wantVar := range c.wantVars {
			gotVar := scp.Get(varName)
			if gotVar.Val != wantVar {
				t.Errorf("Variable assignment should modify global scope. '%s' should be '%s' not '%s'", varName, wantVar, gotVar.Val)
			}
		}
	}
}

func TestParseError(t *testing.T) {
	cases := []struct {
		in   string
		want A.ParseError
	}{
		{"1*=1", A.ParseError{Fallback: "LHS of assignment 'ArithAssignMultiply' is not a variable"}},
	}

	for _, c := range cases {
		_, gotErr := A.Parse(c.in, EmptyScope)
		if gotErr != c.want {
			t.Errorf("Parse should return the error\n%s\nfor the input '%s' not\n%s", c.want, c.in, gotErr)
		}
	}
}
