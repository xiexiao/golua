package gen

import (
	"testing"
)

func TestCheckHeader(t *testing.T) {
	fname := "t.go"
	code := CodeHeader
	err := CheckHeader(fname, code)
	err1 := CheckHeader(fname, "")
	t.Log(err)
	t.Log(err1)
}

func TestGenError(t *testing.T) {
	fname := "t.go"
	code := ``
	buf, err := GenLua(fname, []byte(code), true)
	t.Log(string(buf))
	t.Log(err)
}

func TestGenUseSame(t *testing.T) {
	fname := "t.go"
	code := `package test
type Code struct {
    ID int ` + "`lua`" + `
}

func (c *Code) Test1(a, b string) {
}
    `
	buf, err := GenLua(fname, []byte(code), true)
	t.Log(string(buf))
	t.Log(err)
}

func TestGenBytes(t *testing.T) {
	fname := "t.go"
	code := `package test
type Code struct {
    ID int ` + "`lua`" + `
    Buf []byte ` + "`lua:name`" + `
}

func CodeNew(s byte)(*Code, error) {
    return nil, nil
}

func (c *Code) ResetBuffer(buf []byte) []byte {
    return nil
}

func (c *Code) ResetString(buf []string) []string{
    return nil
}
    `
	buf, err := GenLua(fname, []byte(code), true)
	t.Log(string(buf))
	t.Log(err)
}

func TestGenNew(t *testing.T) {
	fname := "t.go"
	code := `package test
type Code struct {
    ID int ` + "`lua`" + `
    Name string ` + "`lua:name`" + `
}

func CodeNew(s byte, s2 string)(*Code, error) {
    return nil, nil
}

func (c *Code) Check() {
}
    `
	buf, err := GenLua(fname, []byte(code), true)
	t.Log(string(buf))
	t.Log(err)
}

func TestGenTplError(t *testing.T) {
	buf, err := codeExec(codeTemplate, 0, true)
	t.Log(err)
	t.Log(string(buf))
}

func TestGen(t *testing.T) {
	fname := "t.go"
	code := `package test
import (
    "bytes"
)
type Code struct {
    ID int ` + "`lua`" + `
    Name string ` + "`lua:name`" + `
}

func (c *Code) Test(s string) (byte, error) {
    return 0, nil
}

func (c *Code) Test2(s *string) (*bytes.Buffer, error) {
    return nil, nil
}

func (c *Code) Test3(s *string) ([]*bytes.Buffer, error) {
    return nil, nil
}

func (c *Code) Test4(s *string) ([]bytes.Buffer, error) {
    return nil, nil
}

//lua:get
func (c *Code) Get(s *Node, 
    num int16, num1 uint16, 
    num2 int32, num3 uint32, 
    num4 int64, num5 uint64, 
    num6 float32, num7 float64, 
    v0 bool) (int16, uint16,
    int32, uint32, 
    int64, uint64, 
    float32, float64,
    bool){
	return num, num1, num2, num3, num4, num5, num6, num7, v0
}

//lua:print
func (c *Code) Print(s *Node, num int,
	num1 int32, num2 int64, num3 uint32,
	num4 string, num5 byte) int {
	fmt.Println(s)
	fmt.Println(num)
	fmt.Println(num1)
	fmt.Println(num2)
	fmt.Println(num3)
	fmt.Println(num4)
	fmt.Println(num5)
	return num * num
}

//lua:get_strings
func (c *Code) Get_Strings() ([]string, error) {
	list := make([]string, 0)
	list = append(list, "test1")
	return list, nil
}

//lua:print2
func (c *Code) Print2(s string) (uint64, error) {
    return 1, nil
}
    `
	buf, err := GenLua(fname, []byte(code), true)
	t.Log(string(buf))
	t.Log(err)
}

func TestSetMetatable(t *testing.T) {
	fname := "t.go"
	code := `package test
type Cat struct {
    Name string ` + "`lua:name`" + `
}

//lua:new,__make_udata
func (c *Cat) New(a string, b int) (*Dog, error) {
	return nil, nil
}`
	buf, err := GenLua(fname, []byte(code), true)
	if err != nil {
		t.Errorf("Gen Code get error:%s\n", err)
	}
	t.Log(string(buf))
}
