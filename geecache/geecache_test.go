package geecache

import (
	"testing"
	"reflect"
	"fmt"
	"log"
)

func TestGetter(t *testing.T) {
	var my_getter Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil 
	})
	expect := []byte("testkey")
	if v, _ := my_getter.Get("testkey"); !reflect.DeepEqual(v, expect){
		t.Error("callback return value incorrect")
	}
}

var db = map[string]string{
	"tom": "4444",
	"jack": "22",
	"sam": "333",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	gee := NewGroup("scores", 2<<10, GetterFunc(
		func (key string) ([]byte, error) {
			log.Println("simple getter key:", key)
			if v, ok:=db[key]; ok {
				// if _, ok := loadCounts(key); !ok {
				// 	loadCounts[key] = 0
				// }
				loadCounts[key] += 1
				return []byte(v), nil
				
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of key: ", k)
		}
		if _, err := gee.Get(k); err!=nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}

	}

	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknown should be empty but got %s", view)
	}

}

