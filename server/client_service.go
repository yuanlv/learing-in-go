package main

import (
	"fmt"
	"github.com/hprose/hprose-go"
)

type clientStub struct {
	Hello func(string) string
	Sum   func(int, int) int
}

func main(){
	client := hprose.NewClient("http://127.0.0.1:6000/")
	var resp *clientStub
	client.UseService(&resp)

	fmt.Println(resp.Hello("world"))
	fmt.Println(resp.Sum(2, 4))

}