// Package common provides common functions for mock generation.
package common

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/koron-go/srcdom"
)

var ForTest bool = false

type Variable struct {
	Name string
	Typ  string
}

type Vars []*Variable

func (vv *Vars) add(v *Variable) {
	*vv = append(*vv, v)
}

func (vv Vars) NameTypes() string {
	return vv.Join(func(v *Variable) string {
		return v.Name + " " + v.Typ
	})
}

func (vv Vars) Names() string {
	return vv.Join(func(v *Variable) string {
		return v.Name
	})
}

func (vv Vars) NamesPrefix(prefix string) string {
	return vv.Join(func(v *Variable) string {
		return prefix + "." + v.Name
	})
}

func (vv Vars) Types() string {
	return vv.Join(func(v *Variable) string {
		return v.Typ
	})
}

func (vv Vars) Join(fn func(v *Variable) string) string {
	b := &strings.Builder{}
	for i, v := range vv {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(fn(v))
	}
	return b.String()
}

type Method struct {
	Typn string
	Name string
	Args Vars
	Rets Vars
}

func (m *Method) ParamTypeName() string {
	return m.Typn + m.Name + "_P"
}

func (m *Method) ReturnTypeName() string {
	return m.Typn + m.Name + "_R"
}

// varName generates variable name.
func varName(name string, attr string, n int) string {
	if name != "" {
		return name
	}
	return attr + strconv.Itoa(n)
}

// FilterMethods filter methods which match with typname.
func FilterMethods(src []*srcdom.Func, typname string) []*Method {
	var dst []*Method
	for _, f := range src {
		if !f.IsPublic() {
			continue
		}
		m := &Method{Typn: typname, Name: f.Name}
		for i, p := range f.Params {
			m.Args.add(&Variable{
				Name: varName(p.Name, "in", i),
				Typ:  p.Type,
			})
		}
		for i, r := range f.Results {
			m.Rets.add(&Variable{
				Name: varName(r.Name, "Out", i),
				Typ:  r.Type,
			})
		}
		dst = append(dst, m)
	}
	return dst
}

// ToStructFieldType convert arg/parame type name to struct field type name.
func ToStructFieldType(typ string) string {
	if strings.HasPrefix(typ, "...") {
		return "[]" + typ[3:]
	}
	return typ
}

// ToPub convert name as public name.
// It make a first character to upper.
func ToPub(s string) string {
	if s == "" {
		return ""
	}
	_, n := utf8.DecodeRuneInString(s)
	p, r := s[:n], s[n:]
	return strings.ToUpper(p) + r
}
