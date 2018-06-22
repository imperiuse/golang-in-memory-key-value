package kv_storage

import (
	"testing"
	"../safemap"
	"fmt"
)

func TestVerySimpleAbstractKV(t *testing.T) {
	fmt.Printf("\nTestVerySimpleAbstractKV\n")
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
	} else {
		if pReply.ErrNo != NOT_FOUND_KEY {
			t.Errorf("Error test value != empty interface{}! %v", pReply)
		}
	}
}

func TestVerySimpleIMKV(t *testing.T) {
	fmt.Printf("\nTestVerySimpleIMKV\n")
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
	if _, PKVE := KeyValueStorage.Storage.Get("key"); PKVE == nil {
		t.Errorf("Not err about Not Found Key!")
	} else {
		if PKVE.ErrCode != NOT_FOUND_KEY {
			t.Error("Bad error code found!", PKVE.ErrCode)
		}
	}
}
