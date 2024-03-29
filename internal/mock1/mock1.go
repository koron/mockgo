// Package mock1 provides generator of mock version 1
package mock1

import (
	"fmt"
	"io"

	"github.com/koron-go/srcdom"
	"github.com/koron/mockgo/internal/common"
)

// Generate generates a mock (ver.1) for a type.
func Generate(w io.Writer, mockTag, mockTypn, mockPkgn string, typ *srcdom.Type, pkg *srcdom.Package) error {
	origTypn := pkg.Name + "." + typ.Name
	methods := common.FilterMethods(typ.Methods, mockTypn)
	if len(methods) == 0 {
		return fmt.Errorf("no methods in type:%s", typ.Name)
	}

	// write headers.
	if !common.ForTest {
		fmt.Fprintf(w, "//go:build %s\n\n", mockTag)
		fmt.Fprintf(w, "// +build %s\n\n", mockTag)
	}
	fmt.Fprintf(w, "// Code generated by github.com/koron/mockgo; DO NOT EDIT.\n\n")
	fmt.Fprintf(w, "package %s\n\n", mockPkgn)

	// write the mock type.
	fmt.Fprintf(w, "// %s is a mock of %s for test.\n", mockTypn, origTypn)
	fmt.Fprintf(w, "type %s struct {\n", mockTypn)
	for _, m := range methods {
		fmt.Fprintf(w, "\t%s_Ps []*%s\n", m.Name, m.ParamTypeName())
		fmt.Fprintf(w, "\t%s_Rs []*%s\n", m.Name, m.ReturnTypeName())
	}
	fmt.Fprintf(w, "}\n")

	for _, m := range methods {
		fmt.Fprintf(w, "\n")

		// write parameter type for the method.
		fmt.Fprintf(w, "// %s packs input parameters of %s#%s method.\n", m.ParamTypeName(), origTypn, m.Name)
		fmt.Fprintf(w, "type %s struct {\n", m.ParamTypeName())
		for _, a := range m.Args {
			typ := common.ToStructFieldType(a.Typ)
			fmt.Fprintf(w, "\t%s %s\n", common.ToPub(a.Name), typ)
		}
		fmt.Fprintf(w, "}\n\n")

		// write result type for the method.
		fmt.Fprintf(w, "// %s packs output parameters of %s#%s method.\n", m.ReturnTypeName(), origTypn, m.Name)
		fmt.Fprintf(w, "type %s struct {\n", m.ReturnTypeName())
		for _, r := range m.Rets {
			fmt.Fprintf(w, "\t%s %s\n", r.Name, r.Typ)
		}
		fmt.Fprintf(w, "}\n\n")

		// write mock func for the method.
		fmt.Fprintf(w, "// %s is mock of %s#%[1]s method.\n", m.Name, origTypn)
		fmt.Fprintf(w, "func (_m *%s) %s(%s) (%s) {\n", mockTypn, m.Name, m.Args.NameTypes(), m.Rets.Types())
		fmt.Fprintf(w, "\t_m.%s_Ps = append(_m.%[1]s_Ps, &%s{%s})\n", m.Name, m.ParamTypeName(), m.Args.Names())
		fmt.Fprintf(w, "\tvar _r *%s\n", m.ReturnTypeName())
		fmt.Fprintf(w, "\t_r, _m.%[1]s_Rs = _m.%[1]s_Rs[0], _m.%[1]s_Rs[1:]\n", m.Name)
		fmt.Fprintf(w, "\treturn %s\n", m.Rets.NamesPrefix("_r"))
		fmt.Fprintf(w, "}\n")
	}
	return nil
}
