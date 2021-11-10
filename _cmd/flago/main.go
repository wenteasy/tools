//
// The flago command is a command that outputs a constant for flag management with one variable.
//
// example:

//    flago var_gen.go pkg Var A B C D
//
// Use:
//    v := pkg.Var(pkg.VarA|pkg.VarC)
//    fmt.Println(v)                                    -> A|C (e.g. -s
//    fmt.Println(v.On(pkg.VarA))                       -> true
//    fmt.Println(v.On(pkg.VarB))                       -> false
//    fmt.Println(v.On(pkg.VarC))                       -> true
//    fmt.Println(v.On(pkg.VarD))                       -> false
//    fmt.Println(v.Equals(pkg.VarA|pkg.VarC))          -> true
//    fmt.Println(v.Equals(pkg.VarA|pkg.VarC|pkg.VarB)) -> false
//    fmt.Println(v.Equals(5))                          -> true
//    fmt.Println(v.Equals(7))                          -> false
//
// -o Omit the prefix(default:false)
//       VarA -> A
//
// -s Split charactor(default:"|")
//       A|C -> A-C
//
//
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"golang.org/x/xerrors"
)

var fv flagsValue

type flagsValue struct {
	Mark          string
	RunCmd        string
	Package       string
	TypeName      string
	UpperTypeName string
	OriginalType  string
	Splitter      string
	OmitPrefix    bool
	Mod           int
	Values        []value
}

type value struct {
	Name  string
	Value string
}

func init() {
	flag.BoolVar(&fv.OmitPrefix, "o", false, "Omit the prefix")
	flag.StringVar(&fv.Splitter, "s", "|", "Split Value")
	flag.StringVar(&fv.OriginalType, "d", "uint", "Define Type(uint,int,uint8,int8,uint16,int16,uint32,uint64)")
	flag.IntVar(&fv.Mod, "m", 3, "values line break ratio")
}

func main() {

	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "run() error:\n%+v\n", err)
		os.Exit(1)
	}
}

func checkArgs() ([]string, error) {
	flag.Parse()
	args := flag.Args()

	if len(args) < 4 {
		return nil, fmt.Errorf("requires 4 or more arguments[FileName,PackageName,TypeName,Value1,Value2,,,ValueN]")
	}

	return args, nil
}

func run() error {

	args, err := checkArgs()
	if err != nil {
		return xerrors.Errorf("checkArgs() error: %w", err)
	}

	fn := args[0]

	pn := args[1]
	fv.Package = pn

	tn := args[2]
	fv.TypeName = tn
	fv.UpperTypeName = upperCase(tn)

	fv.Mark = fmt.Sprintf("Auto Generated: %v", time.Now())
	fv.RunCmd = fmt.Sprintf("Run Command   : %v", os.Args)

	values := make([]value, 0, len(args)-3)
	for idx := 3; idx < len(args); idx++ {
		v := value{}
		n := args[idx]
		v.Name = n
		if !fv.OmitPrefix {
			v.Name = tn + n
		}
		v.Value = n
		values = append(values, v)
	}
	fv.Values = values

	err = generate(fn)
	if err != nil {
		return xerrors.Errorf("generate() error: %w", err)
	}

	fmt.Println("Generated:", fn)

	err = gofmt(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gofmt() error: %+v", err)
	}

	return nil
}

func generate(fn string) error {

	var err error

	tmpl := template.New("root").Funcs(funcMap())
	tmpl, err = tmpl.Parse(rootTxt)
	if err != nil {
		return xerrors.Errorf("template Parse() error: %w", err)
	}

	f, err := os.Create(fn)
	if err != nil {
		return xerrors.Errorf("os.Create() error: %w", err)
	}
	defer f.Close()

	err = tmpl.Execute(f, fv)
	if err != nil {
		return xerrors.Errorf("template Execute() error: %w", err)
	}
	return nil
}

func gofmt(fn string) error {
	_, err := exec.Command("go", "fmt", fn).Output()
	if err != nil {
		return xerrors.Errorf("exec.Command() error: %w", err)
	}
	fmt.Println("Go Formated.")
	return nil
}

func funcMap() template.FuncMap {
	funcs := map[string]interface{}{
		"mod": mod,
	}
	return funcs
}

func upperCase(v string) string {
	l := v[0]
	U := strings.ToUpper(string(l))
	if len(v) == 1 {
		return U
	}

	U += v[1:]
	return U
}

func mod(v, m int) bool {
	v1 := v + 1
	return (v1 % m) == 0
}

const rootTxt = `{{ define "root" }}
// {{ .Mark }}
// {{ .RunCmd }}
package {{ .Package }}

import (
	"strings"
)

type {{ .TypeName }} {{ .OriginalType }}

const (
    {{ range $i,$v := .Values }} {{ if eq $i 0 }} {{ .Name }} {{ $.TypeName }} = 1 << iota {{ else }} {{ .Name }} {{ end }} {{printf "\n"}} {{ end }}
)

var values{{ .UpperTypeName }} = []{{.TypeName}}{
    {{ range $i,$v := .Values }} {{ .Name }}, {{ if mod $i $.Mod }} {{ printf "\n" }} {{ end }} {{ end }}
}

func (r {{ .TypeName }}) Equals(v {{ .TypeName }}) bool {
	return r == v
}

func (r {{ .TypeName }}) On(v {{ .TypeName }}) bool {
	return (r & v) == v
}

func (r {{ .TypeName }}) String() string {

	vals := make([]string, 0, len(values{{ .UpperTypeName }}))
	for _, v := range values{{ .UpperTypeName }} {
		if r.On(v) {
			vals = append(vals, v.value())
		}
	}
	return strings.Join(vals, "{{.Splitter}}")
}

func (r {{ .TypeName }}) value() string {
	switch r {
    {{ range .Values }} case {{.Name}}: return "{{.Value}}" {{ printf "\n" }} {{ end }}
	}
	return ""
}

{{ end }}`
