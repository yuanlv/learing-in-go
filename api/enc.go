package enc

import(
	"fmt"

)

func encCreate(appId string) (ctx int64, err error){
	fmt.Println("call encCreate begin")
	ctx := 123
	
	return ctx, nil
}

func encUpdate(ctx int64, in []byte) (out []byte, err error){
	fmt.Println("call encUpdate")
	out := []byte("test update result")
	return out, nil
}

func encFinal(ctx int64)(out []byte, err error){
	fmt.Println("call encFinal")
	return []byte("test final result"), nil
}

func release(ctx int64)(err error){
	fmt.Println("call release")
	return nil
}

