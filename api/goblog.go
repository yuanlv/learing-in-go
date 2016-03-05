package main 

import(
	"fmt"
	"net/http"
	"io/ioutil"
)

func getIndex(){

	resp, err := http.Get("http://qxu1606580098.my3w.com/wx/index.php?echostr=123");
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%v\n", string(body))
}

func main(){
	getIndex()
}