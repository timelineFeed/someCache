package example

import (
	"fmt"
	cache "github.com/timelineFeed/someCache"
	"reflect"

	"testing"
)

func TestCache(t *testing.T) {
	c, err := cache.Init(cache.OCacheModel)
	if err != nil {
		t.Errorf("cache Init err=%+v\n", err)
	}
	key := "key"
	value := "value"
	err = c.SetCache(key, value)
	if err != nil {
		t.Errorf("set cache err=%+v\n", err)
	}
	v, err := c.GetCache(key)
	if err != nil {
		t.Errorf("get cache err=%+v\n", err)
	}
	if !reflect.DeepEqual(v, value) {
		t.Errorf("get:%+v,want:%+v\n", v, value)
	}
}

func TestTTL(t *testing.T) {
	c, err := cache.Init(cache.TTLCacheModel)
	if err != nil {
		t.Errorf("cache Init err=%+v\n", err)
	}
	key := "key"
	value := "value"
	err = c.SetCache(key, value)
	if err != nil {
		t.Errorf("set cache err=%+v\n", err)
	}
	v, err := c.GetCache(key)
	if err != nil {
		t.Errorf("get cache err=%+v\n", err)
	}
	if !reflect.DeepEqual(v, value) {
		t.Errorf("get:%+v,want:%+v\n", v, value)
	}
}

func TestDTTL(t *testing.T) {
	c, err := cache.Init(cache.DTTLCacheModel)
	if err != nil {
		t.Errorf("cache Init err=%+v\n", err)
	}
	key := "key"
	value := "value"
	err = c.SetCache(key, value)
	if err != nil {
		t.Errorf("set cache err=%+v\n", err)
	}
	fmt.Println("--------------------")
	v, err := c.GetCache(key)
	if err != nil {
		t.Errorf("get cache err=%+v\n", err)
	}
	fmt.Printf("v:%+v,t:%T\n", v, v)
	if !reflect.DeepEqual(v, value) {
		t.Errorf("get:%+v,want:%+v\n", v, value)
	}
}
