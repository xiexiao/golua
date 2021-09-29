package gen

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"
)

//TODO: tpl use go fmt

func join(s []string) string {
	return strings.Join(s, ", ")
}

func wrap(its ...interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i, it := range its {
		m[fmt.Sprintf("T%d", i)] = it
	}
	return m
}

//go:embed tpl_func.txt
var tplFunc string
var funcTemplate = template.Must(template.New("func").
	Funcs(template.FuncMap{"join": join}).
	Funcs(template.FuncMap{"wrap": wrap}).
	Parse(tplFunc))

/*
//go:embed tpl_func_new.txt
var tplFuncNew string
var funcNewTemplate = template.Must(funcTemplate.New("func_new").
	Parse(tplFuncNew))
*/

//go:embed tpl_func_field.txt
var tplFuncField string
var funcFieldTemplate = template.Must(funcTemplate.New("func_field").
	Parse(tplFuncField))

//go:embed tpl_code.txt
var tplCode string
var codeTemplate = template.Must(funcTemplate.New("code").
	Parse(tplCode))

type Code struct {
	FileName     string
	CodeHeader   string
	PkgName      string
	StructID     string //use id in code
	StructName   string
	StructPrefix string
	ModuleName   string
	Imports      []PkgImport

	Funcs []CodeFunc

	Fields []CodeField

	TypeFuncs []CodeFunc
}

type PkgImport struct {
	Name   string
	Import string
}

type CodeFunc struct {
	IsType  bool
	StartID string

	FuncID     string
	FuncName   string
	FuncPrefix string

	ResultList []string
	ParamList  []string

	Params  []CodeItem
	Results []CodeItem
}

type CodeItem struct {
	Name string
	Code string
}

type CodeField struct {
	FieldName string
	Name      string
	Code      string
}
