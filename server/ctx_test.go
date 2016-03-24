package server

import(
	"testing"
	"fmt"
	//"sync"
	"math/rand"
	"github.com/pborman/uuid"
)

func GetUUid() string{
	id  := uuid.New()
	return id
	//return "aaa-bbbbbbbbbb-ccccccccc"
}

// func TestDuplicGenUUid(t *testing.T){
// 	uuidMap := make(map[string]bool, 1024)
// 	for i:=0; i<1024; i++ {
// 		id := uuid.New()
// 		if uuidMap[id] {
// 			t.Error("get duplic uuid ", id)
// 		}
// 		uuidMap[id] = true
// 		t.Log(id)
// 	}
// }

func TestAddRemoveKeyMap(t *testing.T){
	uuidMap := make(map[string]bool, 1) //初始大小为0
	t.Logf("map len=%d\n", len(uuidMap))
	id := uuid.New()

	uuidMap[id] = true

	delete(uuidMap, id)
	t.Log("delete map key", id)

	id2 := uuid.New()
	uuidMap[id2] = true
	t.Log("add map key", id2)

	id3 := uuid.New()
	uuidMap[id3] = true
	t.Logf("map len=%d", len(uuidMap))
	if _, ok := uuidMap[id3] ; ok == true {
		t.Error("should no space for id", id3)
	}else {
		t.Log("map has no space for new key to add")
	}
}

func TestGenCtxKeyMap(t *testing.T){
	t.Log("test gen uuid => *int map")
	ctxKeyMap := make(map[string](*int), 100)

	uuid := GetUUid()
	t.Log("uuid=", uuid)
	ctx := int(10)
	t.Logf("ctx pointer value=%d\n", &ctx)
	ctxKeyMap[uuid] = &ctx

	getCtx := ctxKeyMap[uuid]
	if getCtx != &ctx {
		t.Error("notfound ctx")
	}else{
		t.Log("get ctx by uuid", *getCtx)
	}

}

func TestGenCtxKey(t *testing.T){
	t.Log("test gen ctx key map")
	ctxKeyMap := make(map[string]int64, 10)
	
	rand.Seed(100)
	ctx := rand.Int63() //int64(1000)
	key := fmt.Sprintf("%x", ctx)
	t.Logf("key value %s", key)

	ctxKeyMap[key] = ctx

	ctx1 := rand.Int63()
	key1 := fmt.Sprintf("%x", ctx1)
	t.Logf("key1 value %s", key1)

	ctxKeyMap[key1] = ctx1	

	if ctxKeyMap["a"] != 0{
		t.Error("notfound map value is not 0")
	}else{
		t.Log("notfound map value is 0")
	}

	if ctxKeyMap[key] != ctx {
		t.Errorf("key is not in map, get ctx value=%d", ctxKeyMap[key])
	}else {
		t.Log("test ok...")
	}
	
}

func TestMapFindError(t *testing.T){

	emptyMap := make(map[string]string, 10)

	if emptyMap["a"] == "" {
		t.Log("map is nil")
	}else{
		t.Error("map is not nil")
	}
}