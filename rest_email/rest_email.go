//rest email api

package rest_email

import (
	"bytes"
	"encoding/json"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jhillyerd/go.enmime" //parse mail body
	"github.com/mxk/go-imap/imap"
	//"github.com/scorredoira/email"
	"github.com/taknb2nch/go-pop3"
	"io"
	"io/ioutil"
	"log"
	"net/mail" //parse mail content
	"os" //for file io
	"strconv"
	"strings"
	"time"
	"net/smtp" //for send email
	//辅助方法
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
	Size       uint64 `json:"size,omitempty"`
}

type Attachment struct {
	Name string `json:"name,omitempty"`
	Body string `json:"b64Body,omitempty"`
}

type EmailContent struct {
	SenderAddr     string      `json:"senderAddr,omitempty"`     // 发送方地址
	SenderName     string      `json:"senderName,omitempty"`     // 发送方姓名
	RecverAddrList interface{} `json:"recverAddrList,omitempty"` // 接收方地址
	CcAddrList     interface{} `json:"ccAddrList,omitempty"`     // 抄送地址
	Subject        string      `json:"subject,omitempty"`        // 邮件主题
	Content        string      `json:"content,omitempty"`        // 正文字符串, 不是base64
	Attachments    interface{} `json:"attachmentList,omitempty"`
}

type EmailHeaderList struct {
	PreheaderList []EmailHeader `json:"preheaderList"`
}

type EmailJsonResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

//解析发送邮件json数据时使用,必须指定实际类型
type EmailSendContent struct{
	RecverAddrList []string `json:"recverAddrList,omitempty"` // 接收方地址
	CcAddrList     []string `json:"ccAddrList,omitempty"`     // 抄送地址
	Subject        string      `json:"subject,omitempty"`        // 邮件主题
	Content        string      `json:"content,omitempty"`        // 正文字符串, 不是base64
	Attachments    []Attachment `json:"attachmentList,omitempty"`
}
type EmailJsonRequest struct {
	Success bool         `json:"success"`
	Data EmailSendContent `json:"data,omitempty"`
}

type ErrorInfo struct {
	Code int    `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}

func (email *Email) GetServiceList(w rest.ResponseWriter, r *rest.Request) {
	log.Printf("get service list...")

	service_list := ServiceList{Service_list: "email"}
	resp := EmailJsonResponse{Success: true, Data: service_list}

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
		errInfo := ErrorInfo{Code: 403, Name: ""}
		errJson := EmailJsonResponse{Success: false, Error: errInfo}
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

//pop3登录认证，成功后获取client对象
func (email *Email) getPop3EmailClient(address, user, pwd string) (client *pop3.Client, err error) {
	client, err = pop3.Dial(address)

	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}

	if err = client.User(user); err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}

	if err = client.Pass(pwd); err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}

	return client, nil
}

func (email *Email) getEmailListUsePop3(w rest.ResponseWriter, r *rest.Request) ([]EmailHeader, error) {
	addr := r.FormValue("pop3Addr")
	port := r.FormValue("pop3Port")
	user := r.FormValue("user")
	pwd := r.FormValue("pwd")
	afterId := r.FormValue("afterId")

	address := addr + ":" + port
	log.Printf("pop3 addr=%s\n", address)

	client, err := email.getPop3EmailClient(address, user, pwd)
	if err != nil {
		return nil, err
	}

	defer func() {
		client.Quit()
		client.Close()
	}()

	var count int
	var size uint64

	if count, size, err = client.Stat(); err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}

	log.Printf("Count: %d, Size: %d\n", count, size)

	//获取指定afterid及以后的所有邮件列表
	headers := make([]EmailHeader, 0)
	id, _ := strconv.Atoi(afterId)
	for i := id; i <= count; i++ {
		log.Printf("=============")
		log.Printf("i=%d, count=%d\n", i, count)
		if _, size, err = client.List(i); err != nil {
			log.Printf("Error: %v\n", err)
			return nil, err
		}

		log.Printf("email Number: %d, Size: %d\n", i, size)

		header := EmailHeader{}
		header.Id = strconv.Itoa(i)

		headers = append(headers, header)
	}

	return headers, nil
}

func parseContent(content string) (string, string, string, error) {
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
func saveFile(filename, content string) (int, error) {
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

func (email *Email) getEmailContentUsePop3(w rest.ResponseWriter, r *rest.Request) {
	addr := r.FormValue("pop3Addr")
	port := r.FormValue("pop3Port")
	user := r.FormValue("user")
	pwd := r.FormValue("pwd")
	id := r.FormValue("id")

	address := addr + ":" + port
	client, err := email.getPop3EmailClient(address, user, pwd)
	if err != nil {
		errorInfo := ErrorInfo{Code: 403, Name: ""}
		resp := EmailJsonResponse{Success: false, Error: errorInfo}
		w.WriteJson(&resp)
	}

	defer func() {
		client.Quit()
		client.Close()
	}()

	var body string
	content := EmailContent{}
	i, _ := strconv.Atoi(id)
	if body, err = client.Retr(i); err != nil {
		log.Printf("Error: %v\n", err)
		errorInfo := ErrorInfo{Code: 403, Name: ""}
		resp := EmailJsonResponse{Success: false, Error: errorInfo}
		w.WriteJson(&resp)
		return
	}

	_, subject, senderAddr, err := parseContent(body)
	//skip parse error
	if err != nil {
		log.Printf("parse email body error")
		errorInfo := ErrorInfo{Code: 403, Name: ""}
		resp := EmailJsonResponse{Success: false, Error: errorInfo}
		w.WriteJson(&resp)
		return
	}

	content.SenderAddr = senderAddr
	content.SenderName = ""
	content.Subject = subject
	content.Content = body

	resp := EmailJsonResponse{Success: true, Data: content}
	w.WriteJson(&resp)
}

func (emaill *Email) getEmailContentUseIMAP(w rest.ResponseWriter, r *rest.Request) {
	addr := r.FormValue("imapAddr")
	port := r.FormValue("imapPort")
	user := r.FormValue("user")
	pwd := r.FormValue("pwd")
	id := r.FormValue("id")

	address := addr + ":" + port
	log.Printf("get email content use imap %s\n", address)

	var (
		c   *imap.Client
		cmd *imap.Command
		rsp *imap.Response
	)

	// Connect to the server
	c, err := imap.Dial(addr)
	if err != nil {
		log.Printf("dial %s error\n", address)
		return
	}

	// Remember to log out and close the connection when finished
	defer c.Logout(30 * time.Second)

	// Print server greeting (first response in the unilateral server data queue)
	log.Printf("Server says hello:%s", c.Data[0].Info)
	c.Data = nil

	// Enable encryption, if supported by the server
	if c.Caps["STARTTLS"] {
		c.StartTLS(nil)
	}

	// Authenticate
	if c.State() == imap.Login {
		c.Login(user, pwd)
	}

	// List all top-level mailboxes, wait for the command to finish
	cmd, _ = imap.Wait(c.List("", "%"))

	// Print mailbox information
	log.Printf("\nTop-level mailboxes:")
	for _, rsp = range cmd.Data {
		log.Printf("|--%s", rsp.MailboxInfo())
	}

	// Check for new unilateral server data responses
	for _, rsp = range c.Data {
		log.Printf("Server data:%s", rsp)
	}
	c.Data = nil

	// Open a mailbox (synchronous command - no need for imap.Wait)
	c.Select("INBOX", true)
	log.Printf("\nMailbox status:%s\n", c.Mailbox)
	if c.Mailbox == nil {
		resp := EmailJsonResponse{Success: true}
		w.WriteJson(&resp)
		return
	}

	// Fetch the headers of the 10 most recent messages
	set, _ := imap.NewSeqSet("")
	set.Add(id)
	// if c.Mailbox.Messages >= 10 {
	// 	set.AddRange(c.Mailbox.Messages-9, c.Mailbox.Messages) //测试只取最新一封邮件
	// } else {
	// 	set.Add("1:*")
	// }
	cmd, _ = c.Fetch(set, "RFC822.HEADER", "RFC822.TEXT") //指定要获取的内容

	// Process responses while the command is running
	log.Printf("\nget mail [%s] messages:", id)
	for cmd.InProgress() {
		// Wait for the next response (no timeout)
		c.Recv(-1)

		// Process command data
		for _, rsp = range cmd.Data {
			header := imap.AsBytes(rsp.MessageInfo().Attrs["RFC822.HEADER"])
			if msg, _ := mail.ReadMessage(bytes.NewReader(header)); msg != nil {
				subject := msg.Header.Get("Subject")
				log.Printf("|--%s", subject)

				realSubject := GetRealSubject(subject)
				log.Printf("in rest_email.go: get real subject")
				log.Printf(realSubject)

				senderAddr := msg.Header.Get("From")
				recverAddrList := msg.Header.Get("To")

				realSenderAddr := GetRealSubject(senderAddr)
				realRecverAddrList := GetRealSubject(recverAddrList)

				body := imap.AsBytes(rsp.MessageInfo().Attrs["RFC822.TEXT"])
				//log.Printf("email body: %s", body)
				//realBody  := GetRealBody(string(body))
				headerAndBody := make([]byte, len(header)+len(body))
				copy(headerAndBody, header)
				copy(headerAndBody[len(header):], body)

				msg, _ := mail.ReadMessage(bytes.NewReader(headerAndBody))
				mime, _ := enmime.ParseMIMEBody(msg)
				realBody := mime.Text //如果原始内容为html，会去掉html元素标签
				log.Printf("real body: %s", realBody)

				//获取MIMEPart所有节点内容
				// log.Printf("root ======================")
				// root := mime.Root
				// if root != nil {
				// 	log.Printf(string(root.Content()))
				// 	log.Printf("child==========")
				// 	if child := root.FirstChild(); child != nil {
				// 		log.Printf(string(child.Content()))
				// 	}

				// }

				attachments := mime.Attachments
				log.Printf("attachments len=%d", len(attachments))
				count := len(attachments)
				var attachmentList []Attachment = nil
				var data EmailContent
				if count > 0 {
					attachmentList = make([]Attachment, count)
					for i := 0; i < len(attachments); i++ {
						name := attachments[i].FileName()
						content := attachments[i].Content() //todo encode by base64
						log.Printf("name===%s", name)
						attachmentList[i] = Attachment{Name: name, Body: string(content)}
					}
					
				}
				data = EmailContent{Subject: realSubject, SenderAddr: realSenderAddr,
					RecverAddrList: realRecverAddrList, Content: realBody,
					Attachments: attachmentList}	
				

				
				resp := EmailJsonResponse{Success: true, Data: data}
				w.WriteJson(resp)
			}
		}
		cmd.Data = nil

		// Process unilateral server data
		for _, rsp = range c.Data {
			log.Printf("Server data:%s", rsp)
		}
		c.Data = nil
	}

	// Check command completion status
	if rsp, err := cmd.Result(imap.OK); err != nil {
		if err == imap.ErrAborted {
			log.Printf("Fetch command aborted\n")
		} else {
			log.Printf("Fetch error:%s\n", rsp.Info)
		}
	}

}

//获取与id对应的邮件内容
func (email *Email) getEmailContent(w rest.ResponseWriter, r *rest.Request) {
	log.Printf("get email content...")
	addr := r.FormValue("pop3Addr")
	//使用pop3Addr参数区分协议类型
	log.Printf("pop3 addr=%s\n", addr)

	if addr != "" {
		email.getEmailContentUsePop3(w, r)
	} else {
		email.getEmailContentUseIMAP(w, r)
	}

}

//防止与email包冲突
func (*Email) SendEmailContent(w rest.ResponseWriter, r *rest.Request) {
	log.Printf("send email...")

	addr := r.FormValue("smtpAddr")
	port := r.FormValue("smtpPort")
	user := r.FormValue("user")
	pwd := r.FormValue("pwd")

	address := addr + ":" + port
	log.Printf("smtp address %s\n", address)

	//从请求报文body中解析json串，获取发文地址及主题、正文、附件等
	if r.Body == nil {
		log.Printf("request body is null")
		return
	}
	reqJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("when read request body error: %s", err)
		return
	}
	defer r.Body.Close()
	log.Printf("request body: %s", string(reqJson))

	var emailReq EmailJsonRequest
	err = json.Unmarshal(reqJson, &emailReq)
	if err != nil {
		log.Printf("unmarshal req json error, %s", err)
		return
	}
	log.Printf("%v", emailReq)

	log.Printf("from addrs: %s", user)
	to := emailReq.Data.RecverAddrList
	log.Printf("send to addrs: %v", to)
	auth := smtp.PlainAuth("", user, pwd, addr)
	str := strings.Replace("From: "+user+"~To: "+ to[0] +"~Subject: "+ "this is subject" +"~~", "~", "\r\n", -1) + "this is body msg"
	//多个收件人时需要循环发送，一次构造一个收件人[]string
	err = smtp.SendMail(address, auth, user, []string{to[0]}, []byte(str))	
	
	//处理附件，从emailReq.Data.Attachements中取出名称和内容

	//使用net/smtp发送邮件


	resp := EmailJsonResponse{}
	if err == nil {
		log.Printf("send mail ok\n")
		resp.Success = true

	} else {
		log.Printf("send mail error: %s", err)
		resp.Success = false
		resp.Error = ErrorInfo{Code: 404, Name: "send mail error"}
	}

	w.WriteJson(&resp)
}

//use imap to delete a mail
func (email *Email) DeleteEmail(w rest.ResponseWriter, r *rest.Request) {
	log.Printf("delete email...")
	addr := r.FormValue("imapAddr")
	//port := r.FormValue("smtpPort")
	user := r.FormValue("user")
	pwd := r.FormValue("pwd")
	id := r.FormValue("id")

	//check params, return error

	//create imap client
	var (
		c   *imap.Client
		cmd *imap.Command
		rsp *imap.Response
	)

	// Connect to the server
	c, _ = imap.Dial(addr)

	// Remember to log out and close the connection when finished
	defer c.Logout(30 * time.Second)

	// Print server greeting (first response in the unilateral server data queue)
	log.Printf("Server says hello:%s", c.Data[0].Info)
	c.Data = nil

	// Enable encryption, if supported by the server
	if c.Caps["STARTTLS"] {
		c.StartTLS(nil)
	}

	// Authenticate
	if c.State() == imap.Login {
		c.Login(user, pwd)
	}

	// List all top-level mailboxes, wait for the command to finish
	cmd, _ = imap.Wait(c.List("", "%"))

	// Print mailbox information
	log.Printf("\nTop-level mailboxes:")
	for _, rsp = range cmd.Data {
		log.Printf("|--%s", rsp.MailboxInfo())
	}

	// Check for new unilateral server data responses
	for _, rsp = range c.Data {
		log.Printf("Server data:%s", rsp)
	}
	c.Data = nil

	// Open a mailbox (synchronous command - no need for imap.Wait)
	c.Select("INBOX", true)
	log.Printf("\nMailbox status:%s\n", c.Mailbox)
	if c.Mailbox == nil {
		resp := EmailJsonResponse{Success: true}
		w.WriteJson(&resp)
		return
	}

	//use Expunge to delete a mail
	set, _ := imap.NewSeqSet("")
	mid, _ := strconv.Atoi(id)
	set.AddNum(uint32(mid))

	//delete mail
	cmd, err := c.Expunge(set)
	if err != nil {
		log.Printf("%v ", cmd)
	} else {
		log.Printf("delete mail ok")
	}

}
