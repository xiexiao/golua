package gen

import (
	"fmt"
	"go/ast"
	"strings"
)

//TODO: Add Imports
//TODO: []string impl

func ParseTag(tag string) string {
	if strings.HasPrefix(tag, "`lua:") {
		str := tag[5:]
		list := strings.Split(str, ",")
		if len(list) > 0 {
			return list[0]
		}
	}
	return ""
}

/*
func AddImports(names []string, imports []PkgImport, name string) {
}
*/

const checkErrorTag = "check_error"

func ParseDoc(cg *ast.CommentGroup) (name string, checkError bool) {
	if cg != nil && len(cg.List) > 0 {
		text := cg.List[0].Text
		if strings.HasPrefix(text, "//lua:") {
			str := text[6:]
			list := strings.Split(str, ",")
			if len(list) > 0 && list[0] != checkErrorTag {
				name = list[0]
			}
			checkError = strings.Contains(str, checkErrorTag)
		}
	}
	return
}

func GenGetCode(vtype, vname string, idx int) string {
	format := ""
	code := ""
	switch vtype {
	case "byte":
		format = "%s := byte(L.CheckInteger(%d))"
	case "int":
		format = "%s := L.CheckInteger(%d)"
	case "uint16":
		format = "%s := uint16(L.CheckInteger(%d))"
	case "int16":
		format = "%s := int16(L.CheckInteger(%d))"
	case "int32":
		format = "%s := int32(L.CheckInteger(%d))"
	case "uint32":
		format = "%s := uint32(L.CheckInteger(%d))"
	case "int64":
		format = "%s := int64(L.CheckInteger(%d))"
	case "uint64":
		format = "%s := uint64(L.CheckInteger(%d))"
	case "float32":
		format = "%s := float32(L.CheckInteger(%d))"
	case "float64":
		format = "%s := float64(L.CheckInteger(%d))"
	case "string":
		format = "%s := L.CheckString(%d)"
	case "bool":
		format = "%s := L.ToBoolean(%d)"
	default:
		format = "%s := &%s{};L.LGetFromTable(%s, %d)"
		code = fmt.Sprintf(format, vname, vtype, vname, idx)
		break
	}
	if format != "" && code == "" {
		code = fmt.Sprintf(format, vname, idx)
	}
	return code
}

func GenPushCode(vtype, vname string, idx int) string {
	format := ""
	code := ""
	switch vtype {
	case "byte":
		format = "L.PushInteger(int64(%s))"
	case "int":
		format = "L.PushInteger(int64(%s))"
	case "uint16", "int16":
		format = "L.PushInteger(int64(%s))"
	case "int32":
		format = "L.PushInteger(int64(%s))"
	case "uint32":
		format = "L.PushInteger(int64(%s))"
	case "int64":
		format = "L.PushInteger(%s)"
	case "uint64":
		format = "L.PushInteger(int64(%s))"
	case "float32":
		format = "L.PushNumber(float64(%s))"
	case "float64":
		format = "L.PushNumber(float64(%s))"
	case "string":
		format = "L.PushString(%s)"
	case "bool":
		format = "L.PushBoolean(%s)"
	case "error":
		format = "L.LPushError(%s)"
	default:
		format = "L.LPushToTable(%s)"
		break
	}
	if format != "" {
		code = fmt.Sprintf(format, vname)
	}
	return code
}

func ParseFunc(fd *ast.FuncDecl, check bool) ([]string, []string,
	[]CodeItem, []CodeItem, string, string) {
	var paramList []string
	var resultList []string
	var params []CodeItem
	var results []CodeItem
	errCheck := ""

	ftype := fd.Type
	//paramNum := len(ftype.Params.List)
	for i, x := range ftype.Params.List {
		ptype := ""
		if t, ok := x.Type.(*ast.Ident); ok {
			ptype = t.Name
		} else if t, ok := x.Type.(*ast.StarExpr); ok {
			if t2, ok2 := t.X.(*ast.Ident); ok2 {
				ptype = t2.Name
			}
		}
		pname := x.Names[0].Name
		pcode := GenGetCode(ptype, pname, i+1)
		paramList = append(paramList, pname)
		if pcode != "" {
			params = append(params, CodeItem{
				Name: pname,
				Code: pcode,
			})
		}
	}

	//resultNum := len(ftype.Results.List)
	start := "v0"
	j := 1
	for {
		got := true
		for _, p := range paramList {
			if strings.HasPrefix(p, start) {
				got = false
			}
		}
		if got {
			break
		}
		start = fmt.Sprintf("v%d", j)
		j = j + 1
	}
	structID := start

	ri := 0
	if ftype.Results != nil {
		for _, x := range ftype.Results.List {
			rtype := ""
			if id, ok := x.Type.(*ast.Ident); ok {
				rtype = id.Name
			} else if se, ok := x.Type.(*ast.StarExpr); ok {
				if ide, ok := se.X.(*ast.Ident); ok {
					rtype = ide.Name
				} else if see, ok := se.X.(*ast.SelectorExpr); ok {
					i1, ok1 := see.X.(*ast.Ident)
					if ok1 {
						rtype = i1.Name + "." + see.Sel.Name
						//fmt.Printf("%+v, %+v\n",see.X, see.Sel)
						//fmt.Printf("%T, %T\n",see.X, see.Sel)
					}
				}
			} else if at, ok := x.Type.(*ast.ArrayType); ok {
				if id2, ok := at.Elt.(*ast.Ident); ok {
					rtype = "[]" + id2.Name
				}
			}
			if rtype != "" {
				rname := fmt.Sprintf("%s%d", start, ri)
				resultList = append(resultList, rname)
				rcode := GenPushCode(rtype, rname, ri)
				if rcode != "" {
					results = append(results, CodeItem{
						Name: rname,
						Code: rcode,
					})
				}
				ri = ri + 1
			}
		}
	}

	return paramList, resultList, params, results, structID, errCheck
}

/*
func GenCheckCode(stype, name string, idx int) string {
	switch stype {
	case "int":
	case "uint16":
	case "int16":
	case "int32":
	case "uint32":
	case "int64":
	case "uint64":
	case "float32":
	case "float64":
	case "string":
	case "boolean":
	default:
		break
	}
	//CheckString
	 // CheckInteger
	  //CheckNumber
	  //ToBoolean()
	return ""
}
*/
