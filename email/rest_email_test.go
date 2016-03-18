package test
import(

	"testing"
	"log"
	"io/ioutil"
    "github.com/mxk/go-imap/imap"
    "net"
    "fmt"
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