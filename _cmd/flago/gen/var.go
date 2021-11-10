// Auto Generated: 2021-11-11 06:30:00.8101892 +0900 JST m=+0.002103301
// Run Command   : [C:\Users\secon\GoApp\wenteasy\tools\_cmd\flago\main.exe ./gen/var.go pkg Var A B C D E F G]
package pkg

import (
	"strings"
)

type Var uint

const (
	VarA Var = 1 << iota
	VarB
	VarC
	VarD
	VarE
	VarF
	VarG
)

var valuesVar = []Var{
	VarA, VarB, VarC,
	VarD, VarE, VarF,
	VarG,
}

func (r Var) Values() []Var {
	return valuesVar
}

func (r Var) Equals(v Var) bool {
	return r == v
}

func (r Var) On(v Var) bool {
	return (r & v) == v
}

func (r Var) Sum(v Var) Var {
	return (r | v)
}

func (r Var) String() string {

	vals := make([]string, 0, len(valuesVar))
	for _, v := range valuesVar {
		if r.On(v) {
			vals = append(vals, v.value())
		}
	}
	return strings.Join(vals, "|")
}

func (r Var) value() string {
	switch r {
	case VarA:
		return "A"
	case VarB:
		return "B"
	case VarC:
		return "C"
	case VarD:
		return "D"
	case VarE:
		return "E"
	case VarF:
		return "F"
	case VarG:
		return "G"

	}
	return ""
}
