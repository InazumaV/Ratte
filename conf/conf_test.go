package conf

import (
	"fmt"
	"testing"
)

func TestConf_Load_Local(t *testing.T) {
	c := New("./config.json5")
	err := c.Load(nil)
	if err != nil {
		t.Error(err)
	}
	t.Log(c)
}

func TestConf_Load_Remote(t *testing.T) {
	c := New("http://127.0.0.1:9000/config.json5")
	err := c.Load(nil)
	if err != nil {
		t.Error(err)
	}
}

func TestConf_Watch(t *testing.T) {
	c := New("./config.json5")
	err := c.Load(nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(c)
	c.SetEventHandler(func(event uint, target ...string) {
		switch event {
		case ConfigFileChangedEvent:
			t.Log("Event:", "ConfigFileChangedEvent", "target:", target)
		case CoreDataPathChangedEvent:
			t.Log("Event:", "CoreDataPathChangedEvent", "target:", target)
		}
	})
	c.SetErrorHandler(func(err error) {
		t.Error(err)
	})
	err = c.Watch()
	if err != nil {
		t.Error(err)
	}
	t.Log("press any key to done.")
	fmt.Scan()
}
