package rest_email
import(

	"testing"
	"log"
	//"io/ioutil"
    //"github.com/mxk/go-imap/imap"
    //"net"
    "fmt"
    "encoding/base64"
    "github.com/yuanlv/mahonia" //for gbk to utf8
    "strings"
)

func TestFor(*testing.T){
	for i:=1; i<3; i++ {
		log.Printf("i=%d\n", i)
	}
}


func get1st(a, b interface{}) interface{} {
    return a
}

func TestIcmp(*testing.T) {
 
}

func TestBase64(*testing.T){
	msg := "Hello, 世界"
	encoded := base64.StdEncoding.EncodeToString([]byte(msg))
	fmt.Println(encoded)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
	    fmt.Println("decode error:", err)
	    return
	}
	fmt.Println(string(decoded))
}


//根据编码格式获取转码后的主题内容
func getRealSubject(lang, encoding, content string) ([]byte){
	encoding = strings.ToLower(encoding)
	if encoding == "b" {
		decoded, err := base64.StdEncoding.DecodeString(content)
		if err != nil {
			log.Printf("decode error:%s", err)
			return nil
		}

		if lang == "utf-8" {
			log.Printf(string(decoded))
			return decoded
		}

		//非u8编码语言需要统一转为u8
		dec := mahonia.NewDecoder(lang)
		if decoded, ok := dec.ConvertStringOK(string(decoded)); ok {
			log.Printf(string(decoded))
		}

		return decoded
	}else if encoding == "q" {
		return nil
	}else {
		return nil
	}
	
}

func getEncodingStr(content string) (string, string, string) {
	elements := strings.Split(content, "?")
	if len(elements) < 3  {
		return "utf-8", "", "" //默认字符集
	}

	log.Printf("elements size %d", len(elements))
	for _, e := range elements {
		log.Printf(e)
	}
	//language, encoding, real content 
	return elements[1], elements[2], elements[3]
}

func TestSubjectBase64(t *testing.T){
	base64Subject := "?UTF-8?B?5a6J5YWo5o+Q6YaS77ya5a+G56CB5L+u5pS5?="
	//quotedSubject := "?utf-8?Q?DaoVoice=EF=BC=8C=E5=89=8D=E6=89=80=E6=9C=AA=E6=9C=89=E7=9A=84=E8=BF=90=E8=90=A5=E4=B9=8B=E9=81=93?="

	//getSubject(subject)
	
	lang, encoding, content := getEncodingStr(base64Subject)
	log.Printf("encoding: %s  content：%s", encoding, content)
	lang = strings.ToLower(lang)
	getRealSubject(lang, encoding, content)
}

func TestBodyBase64(t *testing.T){
	body := "1eLKx9K7t+JIVE1MuPHKvdPKvP6jrMfrx9C7u7W9SFRNTMrTzbw="

	getRealSubject("gb18030", "b", body)

}