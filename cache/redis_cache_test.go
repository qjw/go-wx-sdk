package cache

import (
	"testing"
	"time"
)

func TestMemcache(t *testing.T) {
	redis,_ := NewRedisCache("localhost:6379","",0)
	var err error
	timeoutDuration := 10 * time.Second
	if err = redis.Set("username", "silenceper", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	if !redis.IsExist("username") {
		t.Error("IsExist Error")
	}

	name := redis.Get("username").(string)
	if name != "silenceper" {
		t.Error("get Error")
	}

	if err = redis.Delete("username"); err != nil {
		t.Errorf("delete Error , err=%v", err)
	}
}
