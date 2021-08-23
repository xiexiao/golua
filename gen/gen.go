package gen

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"text/template"
)

func GenLua(filename string, in []byte) ([]byte, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, in, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	code := &Code{}
	code.PkgName = f.Name.Name
	code.ModuleName = strings.ToLower(code.PkgName)
	structName := ""
	var fields []CodeField

	for name, obj := range f.Scope.Objects {
		spec, ok := obj.Decl.(*ast.TypeSpec)
		if ok {
			stype, ok := spec.Type.(*ast.StructType)
			if ok {
				structName = name
				for i, x := range stype.Fields.List {
					if x.Tag != nil && len(x.Names) > 0 && strings.HasPrefix(x.Tag.Value, "`lua") {
						ftype := x.Type.(*ast.Ident).Name
						fname := x.Names[0].Name
						name1 := ParseTag(x.Tag.Value)
						if name1 != "" {
							fname = name1
						}
						vname := fmt.Sprintf("v0.%s", fname)
						fcode := GenPushCode(ftype, vname, 0)
						if fcode != "" {
							fields = append(fields, CodeField{
								StructName: structName,
								Name:       fname,
								Code:       fcode,
							})
						}
					}
				}
				break
			}
		}
	}

	if len(fields) > 0 {
		code.Fields = fields
	}
	code.StructName = structName

	var funcNew *CodeFunc
	var funcs []CodeFunc
	for _, d := range f.Decls {
		fd, ok := d.(*ast.FuncDecl)
		checkError := false
		if ok {
			funcName := fd.Name.Name
			name := funcName
			if fd.Doc != nil {
				funcName1, checkError1 := ParseDoc(fd.Doc)
				if funcName1 != "" {
					name = funcName1
				}
				if checkError1 {
					checkError = checkError1
				}
			}
			paramList, resultList, params,
				results, structID, errCheck := ParseFunc(fd, checkError)
			if fd.Recv != nil {
				rtype := fd.Recv.List[0].Type
				if rt1, ok := rtype.(*ast.StarExpr); ok {
					tname := rt1.X.(*ast.Ident).Name // Type Name
					if tname == structName {         // Recv Must be Current Type
						cf := CodeFunc{
							Name:       name,
							FuncName:   funcName,
							StructID:   structID,
							StructName: structName,
							ErrCheck:   errCheck,
							ParamList:  paramList,
							ResultList: resultList,
							Params:     params,
							Results:    results,
						}
						funcs = append(funcs, cf)
					}
				}
			} else {
				if funcName == structName+"New" {
					funcNew = &CodeFunc{
						StructName: structName,
						ErrCheck:   errCheck,
						ParamList:  paramList,
						ResultList: resultList,
						Params:     params,
						Results:    results,
					}
					code.FuncNew = funcNew
				}
			}
		}
	}
	if len(funcs) > 0 {
		code.Funcs = funcs
	}
	return codeExec(codeTemplate, code)
}

func codeExec(t *template.Template, data interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
