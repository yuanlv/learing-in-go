package main

import (
	"fmt"
	//"strings"
)

type NTP struct{
	Version string
}

type NTPS struct{
	NTP  //匿名字段
	mac []byte
}

func (ntp NTP)Show(){
	fmt.Println("show ntp info")
}

/*
func (ntps NTPS)Show(){
	fmt.Println("show ntps info")
}*/

func main() {
	ntps := new(NTPS)
	ntps.Show()
}