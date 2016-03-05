package main

import (
	"fmt"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
)
const API_URL = "http://zuche.ztu.wang/api.php?op="

type Resp struct {
	Msg string
	err string
	userid string
	username string
}


func register(){
	apiUrl := API_URL + "register"
	fmt.Println("url=" + apiUrl)
	resp, err := http.PostForm(apiUrl, 
		url.Values{"username":{"18001152652"}, "password":{"fke123456"}})

	if err != nil {
		fmt.Println("call api error")
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%v\n", string(body))

	//get msg
	var msg Resp
	err = json.Unmarshal(body, &msg)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println("msg:", msg.Msg)

}

func login(){
	apiUrl := API_URL + "login"
	fmt.Println("url=" + apiUrl)
	resp, err := http.PostForm(apiUrl,
		url.Values{"username":{"18001152652"}, "password":{"fke123456"}})

	if err != nil {
		fmt.Println("call api error")
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()


	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%v\n", string(body))
	
	//get msg
	var msg Resp
	err = json.Unmarshal(body, &msg)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println("msg:", msg.Msg)
}

func main() {

	register()
	login()

}