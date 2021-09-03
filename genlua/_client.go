package main

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
func (c *Client) Get_Nodes() ([]Node, error) {
    var list []Node
    list = append(list, Node{
        ID : 1,
        Name: "test1",
    })

    list = append(list, Node{
        ID : 2,
        Name: "test2",
    })
	return list, nil
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

type Node struct {
    ID int
    Name string
}
