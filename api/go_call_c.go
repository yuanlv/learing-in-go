package main


import(
	"fmt"
	"syscall"
	//"unsafe"
)

func testFail(){
	h, err := syscall.LoadLibrary("msvcrt")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer syscall.FreeLibrary(h)

	proc, err := syscall.GetProcAddress(h, "memcpy")
	if err != nil {
		fmt.Println(err)
		return
	}

	// dst := unsafe.Pointer(uintptr(1))
	// src := unsafe.Pointer(uintptr(2))
	// count := unsafe.Pointer(uintptr(4))
	dst := uintptr(1)
	src := uintptr(2)
	count := uintptr(3)

	fmt.Println("test memcpy error in go...")
	r, _, _ := syscall.Syscall(uintptr(proc), 3, dst, src, count)

	fmt.Println(r)
	fmt.Println("should not run here")

}


func main(){
	testFail()

}