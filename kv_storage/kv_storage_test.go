package kv_storage

import (
	"testing"
	"../safemap"
)

func VerySimpleAbstractKVTest(t *testing.T) {
	KeyValueStorage := KeyValue{Storage: IMKV{SM: safemap.New(1)}}
	pArgs := &Args{Key: "key", Data: "value"}
	pReply := &Reply{}

	if err := KeyValueStorage.Set(pArgs, pReply); err != nil {
		t.Errorf("Error not nil! %v", err.Error())
	}
	if err := KeyValueStorage.Get(pArgs, pReply); err != nil {
		t.Errorf("Error not nil! %v", err.Error())
	} else {
		if pReply.Data != "value" && pReply.ErrNo == 0 {
			t.Errorf("Error test value != 'value'! %v", pReply)
		}
	}
	if err := KeyValueStorage.Delete(pArgs, pReply); err != nil {
		t.Errorf("Error not nil! %v", err.Error())
	}
	if err := KeyValueStorage.Get(pArgs, pReply); err != nil {
		t.Errorf("Error not nil! %v", err.Error())
		if pReply.ErrNo != NOT_FOUND_KEY {
			t.Errorf("Error test value != empty interface{}! %v", pReply)
		}
	}
}

func VerySimpleIMKVTest(t *testing.T) {
	KeyValueStorage := KeyValue{Storage: IMKV{SM: safemap.New(1)}}

	if PKVE := KeyValueStorage.Storage.Set("key", "value"); PKVE != nil {
		t.Errorf("Error not nil! %v", PKVE.ToString())
	}
	if value, PKVE := KeyValueStorage.Storage.Get("key"); PKVE != nil {
		t.Errorf("Error not nil! %v", PKVE.ToString())
	} else {
		if value != "value" {
			t.Errorf("Error test value != 'value'! %v", value)
		}
	}
	if PKVE := KeyValueStorage.Storage.Delete("key"); PKVE != nil {
		t.Errorf("Error not nil! %v", PKVE.ToString())
	}
	if value, PKVE := KeyValueStorage.Storage.Get("key"); PKVE != nil {
		t.Errorf("Error not nil! %v", PKVE.ToString())
		if value != new(interface{}) {
			t.Errorf("Error test value != empty interface{}! %v", value)
		}
	}
}
