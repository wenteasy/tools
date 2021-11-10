package pkg_test

import (
	"fmt"
	"testing"

	"pkg"
)

func TestVar_Values(t *testing.T) {
	v := pkg.Var(0)
	if len(v.Values()) != 7 {
		t.Errorf("pkg.Var Values() length:want %d, got %d", 7, len(v.Values()))
	}
}

func TestVar_String(t *testing.T) {
	v := pkg.Var(0)
	vars := v.Values()
	wants := []string{"A", "B", "C", "D", "E", "F", "G"}
	for idx, v := range vars {
		want := wants[idx]
		got := v.String()
		if want != got {
			t.Errorf("pkg.Var String():want %s, got %s", want, got)
		}
	}
}

func TestVar_On(t *testing.T) {

	v := pkg.Var(pkg.VarA | pkg.VarB | pkg.VarE | pkg.VarG)
	vals := v.Values()
	wants := []bool{true, true, false, false, true, false, true}
	for idx, val := range vals {
		got := v.On(val)
		want := wants[idx]
		if got != want {
			t.Errorf("pkg.Var On():want %t, got %t", want, got)
		}
	}

}

func TestVar_Sum(t *testing.T) {
	v := pkg.Var(pkg.VarA)
	if !v.On(pkg.VarA) {
		t.Errorf("VarA on VarA")
	}
	if v.On(pkg.VarB) {
		t.Errorf("VarA off VarB")
	}
	if v.On(pkg.VarC) {
		t.Errorf("VarA off VarC")
	}

	v = v.Sum(pkg.VarB)
	if !v.On(pkg.VarA) {
		t.Errorf("pkg.Var Sum(VarB) VarA|B on VarA")
	}
	if !v.On(pkg.VarB) {
		t.Errorf("pkg.Var Sum(VarB) VarA|B on VarB")
	}
	if v.On(pkg.VarC) {
		t.Errorf("pkg.Var Sum(VarB) VarA|B off VarC")
	}

}

func TestVar_Equals(t *testing.T) {
	v1 := pkg.Var(pkg.VarA | pkg.VarB | pkg.VarE | pkg.VarG)
	v2 := pkg.Var(pkg.VarA | pkg.VarB | pkg.VarE | pkg.VarF | pkg.VarG)
	v3 := pkg.Var(pkg.VarA | pkg.VarB | pkg.VarG)
	v4 := v3.Sum(pkg.VarE)

	if v1.Equals(v2) {
		t.Errorf("Var Equals() %v not equals %v", v1, v2)
	}
	if v1.Equals(v3) {
		t.Errorf("Var Equals() %v not equals %v", v1, v3)
	}
	if !v1.Equals(v4) {
		t.Errorf("Var Equals() %v equals %v", v1, v4)
	}
}

func ExampleVar() {

	v := pkg.Var(pkg.VarA | pkg.VarE | pkg.VarG)
	fmt.Println(v)
	v = pkg.Var(pkg.VarB | pkg.VarF | pkg.VarD)
	fmt.Println(v)
	v = pkg.Var(pkg.VarC)
	fmt.Println(v)

	// Output:
	// A|E|G
	// B|D|F
	// C
	//
}
