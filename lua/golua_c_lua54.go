//+build lua54

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
