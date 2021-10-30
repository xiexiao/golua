//go:build lua54
// +build lua54

package lua

/*
#include <lua.h>
#include <lauxlib.h>
#include <lualib.h>
#include <stdlib.h>

typedef struct _chunk {
	int size; // chunk size
	char *buffer; // chunk data
	char* toread; // chunk to read
} chunk;

LUA_API void *lua_newuserdata (lua_State *L, size_t size) {
    return lua_newuserdatauv(L, size, 1);
}

LUA_API int (lua_gc_compat) (lua_State *L, int what, int data) {
    return lua_gc(L, what, data);
}

static const char * reader (lua_State *L, void *ud, size_t *sz) {
	chunk *ck = (chunk *)ud;
	if (ck->size > LUAL_BUFFERSIZE) {
		ck->size -= LUAL_BUFFERSIZE;
		*sz = LUAL_BUFFERSIZE;
		ck->toread = ck->buffer;
		ck->buffer += LUAL_BUFFERSIZE;
	}else{
		*sz = ck->size;
		ck->toread = ck->buffer;
		ck->size = 0;
	}
	return ck->toread;
}

static int writer (lua_State *L, const void* b, size_t size, void* B) {
	static int count=0;
	(void)L;
	luaL_addlstring((luaL_Buffer*) B, (const char *)b, size);
	return 0;
}

// load function chunk dumped from dump_chunk
int load_chunk(lua_State *L, char *b, int size, const char* chunk_name) {
	chunk ck;
	ck.buffer = b;
	ck.size = size;
	int err;
	err = lua_load(L, reader, &ck, chunk_name, NULL);
	if (err != 0) {
		return luaL_error(L, "unable to load chunk, err: %d", err);
	}
	return 0;
}

void clua_openio(lua_State* L)
{
	luaL_requiref(L, "io", &luaopen_io, 1);
	lua_pop(L, 1);
}

void clua_openmath(lua_State* L)
{
	luaL_requiref(L, "math", &luaopen_math, 1);
	lua_pop(L, 1);
}

void clua_openpackage(lua_State* L)
{
	luaL_requiref(L, "package", &luaopen_package, 1);
	lua_pop(L, 1);
}

void clua_openstring(lua_State* L)
{
	luaL_requiref(L, "string", &luaopen_string, 1);
	lua_pop(L, 1);
}

void clua_opentable(lua_State* L)
{
	luaL_requiref(L, "table", &luaopen_table, 1);
	lua_pop(L, 1);
}

void clua_openos(lua_State* L)
{
	luaL_requiref(L, "os", &luaopen_os, 1);
	lua_pop(L, 1);
}

void clua_opencoroutine(lua_State *L)
{
	luaL_requiref(L, "coroutine", &luaopen_coroutine, 1);
	lua_pop(L, 1);
}

void clua_opendebug(lua_State *L)
{
	luaL_requiref(L, "debug", &luaopen_debug, 1);
	lua_pop(L, 1);
}

// dump function chunk from luaL_loadstring
int dump_chunk (lua_State *L) {
	luaL_Buffer b;
	luaL_checktype(L, -1, LUA_TFUNCTION);
	lua_settop(L, -1);
	luaL_buffinit(L,&b);
	int err;
	err = lua_dump(L, writer, &b, 0);
	if (err != 0){
	return luaL_error(L, "unable to dump given function, err:%d", err);
	}
	luaL_pushresult(&b);
	return 0;
}

int clua_seri_unpack(lua_State *L, int n, int sz);
int clua_seri_pack(lua_State *L);
void clua_seri_free(void* ud);
*/
import "C"

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func luaToInteger(s *C.lua_State, n C.int) C.longlong {
	return C.lua_tointegerx(s, n, nil)
}

func luaToNumber(s *C.lua_State, n C.int) C.double {
	return C.lua_tonumberx(s, n, nil)
}

func lualLoadFile(s *C.lua_State, filename *C.char) C.int {
	return C.luaL_loadfilex(s, filename, nil)
}

// lua_equal
func (L *State) Equal(index1, index2 int) bool {
	return C.lua_compare(L.s, C.int(index1), C.int(index2), C.LUA_OPEQ) == 1
}

// lua_lessthan
func (L *State) LessThan(index1, index2 int) bool {
	return C.lua_compare(L.s, C.int(index1), C.int(index2), C.LUA_OPLT) == 1
}

func (L *State) ObjLen(index int) uint {
	return uint(C.lua_rawlen(L.s, C.int(index)))
}

// lua_len
func (L *State) Len(index int) {
	C.lua_len(L.s, C.int(index))
}

// lua_tointeger
func (L *State) ToInteger(index int) int {
	return int(C.lua_tointegerx(L.s, C.int(index), nil))
}

// lua_tonumber
func (L *State) ToNumber(index int) float64 {
	return float64(C.lua_tonumberx(L.s, C.int(index), nil))
}

// lua_yield
func (L *State) Yield(nresults int) int {
	return int(C.lua_yieldk(L.s, C.int(nresults), 0, nil))
}

func (L *State) pcall(nargs, nresults, errfunc int) int {
	return int(C.lua_pcallk(L.s, C.int(nargs), C.int(nresults), C.int(errfunc), 0, nil))
}

// Pushes on the stack the value of a global variable (lua_getglobal)
func (L *State) GetGlobal(name string) {
	Ck := C.CString(name)
	defer C.free(unsafe.Pointer(Ck))
	C.lua_getglobal(L.s, Ck)
}

// lua_resume
func (L *State) Resume(narg int) int {
	return int(C.lua_resume(L.s, nil, C.int(narg), nil))
}

// lua_setglobal
func (L *State) SetGlobal(name string) {
	Cname := C.CString(name)
	defer C.free(unsafe.Pointer(Cname))
	C.lua_setglobal(L.s, Cname)
}

// Calls luaopen_debug
func (L *State) OpenDebug() {
	C.clua_opendebug(L.s)
}

// Calls luaopen_coroutine
func (L *State) OpenCoroutine() {
	C.clua_opencoroutine(L.s)
}

// lua_insert
func (L *State) Insert(index int) { C.lua_rotate(L.s, C.int(index), 1) }

// lua_remove
func (L *State) Remove(index int) {
	C.lua_rotate(L.s, C.int(index), -1)
	C.lua_settop(L.s, C.int(-2))
}

// lua_replace
func (L *State) Replace(index int) {
	C.lua_copy(L.s, -1, C.int(index))
	C.lua_settop(L.s, -2)
}

// lua_rawgeti
func (L *State) RawGeti(index int, n int) {
	C.lua_rawgeti(L.s, C.int(index), C.longlong(n))
}

// lua_geti
func (L *State) Geti(index int, n int) {
	C.lua_geti(L.s, C.int(index), C.longlong(n))
}

// lua_rawseti
func (L *State) RawSeti(index int, n int) {
	C.lua_rawseti(L.s, C.int(index), C.longlong(n))
}

// lua_seti
func (L *State) Seti(index int, n int) {
	C.lua_seti(L.s, C.int(index), C.longlong(n))
}

// lua_gc
func (L *State) GC(what, data int) int {
	return int(C.lua_gc_compat(L.s, C.int(what), C.int(data)))
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

func luaGC(L *State) int {
	tname := L.CheckString(2)
	L.RemoveUdata(1, tname)
	return 0
}

// SetGC SetGC
func (L *State) SetGC(n int, tname string) int {
	L.LoadString(`
local f, tname = ...
if f and type(f) == 'userdata' and type(tname) == "string" then
    return function(...)
        local t = ... 
        f(t, tname)
    end
end 
return nil
	`)
	L.PushGoFunction(luaGC)
	L.PushString(tname)
	err := L.Call(2, 1)
	if err != nil {
		L.RaiseError(err.Error())
	}
	L.SetField(n, "__gc")
	return 0
}

// SetMetaFunc SetMetaFunc
// (1, 'call'), will set the meta funct __call
func (L *State) SetMetaFunc(n int, name string) int {
	field := "__" + name
	L.LoadString(fmt.Sprintf(`
local t = ...
if t.%s and type(t.%s) == 'userdata' then
    t.%s(t)
end 
	`, name, name, name))
	L.SetField(n, field)
	return 0
}

// LPushError LPushError
func (L *State) LPushError(err error) int {
	if err == nil {
		L.PushNil()
		return 1
	}
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
	val := reflect.ValueOf(ptr)
	pushV(L, val)
	return 1
}

// LPushToTableOrNil LPushToTableOrNil
func (L *State) LPushToTableOrNil(ptr interface{}) int {
	val := reflect.ValueOf(ptr)
	if val.IsNil() {
		L.PushNil()
	} else {
		pushV(L, val)
	}
	return 1
}

func pushV(L *State, val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			return false
		}
		return pushV(L, val.Elem())
	case reflect.String:
		L.PushString(val.String())
	case reflect.Int, reflect.Int64, reflect.Int16:
		L.PushInteger(val.Int())
	case reflect.Uint32, reflect.Uint64, reflect.Uint16:
		L.PushInteger(int64(val.Uint()))
	case reflect.Float32, reflect.Float64:
		L.PushNumber(val.Float())
	case reflect.Bool:
		L.PushBoolean(val.Bool())
	case reflect.Struct:
		if val.Type() == reflect.TypeOf(time.Time{}) {
			t := val.Interface().(time.Time)
			L.PushString(t.Format("2006-01-02 15:04:05"))
		} else {
			fields, _ := getLuaField(val)
			L.NewTable()
			for name, v := range fields {
				if pushV(L, v) {
					L.SetField(-2, name)
				}
			}
		}
	case reflect.Slice:
		L.NewTable()
		j := 1
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i)
			if pushV(L, v) {
				L.RawSeti(-2, j)
				j++
			}
		}
	default:
		return false
	}
	return true
}

// LGetFromTable LGetFromTable
func (L *State) LGetFromTable(ptr interface{}, n int) {
	val := reflect.ValueOf(ptr)
	fields, defaluts := getLuaField(val)
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

// LGetFromGlobal LGetFromGlobal
func (L *State) LGetFromGlobal(ptr interface{}) {
	val := reflect.ValueOf(ptr)
	fields, defaluts := getLuaField(val)
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
	vv := v
	if v.Kind() == reflect.Ptr {
		vv = v.Elem()
	}
	tt := vv.Type()
	for i := 0; i < tt.NumField(); i++ {
		fieldInfo := tt.Field(i)
		if fieldInfo.Anonymous {
			getFields(v.Field(i), fields, defaults)
			continue
		}
		tag := fieldInfo.Tag
		name := tag.Get("lua")
		if name != "" && strings.Contains(name, ",") {
			names := strings.Split(name, ",")
			if len(names) > 1 {
				name = names[0]
				if name == "" {
					name = strings.ToLower(fieldInfo.Name)
				}
				defaults[name] = names[1]
			}
		}
		if name == "" {
			name = fieldInfo.Name
		}
		if name != "" {
			fields[name] = vv.Field(i)
		}
	}
}

func getLuaField(v reflect.Value) (map[string]reflect.Value, map[string]string) {
	fields := make(map[string]reflect.Value)
	defaults := make(map[string]string)
	elem := v //.Elem()
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

// ToArray ToArray
func (L *State) ToArray(ptr interface{}, n int) int {
	if !L.IsTable(n) {
		return L.ArgError(n, "arg must table")
	}
	tt := reflect.TypeOf(ptr)
	if tt.Kind() != reflect.Slice {
		return L.ArgError(n, "val must be slice")
	}
	tlen := int(L.ObjLen(n))
	if tlen > 0 {
		elemT := tt.Elem()
		vals := make([]reflect.Value, 0, tlen)
		for i := 1; i <= tlen; i++ {
			L.RawGeti(n, i)
			v := getV(L, elemT, -1)
			if v != nil {
				vals = append(vals, *v)
			}
			L.Pop(1)
		}
		a0 := reflect.ValueOf(&ptr).Elem()
		vals_ := reflect.Append(a0, vals...)
		a0.Set(vals_)
	}
	return 0
}

func getV(L *State, t reflect.Type, n int) *reflect.Value {
	var val reflect.Value
	switch t.Kind() {
	//case reflect.Ptr:
	//return getV(L, t.Elem())
	case reflect.String:
		val = reflect.ValueOf(L.CheckString(n))
	case reflect.Int, reflect.Int64, reflect.Uint32, reflect.Uint64:
		val = reflect.ValueOf(L.CheckInteger(n))
	case reflect.Float32, reflect.Float64:
		val = reflect.ValueOf(L.CheckNumber(n))
	case reflect.Bool:
		val = reflect.ValueOf(L.ToBoolean(n))
	default:
		return nil
	}
	return &val
}
