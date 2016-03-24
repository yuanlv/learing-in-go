package main


// #include <stdio.h>
// #include <stdint.h>
// #include <string.h>
/*struct columns {
    int column1;
    int column2;
    int column3;
};
int sum_columns(struct columns a) {
    return a.column1 + a.column2 + a.column3;
}
int sum_vals(int a, int b) {
    return a + b;
}

int update(unsigned char* in, unsigned int in_len, unsigned char* out, unsigned int *out_len){
    printf("input data len=%d\n", in_len);
    *out_len = 10;
    memset(out, 0x10, *out_len);

    return 0;

}
*/
import "C"
import "fmt"
import "unsafe"

func main() {
    c := C.struct_columns{15, 30, 45}
    sum := C.sum_columns(c)
    fmt.Println(sum)
    var a int = 15
    var b int = 30
    s := C.sum_vals((C.int)(a), (C.int)(b))//调用C语言的函数
    fmt.Println(s)
    var goSum int = int(sum)//C语言的int转成Golang的int类型
    fmt.Println(goSum)

    
    in := []byte("test")
    inLen := len(in)
    out := (*C.uchar)(C.malloc(16))
    outLen := C.uint(16)
    r := C.update((*C.uchar)(&in[0]), C.uint(inLen),  out, &outLen)
    if r == 0 {
        fmt.Printf("outlen=%d\n", outLen)

        dst := fmt.Sprintf("%02X", C.GoBytes(unsafe.Pointer(out), C.int(outLen)) )
        fmt.Println(dst)
    }

    //C.free(unsafe.Pointer(out))

}