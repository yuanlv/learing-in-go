package rest_email

import (
 	"github.com/yuanlv/mahonia" //for gbk to utf8
    "strings"
    "log"
    "encoding/base64"
    "io/ioutil"
    "mime/quotedprintable"
    "github.com/jhillyerd/go.enmime" //parse mail body
)

//根据编码格式获取转码后的主题内容
func GetRealSubject(subject string) (string){
	//先根据空格截取多个主题内容项, 合并输出最终的主题内容
	elements := strings.Split(subject, " ")
	var realSubject string
	for i:=0; i<len(elements); i++ {
		subject = elements[i]
		log.Printf(subject)
		curSub := getRealSubject(subject)
		// log.Printf("current subject in for===")
		// log.Printf(curSub)
		realSubject += curSub
	}
	log.Printf("in util.go: get real subject")
	log.Printf(realSubject)
	return realSubject	
}

func getRealSubject(subject string) (string){
	lang, encoding, content := getEncodingStr(subject)
	if content == "" {
		return ""
	}

	lang = strings.ToLower(lang)
	encoding = strings.ToLower(encoding)
	
	return getDecodingContent(lang, encoding, content)
}

func getDecodingContent(lang, encoding, content string) (string){

	if encoding == "b" {
		decoded, err := base64.StdEncoding.DecodeString(content)
		if err != nil {
			log.Printf("decode error:%s", err)
			return ""
		}

		if lang == "utf-8" {
			//log.Printf("utf-8, noneed to convert: %s\n", string(decoded))
			return string(decoded)
		}

		//非u8编码语言需要统一转为u8
		dec := mahonia.NewDecoder(lang)
		if u8Decoded, ok := dec.ConvertStringOK(string(decoded)); ok {
			//log.Printf("%s convert to utf-8: %s\n", lang, string(decoded))
			return string(u8Decoded) //注意作用域
		}
		
	}else if encoding == "q" {
		//通过字符串构造Reader
		r := quotedprintable.NewReader(strings.NewReader(content)) 
		body, err := ioutil.ReadAll(r)
		if err != nil {
			log.Printf("read content error: %s", err)
			return ""
		}
		log.Printf("quoted decode: %s\n", string(body))
		return string(body)
	}
	
	//转换失败，返回空
	return ""
}


//解析带编码格式邮件主题内容, 也可能是没有任何编码（纯ascii英文字符，待确认)
func getEncodingStr(content string) (string, string, string) { 
	elements := strings.Split(content, "?")
	if len(elements) < 3  {
		return "utf-8", "", "" //默认字符集
	}

	log.Printf("elements size %d", len(elements))
	// for _, e := range elements {
	// 	log.Printf(e)
	// }
	//language, encoding, real content 
	return elements[1], elements[2], elements[3]
}


//获取正文部分解码后的内容
/* 正文格式

------=_Part_52_1608986047.1312859369078
Content-Type: multipart/alternative;
    boundary="----=_Part_53_1114408905.1312859369078"

------=_Part_53_1114408905.1312859369078
Content-Type: text/plain; charset="gb18030"
Content-Transfer-Encoding: base64

1eLKx9K7t+JIVE1MuPHKvdPKvP6jrMfrx9C7u7W9SFRNTMrTzbw=

------=_Part_53_1114408905.1312859369078
Content-Type: text/html;charset=UTF-8
Content-Transfer-Encoding: base64



*/

func GetRealBody(content string) (string){
	//check encoding: base64 or quoted
	var realBody string
	var charset string
	if strings.Contains(content, "UTF-8"){
		charset = "utf-8"
		log.Printf("body charset: utf-8")
	}

	if strings.Contains(content, "base64"){
		log.Printf("body encoding is base64")
		realBody = getDecodingContent(charset, "b", content)
	}else if strings.Contains(content, "quoted-printable") {
		log.Printf("body encoding is quoted-printable")
		realBody = getDecodingContent(charset, "q", content)
	}else{
		realBody = content
	}

	//test
	realBody = getDecodingContent(charset, "q", content)

	return realBody
}

func GetAllMimePartBody(mimeBody *enmime.MIMEBody) string{
	log.Printf("enter into GetAllMimePartBody")
	root := mimeBody.Root
	log.Printf(string(root.Content()))

	if child := root.FirstChild(); child != nil {
		log.Printf(string(child.Content()))
	}

	return ""
}