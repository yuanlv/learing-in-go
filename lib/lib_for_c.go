/**
	use go code, generate c shared library
	to learn more detail, please see this blog
	https://blog.filippo.io/building-python-modules-with-go-1-5/

	build cmd: go build -buildmode=c-shared -o lib_for_c.so lib_for_c.go
	@author: yuanlv
	@date  : 2016-01-30
*/

package main

import "C"

//export Sum
func Sum(a, b int) int {
	return a + b
}


func main(){

}