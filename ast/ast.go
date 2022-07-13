package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type Model struct {
	Name            string
	Package         string
	Fields          []Field
	parsedStruct    bool
	parsedOverrides bool
}

func Parse(filename string) *Model {
	fset := token.NewFileSet()
	//f, err := parser.ParseFile(fset, filename, nil, parser.Trace)
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		fmt.Printf("Oops! Can't parse the source: %v\n", err)
		return nil
	}

	model := &Model{}
	ast.Walk(model, f)

	fmt.Printf("MODEL: %s\n", model.Name)
	for idx, f := range model.Fields {
		fmt.Printf("  %d - %#v\n", idx, f)
		//spew.Dump(f)
	}

	return model
}

func (v *Model) Visit(node ast.Node) (w ast.Visitor) {
	if v.parsedStruct {
		return v
	}

	switch t := node.(type) {
	case *ast.File:
		v.Package = t.Name.Name

	case *ast.TypeSpec:
		str, ok1 := t.Type.(*ast.StructType)

		if ok1 {
			v.Name = t.Name.Name

			for _, inp := range str.Fields.List {
				var out Field
				if inp.Type != nil {
					typ, ok := (inp.Type).(*ast.Ident)
					if ok {
						out.Type = typ.Name
					} else {
						v, ok := (inp.Type).(*ast.StarExpr)
						if ok {
							out.Type = exprName(v)
							//							fmt.Printf("X2 : %#v\n", exprName(v))
						}

					}

				}
				if inp.Tag != nil {
					out.Tag = strings.Replace(inp.Tag.Value, "`", "", -1)
				}
				if len(inp.Names) == 1 {
					out.Name = inp.Names[0].Name
				} else {
					panic("Couldn't find field name")
				}
				//				fmt.Printf("   FIELD: %#v \n", out)

				v.Fields = append(v.Fields, out)
			}
			v.parsedStruct = true
		}
	}

	return v
}

func exprName(e ast.Expr) string {
	switch e.(type) {
	case *ast.Ident:
		return e.(*ast.Ident).Name
	case *ast.StarExpr:
		return exprName(e.(*ast.StarExpr).X)
	case *ast.ArrayType:
		return exprName(e.(*ast.ArrayType).Elt)
	case *ast.MapType:
		mt := e.(*ast.MapType)
		return fmt.Sprintf("map[%s]%s", exprName(mt.Key), exprName(mt.Value))
	case *ast.SelectorExpr:
		s := e.(*ast.SelectorExpr)
		return fmt.Sprintf("%s.%s", exprName(s.X), s.Sel.Name)
	case *ast.ChanType:
		ch := e.(*ast.ChanType)
		var chtype string
		if ch.Dir == ast.SEND {
			chtype = "chan<-"
		} else if ch.Dir == ast.RECV {
			chtype = "<-chan"
		} else {
			chtype = "chan"
		}
		return fmt.Sprintf("%s %s", chtype, exprName(ch.Value))
	case *ast.StructType:
		return "struct{} (unknown name)"
	case *ast.InterfaceType:
		return "interface (unknown name)"
	default:
		return fmt.Sprintf("unhandled expr (%T)", e)
	}
}

func (m *Model) FieldSlice() []string {
	out := []string{}
	for _, v := range m.Fields {
		out = append(out, v.Name)
	}
	return out
}

func (m *Model) FieldSliceWithoutID() []string {
	out := []string{}
	for _, v := range m.Fields {
		if v.Name != "ID" {
			out = append(out, v.Name)
		}
	}
	return out
}

func (m *Model) parseOverrides() {
	for _, f := range m.Fields {
		f.parseOverrides()
	}
}
