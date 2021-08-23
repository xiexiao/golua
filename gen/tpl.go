package gen

import (
	_ "embed"
	"strings"
	"text/template"
)

//TODO: tpl use go fmt

func join(s []string) string {
	return strings.Join(s, ", ")
}

//go:embed tpl_func.txt
var tplFunc string
var funcTemplate = template.Must(template.New("func").
	Funcs(template.FuncMap{"join": join}).
	Parse(tplFunc))

//go:embed tpl_func_new.txt
var tplFuncNew string
var funcNewTemplate = template.Must(funcTemplate.New("func_new").
	Parse(tplFuncNew))

//go:embed tpl_func_field.txt
var tplFuncField string
var funcFieldTemplate = template.Must(funcTemplate.New("func_field").
	Parse(tplFuncField))

//go:embed tpl_code.txt
var tplCode string
var codeTemplate = template.Must(funcTemplate.New("code").
	Parse(tplCode))

type Code struct {
	PkgName    string
	StructName string
	ModuleName string
	Imports    []PkgImport

	Fields []CodeField

	FuncNew *CodeFunc

	Funcs []CodeFunc
}

type PkgImport struct {
	Name   string
	Import string
}

type CodeFunc struct {
	Name       string
	StructID   string
	StructName string
	FuncName   string

	ResultList []string
	ParamList  []string

	ErrCheck string

	Params  []CodeItem
	Results []CodeItem
}

type CodeItem struct {
	Name string
	Code string
}

type CodeField struct {
	StructName string
	Name       string
	Code       string
}
