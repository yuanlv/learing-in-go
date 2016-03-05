package main 

import(
	"io"
	"fmt"
	"net/http"
	"github.com/yuanlv/learning-in-go/api" //引用目录下的包，写到包所在目录即可
	"strings"
	"strconv"

)

func process(w http.ResponseWriter, r *http.Request){
	//检查url路径，转到实际的处理方法
	fmt.Println(r.URL.Path)
	io.WriteString(w, r.URL.Path)
	io.WriteString(w, "<br/>")
	//fmt.Fprintf(w, r.URL.Path)
	funcName := r.URL.Path

	var ctx = (int64)(0)
	//get ctx from args
	queryValues := r.URL.Query()
	ctxStr := queryValues.Get("ctx")
	fmt.Println("get ctx from query: ", ctxStr)

	if(ctxStr != ""){
		ctx, err := strconv.ParseInt(ctxStr, 10, 32)
		if err != nil {
			fmt.Println("get ctx error")
			fmt.Printf("parseInt result=%d", ctx)
			io.WriteString(w, "get ctx error")
			//send error msg
			return
		}
	}

	if(strings.Contains(funcName, "encCreate")){
		ctx, err := enc.EncCreate("testid")
		if err == nil {
			fmt.Printf("get enc ctx %d", ctx)
			fmt.Fprintf(w, "%d", ctx)
		}
		
	}else if(strings.Contains(funcName, "encUpdate")){
		out, err := enc.EncUpdate(ctx, []byte("test in data"))
		if err == nil {
			fmt.Println("enc out===")
			fmt.Println(string(out))
			io.WriteString(w, string(out))
		}		
	}else if(strings.Contains(funcName, "encFinal")){
		out, err := enc.EncFinal(ctx)
		if err == nil {
			io.WriteString(w, string(out))
		}		
	}else if(strings.Contains(funcName, "encRelease")){
		err := enc.Release(ctx)
		if err == nil {
			io.WriteString(w, "release ok")
		}
	}

}

func startServer(port string){
	
	http.HandleFunc("/api/encCreate", process)
	http.HandleFunc("/api/encUpdate", process)
	http.HandleFunc("/api/encFinal", process)
	http.HandleFunc("/api/encRelease", process)

	http.ListenAndServe(":"+port, nil)
}

func main(){
	fmt.Println("start api server...")
	startServer("6000")
}