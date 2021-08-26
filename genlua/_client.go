package gen

import (
	"fmt"
)

type Client struct {
	EngineAddr string `lua`
	Port       int    `lua`
}

func ClientNew(addr string, port int) (*Client, error) {
	clt := &Client{
		EngineAddr: addr,
		Port:       port,
	}
	return clt, nil
}

//lua:get_nodes,check_error
func (c *Client) Get_Nodes() (string, error) {
	return "", nil
}

//lua:print
func (c *Client) Print(s string, num int,
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

//lua:get_summary,check_error
func (clt *Client) Get_Summary() (string, error) {
	return "", nil
}
