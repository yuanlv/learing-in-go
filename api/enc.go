package enc

import(
	"fmt"

)

func EncCreate(appId string) (ctx int64, err error){
	fmt.Println("call encCreate begin")
	ctx = 123
	
	return ctx, nil
}

func EncUpdate(ctx int64, in []byte) (out []byte, err error){
	fmt.Println("call encUpdate")
	out = []byte("test update result")
	return out, nil
}

func EncFinal(ctx int64)(out []byte, err error){
	fmt.Println("call encFinal")
	return []byte("test final result"), nil
}

func Release(ctx int64)(err error){
	fmt.Println("call release")
	return nil
}

