package main

import (
	//"errors"
	"github.com/hprose/hprose-go"
	"net/http"
	"fmt"
)

func hello(name string) string {
	fmt.Println("func hello was called...")
	return "server echo : hello " + name
}

type myService struct{}

func (myService) Sum(a int, b int) (int) {
	fmt.Println("Sum hell was called...")
	return a+b
}

func main(){
	service := hprose.NewHttpService()
	service.AddFunction("hello", hello)
	service.AddMethods(myService{})
	http.ListenAndServe(":6000", service)
}