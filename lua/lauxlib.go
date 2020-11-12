package lua

//#include <lua.h>
//#include <lauxlib.h>
//#include <lualib.h>
//#include <stdlib.h>
//#include "golua.h"
import "C"
import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

type LuaError struct {
	code       int
	message    string
	stackTrace []LuaStackEntry
}

func (err *LuaError) Error() string {
	return err.message
}

func (err *LuaError) Code() int {
	return err.code
}

func (err *LuaError) StackTrace() []LuaStackEntry {
	return err.stackTrace
}

// luaL_argcheck
// WARNING: before b30b2c62c6712c6683a9d22ff0abfa54c8267863 the function ArgCheck had the opposite behaviour
func (L *State) Argcheck(cond bool, narg int, extramsg string) {
	if !cond {
		Cextramsg := C.CString(extramsg)
		defer C.free(unsafe.Pointer(Cextramsg))
		C.luaL_argerror(L.s, C.int(narg), Cextramsg)
	}
}

// luaL_argerror
func (L *State) ArgError(narg int, extramsg string) int {
	Cextramsg := C.CString(extramsg)
	defer C.free(unsafe.Pointer(Cextramsg))
	return int(C.luaL_argerror(L.s, C.int(narg), Cextramsg))
}

// luaL_callmeta
func (L *State) CallMeta(obj int, e string) int {
	Ce := C.CString(e)
	defer C.free(unsafe.Pointer(Ce))
	return int(C.luaL_callmeta(L.s, C.int(obj), Ce))
}

// luaL_checkany
func (L *State) CheckAny(narg int) {
	C.luaL_checkany(L.s, C.int(narg))
}

// luaL_checkinteger
func (L *State) CheckInteger(narg int) int {
	return int(C.luaL_checkinteger(L.s, C.int(narg)))
}

// luaL_checknumber
func (L *State) CheckNumber(narg int) float64 {
	return float64(C.luaL_checknumber(L.s, C.int(narg)))
}

// luaL_checkstring
func (L *State) CheckString(narg int) string {
	var length C.size_t
	return C.GoString(C.luaL_checklstring(L.s, C.int(narg), &length))
}

// luaL_checkoption
//
// BUG(everyone_involved): not implemented
func (L *State) CheckOption(narg int, def string, lst []string) int {
	//TODO: complication: lst conversion to const char* lst[] from string slice
	return 0
}

// luaL_checktype
func (L *State) CheckType(narg int, t LuaValType) {
	C.luaL_checktype(L.s, C.int(narg), C.int(t))
}

// luaL_checkudata
func (L *State) CheckUdata(narg int, tname string) unsafe.Pointer {
	Ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(Ctname))
	return unsafe.Pointer(C.luaL_checkudata(L.s, C.int(narg), Ctname))
}

// UdataToBytes UdataToBytes
func (L *State) UdataToBytes(narg int, nsize int) []byte {
	ud := L.ToUserdata(narg)
	sz := L.ToInteger(nsize)
	buf := C.GoBytes(ud, C.int(sz))
	return buf
}

// Executes file, returns nil for no errors or the lua error string on failure
func (L *State) DoFile(filename string) error {
	if r := L.LoadFile(filename); r != 0 {
		return &LuaError{r, L.ToString(-1), L.StackTrace()}
	}
	return L.Call(0, LUA_MULTRET)
}

// Executes the string, returns nil for no errors or the lua error string on failure
func (L *State) DoString(str string) error {
	if r := L.LoadString(str); r != 0 {
		return &LuaError{r, L.ToString(-1), L.StackTrace()}
	}
	return L.Call(0, LUA_MULTRET)
}

// Like DoString but panics on error
func (L *State) MustDoString(str string) {
	if err := L.DoString(str); err != nil {
		panic(err)
	}
}

// luaL_getmetafield
func (L *State) GetMetaField(obj int, e string) bool {
	Ce := C.CString(e)
	defer C.free(unsafe.Pointer(Ce))
	return C.luaL_getmetafield(L.s, C.int(obj), Ce) != 0
}

// luaL_getmetatable
func (L *State) LGetMetaTable(tname string) {
	Ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(Ctname))
	C.lua_getfield(L.s, LUA_REGISTRYINDEX, Ctname)
}

// luaL_gsub
func (L *State) GSub(s string, p string, r string) string {
	Cs := C.CString(s)
	Cp := C.CString(p)
	Cr := C.CString(r)
	defer func() {
		C.free(unsafe.Pointer(Cs))
		C.free(unsafe.Pointer(Cp))
		C.free(unsafe.Pointer(Cr))
	}()

	return C.GoString(C.luaL_gsub(L.s, Cs, Cp, Cr))
}

// luaL_loadfile
func (L *State) LoadFile(filename string) int {
	Cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(Cfilename))
	return int(lualLoadFile(L.s, Cfilename))
}

// luaL_loadstring
func (L *State) LoadString(s string) int {
	Cs := C.CString(s)
	defer C.free(unsafe.Pointer(Cs))
	return int(C.luaL_loadstring(L.s, Cs))
}

// lua_dump
func (L *State) Dump() int {
	ret := int(C.dump_chunk(L.s))
	return ret
}

// lua_load
func (L *State) Load(bs []byte, name string) int {
	chunk := C.CString(string(bs))
	ckname := C.CString(name)
	defer C.free(unsafe.Pointer(chunk))
	defer C.free(unsafe.Pointer(ckname))
	ret := int(C.load_chunk(L.s, chunk, C.int(len(bs)), ckname))
	if ret != 0 {
		return ret
	}
	return 0
}

// luaL_newmetatable
func (L *State) NewMetaTable(tname string) bool {
	Ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(Ctname))
	return C.luaL_newmetatable(L.s, Ctname) != 0
}

// luaL_newstate
func NewState() *State {
	ls := (C.luaL_newstate())
	if ls == nil {
		return nil
	}
	L := newState(ls)
	return L
}

// luaL_openlibs
func (L *State) OpenLibs() {
	C.luaL_openlibs(L.s)
	C.clua_hide_pcall(L.s)
}

// luaL_optinteger
func (L *State) OptInteger(narg int, d int) int {
	return int(C.luaL_optinteger(L.s, C.int(narg), C.lua_Integer(d)))
}

// luaL_optnumber
func (L *State) OptNumber(narg int, d float64) float64 {
	return float64(C.luaL_optnumber(L.s, C.int(narg), C.lua_Number(d)))
}

// luaL_optstring
func (L *State) OptString(narg int, d string) string {
	var length C.size_t
	Cd := C.CString(d)
	defer C.free(unsafe.Pointer(Cd))
	return C.GoString(C.luaL_optlstring(L.s, C.int(narg), Cd, &length))
}

// luaL_ref
func (L *State) Ref(t int) int {
	return int(C.luaL_ref(L.s, C.int(t)))
}

// luaL_typename
func (L *State) LTypename(index int) string {
	return C.GoString(C.lua_typename(L.s, C.lua_type(L.s, C.int(index))))
}

// luaL_unref
func (L *State) Unref(t int, ref int) {
	C.luaL_unref(L.s, C.int(t), C.int(ref))
}

// luaL_where
func (L *State) Where(lvl int) {
	C.luaL_where(L.s, C.int(lvl))
}

// clua_seri_unpack
func (L *State) LSeriUnpack(nargs, nsz int) int {
	return int(C.clua_seri_unpack(L.s,
		C.int(nargs), C.int(nsz)))
}

// clua_seri_pack
func (L *State) LSeriPack() int {
	return int(C.clua_seri_pack(L.s))
}

// LSeriFree LSeriFree
// clua_seri_free
func LSeriFree(ud unsafe.Pointer) {
	C.clua_seri_free(ud)
}

// SetGoFuncs SetGoFuncs
func (L *State) SetGoFuncs(n int,
	funcs map[string]LuaGoFunction) {
	for k, v := range funcs {
		L.PushGoFunction(v)
		L.SetField(n, k)
	}
}

// SetGoGC SetGoGC
func (L *State) SetGoGC(n int, name string) int {
	err := L.DoString(fmt.Sprintf(`
	function g__gc(t) 
		if t.%s and type(t.%s) == 'userdata' then
			t.%s(t)
		end 
	end
	`, name, name, name))
	if err != nil {
		L.RaiseError(err.Error())
	}
	L.GetGlobal("g__gc")
	L.SetField(n, "__gc")
	err = L.DoString(`g__gc = nil`)
	if err != nil {
		L.RaiseError(err.Error())
	}
	return 0
}

// SetMetaFunc SetMetaFunc
// (1, 'call'), will set the meta funct __call
func (L *State) SetMetaFunc(n int, name string) int {
	field := "__" + name
	fname := "g" + field
	err := L.DoString(fmt.Sprintf(`
	function %s(t) 
		if t.%s and type(t.%s) == 'userdata' then
			t.%s(t)
		end 
	end
	`, fname, name, name, name))
	if err != nil {
		L.RaiseError(err.Error())
	}
	L.GetGlobal(field)
	L.SetField(n, field)
	err = L.DoString(fname + ` = nil`)
	if err != nil {
		L.RaiseError(err.Error())
	}
	return 0
}

// LPushError LPushError
func (L *State) LPushError(err error) int {
	return L.LPushErrorStr(err.Error())
}

// LPushErrorStr LPushErrorStr
func (L *State) LPushErrorStr(s string) int {
	L.PushBoolean(false)
	L.PushString(s)
	return 2
}

// LPushToTable LPushToTable
func (L *State) LPushToTable(ptr interface{}) int {
	fields, _ := getLuaField(ptr)
	L.NewTable()
	for name, v := range fields {
		setV(L, v, name)
	}
	return 1
}

func setV(L *State, v reflect.Value, name string) {
	switch v.Kind() {
	case reflect.Ptr:
		if !v.IsNil() {
			setV(L, v.Elem(), name)
		}
	case reflect.String:
		L.PushString(v.String())
		L.SetField(-2, name)
	case reflect.Int:
		L.PushInteger(v.Int())
		L.SetField(-2, name)
	case reflect.Float64:
		L.PushNumber(v.Float())
		L.SetField(-2, name)
	case reflect.Bool:
		L.PushBoolean(v.Bool())
		L.SetField(-2, name)
	}
}

// LGetFromTable LGetFromTable
func (L *State) LGetFromTable(ptr interface{}, n int) {
	fields, defaluts := getLuaField(ptr)
	L.fromTable(ptr, n, fields, defaluts)
}

// LFromTable LFromTable
func (L *State) fromTable(ptr interface{}, n int,
	fields map[string]reflect.Value,
	defaluts map[string]string) {
	for name, v := range fields {
		L.GetField(n, name)
		var val string
		if !L.IsNoneOrNil(-1) {
			val = L.ToString(-1)
		} else {
			if val2, ok := defaluts[name]; ok {
				val = val2
			}
		}
		if val != "" {
			setVal(v, val)
		}
	}
}

// LGetGlobal LGetGlobal
func (L *State) LGetGlobal(ptr interface{}) {
	fields, defaluts := getLuaField(ptr)
	for name, v := range fields {
		L.GetGlobal(name)
		var val string
		if !L.IsNoneOrNil(-1) {
			val = L.ToString(-1)
		} else {
			if val2, ok := defaluts[name]; ok {
				val = val2
			}
		}
		if val != "" {
			setVal(v, val)
		}
	}
}

func getFields(v reflect.Value,
	fields map[string]reflect.Value,
	defaults map[string]string) {
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		if fieldInfo.Anonymous {
			getFields(v.Field(i), fields, defaults)
			continue
		}
		tag := fieldInfo.Tag
		name := tag.Get("lua")
		if strings.Contains(name, ",") {
			names := strings.Split(name, ",")
			if len(names) > 1 {
				name = names[0]
				if name == "" {
					name = strings.ToLower(fieldInfo.Name)
				}
				defaults[name] = names[1]
			}
		}
		if name != "" {
			fields[name] = v.Field(i)
		}
	}
}

func getLuaField(ptr interface{}) (map[string]reflect.Value, map[string]string) {
	fields := make(map[string]reflect.Value)
	defaults := make(map[string]string)
	elem := reflect.ValueOf(ptr).Elem()
	getFields(elem, fields, defaults)
	return fields, defaults
}

func setVal(v reflect.Value, val string) {
	switch v.Kind() {
	case reflect.String:
		v.SetString(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			v.SetInt(i)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if i, err := strconv.ParseUint(val, 10, 64); err == nil {
			v.SetUint(i)
		}
	case reflect.Float32, reflect.Float64:
		if i, err := strconv.ParseFloat(val, 64); err == nil {
			v.SetFloat(i)
		}
	case reflect.Bool:
		if b, err := strconv.ParseBool(val); err == nil {
			v.SetBool(b)
		}
	}
}
