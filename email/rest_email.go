//rest email api

package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	
	"io"
	"io/ioutil"
	"os" //for file io
	"strconv"
	"strings"
	"github.com/taknb2nch/go-pop3"
	"net/mail" //parse mail content
)

/**
	测试获取服务列表
	curl -i http://127.0.0.1:8080/
	
	测试获取邮件列表 pop3
	curl -i http://127.0.0.1:8080/email?pop3Addr=pop.163.com&pop3Port=22&user=myusername&pwd=mypasswords&afterId=***
	
*/

type Email struct {
}

//服务列表响应结构定义
var emails []string

type ServiceList struct {
	Service_list interface{} `json:"service_list,omitempty"` //标签小写符合响应json串格式
}


/**
邮件列表响应结构定义
{"id":"***", "reachTime":"", "subject":"", "senderAddr":"", "size":12345}
pop3应该只能得到id

*/
type EmailHeader struct {
	Id         string `json:"id,omitempty"`
	ReachTime  string `json:"reachTime,omitempty"`
	Subject    string `json:"subject,omitempty"`
	SenderAddr string `json:"senderAddr,omitempty"`
	Size       uint64    `json:"size,omitempty"`
}

type EmailContent struct {
	SenderAddr 		string 	 `json:"senderAddr,omitempty"` // 发送方地址
	SenderName 		string   `json:"senderName,omitempty"` // 发送方姓名
	RecverAddrList  string   `json:"recverAddrList",omitempty` // 接收方地址
	CcAddrList      string 	 `json:"ccAddrList,omitempty"` // 抄送地址
	Subject   		string 	 `json:"subject,omitempty"` // 邮件主题
	Content 		string 	 `json:"content,omitempty"` // 正文字符串, 不是base64
	AttachmentList 	string 	 `json:"attachmentList,omitempty"`
}

type EmailHeaderList struct {
	PreheaderList []EmailHeader `json:"preheaderList"`
}


type EmailJsonResponse struct {
	Success bool          `json:"success"`
	Data    interface{}	  `json:"data,omitempty"`
	Error   interface{}   `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code	int    `json:"code,omitempty"`
	Name    string `json:"name,omitempty"`
}

//邮件内容响应结构定义


func init() {
	emails = make([]string, 0)
	emails = append(emails, "yuanlv@126.com")
	emails = append(emails, "yuanlvxg@163.com")
}

func (email *Email) GetServiceList(w rest.ResponseWriter, r *rest.Request) {
	log.Printf("get service list...")
	//w.WriteJson(&email)

	listLen := len(emails)
	log.Printf("service list size=%d", listLen)
	for _, name := range emails {
		log.Printf(name)
	}

	service_list := ServiceList{emails}
	resp := EmailJsonResponse{Success:true, Data:service_list}

	w.WriteJson(&resp)
}

/*
	GET /email 统一入口，通过参数区分具体的请求
 GET /email?pop3Addr=pop.163.com&pop3Port=22&user=myusername&pwd=mypasswords&afterId=***
 GET /email?pop3Addr=pop.163.com&pop3Port=22&user=myusername&pwd=mypasswords&id=***

 tes url
 http://localhost:8080/email?pop3Addr=pop.126.com&pop3Port=110&user=yuanlv@126.com&pwd=***&afterId=1
*/
func (email *Email) GetEmailListAndContent(w rest.ResponseWriter, r *rest.Request) {
	log.Printf("enter into get email list and content...")
	//检验参数是否完整

	getListArgs := r.FormValue("afterId")
	getContentArgs := r.FormValue("id")

	log.Printf("url query params: afterId=%s, id=%s", getListArgs, getContentArgs)

	if getListArgs != "" {
		email.getEmailList(w, r)
	} else if getContentArgs != "" {
		email.getEmailContent(w, r)
	}

	//回复错误信息

}

func (email *Email) getEmailList(w rest.ResponseWriter, r *rest.Request) {
	log.Printf("get email list...")
	emailHeaders := make([]EmailHeader, 0)

	emailHeaders, err := email.getEmailListUsePop3(w, r)
	if err != nil {
		errInfo := ErrorInfo{Code: 403, Name:""}
		errJson := EmailJsonResponse{Success:false, Error:errInfo}
		w.WriteJson(&errJson)
		return
	}

	//get email heaer test
	// emailHeader1 := EmailHeader{Id: "123", ReachTime: "10:00:33", Subject: "subject", SenderAddr: "admin@126.com", Size: 800}
	// emailHeader2 := EmailHeader{Id: "124", ReachTime: "10:10:33", Subject: "subject2", SenderAddr: "admin2@126.com", Size: 8200}
	// emailHeaders = append(emailHeaders, emailHeader1)
	// emailHeaders = append(emailHeaders, emailHeader2)

	preheaderList := EmailHeaderList{emailHeaders}
	preheaderListResp := EmailJsonResponse{Success: true, Data: preheaderList}

	w.WriteJson(&preheaderListResp)
}

//pop3获取到邮件列表后的回调处理函数
func process_pop3(number int, uid, data string, err error) (bool, error) {
	log.Printf("%d, %s\n", number, uid)

	// implement your own logic here
	log.Printf("%s\n", uid)
	return false, nil
}

func (email *Email) getEmailListUsePop3(w rest.ResponseWriter, r *rest.Request) ([]EmailHeader, error) {
	addr := r.FormValue("pop3Addr")
	port := r.FormValue("pop3Port")
	user := r.FormValue("user")
	pwd := r.FormValue("pwd")
	afterId  := r.FormValue("afterId")

	address := addr + ":" + port
	log.Printf("pop3 addr=%s\n", address)
	// if err := pop3.ReceiveMail(addrs, user, pwd, process_pop3); err != nil {
	// 	log.Printf("%v\n", err)
	// 	return err
	// }

	client, err := pop3.Dial(address)

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	defer func() {
		client.Quit()
		client.Close()
	}()

	if err = client.User(user); err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}

	if err = client.Pass(pwd); err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}

	var count int
	var size uint64

	if count, size, err = client.Stat(); err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}

	log.Printf("Count: %d, Size: %d\n", count, size)

	//获取最新一封邮件, 指定id
	headers := make([]EmailHeader, 0)
	id, _ := strconv.Atoi(afterId)
	for i:= id; i <= count; i++ {
		log.Printf("=============")
		log.Printf("i=%d, count=%d\n", i, count)
		if _, size, err = client.List(i); err != nil {
			log.Printf("Error: %v\n", err)
			return nil, err
		}

		log.Printf("email Number: %d, Size: %d\n", i, size)

		var content  string
		header := EmailHeader{}
		if content, err = client.Retr(id); err != nil {
	    	log.Printf("Error: %v\n", err)
	    	return nil, err
		}

		reachTime, subject, senderAddr, err := parseContent(content)
		//skip parse error
		if err != nil {
			log.Printf("parse error id=%d\n", i)
			continue
		}

		header.Id = strconv.Itoa(i)
		header.ReachTime =  reachTime //"10:33:11"
		header.Subject =  subject //"hello"
		header.SenderAddr =  senderAddr //"admin@126.com"
		header.Size = size
		log.Printf("Content length:%d\n", len(content))

		headers = append(headers, header)
	}



	//log.Printf("Content:\n%s\n", content)
	//save to file for test
	//saveFile("d:/test/email.html", content)
	return headers, nil
}

func parseContent(content string) (string, string, string, error){
	r := strings.NewReader(content)
	m, err := mail.ReadMessage(r)
	if err != nil {
		return "", "", "", nil
	}

	header := m.Header

	body, err := ioutil.ReadAll(m.Body)
	log.Printf("%s", body)
	return header.Get("Date"), header.Get("Subject"), header.Get("From"), nil

	//return "2016-3-16", "hello", "admin@126.com", nil

}

//ioutil
func saveFile(filename, content string) (int, error){
	fl, err := os.OpenFile(filename, os.O_CREATE, 0644)
	if err != nil {
	    return 0, err
	}
	defer fl.Close()
	n, err := fl.Write([]byte(content))
	if err == nil && n < len(content) {
	    err = io.ErrShortWrite
	}
	return n, err
}

func (email *Email) getEmailContent(w rest.ResponseWriter, r *rest.Request) {
	log.Printf("get email content...")
}

func (email *Email) SendEmailContent(w rest.ResponseWriter, r *rest.Request) {
	log.Printf("send email...")

}

func (email *Email) DeleteEmail(w rest.ResponseWriter, r *rest.Request) {
	log.Printf("delete email...")
}

func main() {
	email := Email{}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/", email.GetServiceList),
		rest.Get("/email", email.GetEmailListAndContent),
		rest.Put("/email", email.SendEmailContent),
		rest.Delete("/email", email.DeleteEmail),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Printf("email api server start...")
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

