package gen

import (
	"fmt"
	"go/ast"
	"log"
	"strings"
)

//TODO: Add Imports
//TODO: []string impl

func ParseTag(tag string) string {
	if strings.HasPrefix(tag, "`lua:") {
		str := tag[5 : len(tag)-1]
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

const makeUdataTag = "__make_udata" //返回用户定义类型，绑定到特定的metatable上的

func ParseDoc(cg *ast.CommentGroup) (name string, makeUdata bool) {
	if cg != nil && len(cg.List) > 0 {
		text := cg.List[0].Text
		if strings.HasPrefix(text, "//lua:") {
			str := text[6:]
			list := strings.Split(str, ",")
			if len(list) > 0 && list[0] != makeUdataTag {
				name = list[0]
			}
			makeUdata = strings.Contains(str, makeUdataTag)
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
	case "[]byte":
		format = "%s := L.ToBytes(%d)"
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
		if strings.HasPrefix(vtype, "[]") {
			format = "%s := %s; L.ToArray(%s, %d)"
			code = fmt.Sprintf(format, vname, vtype, vname, idx)
		} else {
			format = "%s := &%s{}; L.LGetFromTable(%s, %d)"
			code = fmt.Sprintf(format, vname, vtype, vname, idx)
		}
		break
	}
	if format != "" && code == "" {
		code = fmt.Sprintf(format, vname, idx)
	}
	return code
}

func GenPushCode(vtype, vname string, idx int, makeUdata bool) string {
	format := ""
	code := ""
	switch vtype {
	case "byte":
		format = "L.PushInteger(int64(%s))"
	case "[]byte":
		format = "L.PushBytes(%s)"
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
		if makeUdata {
			return fmt.Sprintf("L.NewUdataWithField(%s, luaMT_%s)", vname, vtype)
		}
		format = "L.LPushToTable(%s)"
		break
	}
	if format != "" {
		code = fmt.Sprintf(format, vname)
	}
	return code
}

func ParseFunc(fd *ast.FuncDecl, idx int, makeUdata bool) ([]string, []string,
	[]CodeItem, []CodeItem, string) {
	var paramList []string
	var resultList []string
	var params []CodeItem
	var results []CodeItem

	ftype := fd.Type
	//paramNum := len(ftype.Params.List)
	for i, x := range ftype.Params.List {
		ptype := ""
		switch tt := x.Type.(type) {
		case *ast.Ident:
			ptype = tt.Name
		case *ast.StarExpr:
			if t2, ok2 := tt.X.(*ast.Ident); ok2 {
				ptype = t2.Name
			}
		case *ast.ArrayType:
			if id2, ok := tt.Elt.(*ast.Ident); ok {
				ptype = "[]" + id2.Name
			} else if see, ok := tt.Elt.(*ast.SelectorExpr); ok {
				i1, ok1 := see.X.(*ast.Ident)
				if ok1 {
					ptype = i1.Name + "." + see.Sel.Name
				}
			}
		default:
			log.Printf("-->x.Type:%t,%+v\n", tt, tt)
			break
		}
		if ptype != "" {
			pname := x.Names[0].Name
			pcode := GenGetCode(ptype, pname, i+idx)
			paramList = append(paramList, pname)
			if pcode != "" {
				params = append(params, CodeItem{
					Name: pname,
					Code: pcode,
				})
			}
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
			switch tt := x.Type.(type) {
			case *ast.Ident:
				rtype = tt.Name
			case *ast.StarExpr:
				if ide, ok := tt.X.(*ast.Ident); ok {
					rtype = ide.Name
				} else if see, ok := tt.X.(*ast.SelectorExpr); ok {
					i1, ok1 := see.X.(*ast.Ident)
					if ok1 {
						rtype = i1.Name + "." + see.Sel.Name
					}
				}
			case *ast.ArrayType:
				if id2, ok := tt.Elt.(*ast.Ident); ok {
					rtype = "[]" + id2.Name
				} else if see, ok := tt.Elt.(*ast.SelectorExpr); ok {
					i1, ok1 := see.X.(*ast.Ident)
					if ok1 {
						rtype = i1.Name + "." + see.Sel.Name
					}
				}
			default:
				log.Printf("-->x.Type:%t,%+v\n", tt, tt)
				break
			}
			if rtype != "" {
				rname := fmt.Sprintf("%s%d", start, ri)
				resultList = append(resultList, rname)
				rcode := GenPushCode(rtype, rname, ri, makeUdata)
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

	return paramList, resultList, params, results, structID
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
