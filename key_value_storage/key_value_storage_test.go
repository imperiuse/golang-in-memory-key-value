package key_value_storage

import (
	"testing"
	"../safemap"
	"fmt"
)

func TestVerySimpleAbstractKV(t *testing.T) {
	fmt.Printf("\nTestVerySimpleAbstractKV\n")
	KeyValueStorage, _ := CreateKeyValueStorage(&(IMKV{safemap.New(1)}))
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
		if pReply.ErrNo != NotFoundKey {
			t.Errorf("Error test value != empty interface{}! %v", pReply)
		}
	}
}

func TestVerySimpleIMKV(t *testing.T) {
	fmt.Printf("\nTestVerySimpleIMKV\n")
	Storage := IMKV{safemap.New(1)}

	if PKVE := Storage.Set("key", "value"); PKVE != nil {
		t.Errorf("Error not nil! %v", PKVE.ToString())
	}
	if value, PKVE := Storage.Get("key"); PKVE != nil {
		t.Errorf("Error not nil! %v", PKVE.ToString())
	} else {
		if value != "value" {
			t.Errorf("Error test value != 'value'! %v", value)
		}
	}
	if PKVE := Storage.Delete("key"); PKVE != nil {
		t.Errorf("Error not nil! %v", PKVE.ToString())
	}
	if _, PKVE := Storage.Get("key"); PKVE == nil {
		t.Errorf("Not err about Not Found Key!")
	} else {
		if PKVE.ErrCode != NotFoundKey {
			t.Error("Bad error code found!", PKVE.ErrCode)
		}
	}
}

type testPairAR struct {
	args  Args
	reply Reply
}

type testPairKV struct {
	key   string
	value interface{}
}

func TestKeyValue(t *testing.T) {
	fmt.Printf("\nTestKeyValue\n")
	KeyValueStorage, _ := CreateKeyValueStorage(&(IMKV{safemap.New(1)}))

	tests := []testPairAR{{Args{}, Reply{}}, {Args{"key", "value"}, Reply{ErrNo: 0, Data: "value"}},
		{Args{"1", "1"}, Reply{ErrNo: 0, Data: "1"}}, {Args{"key", nil}, Reply{ErrNo: 0, Data: nil}}}

	for n, test := range tests {
		reply := Reply{}
		if err := KeyValueStorage.Get(&test.args, &reply); err != nil {
			t.Errorf("Problem in %v test. In Get func (Get not exist key)", n)
		} else {
			if reply.ErrNo != NotFoundKey {
				t.Errorf("Problem in %v test. In Get func (Get not exist key) ErrCode != NotFound %v", n, reply.ErrNo)
			}
			if reply.Data != nil {
				t.Errorf("Problem in %v test. In Get func (Get not exist key) Data!=nil %v ", n, reply.Data)
			}
		}
		if err := KeyValueStorage.Set(&test.args, &reply); err != nil {
			t.Errorf("Problem in %v test. In Set func (Set not exist key)", n)
		} else {
			if reply.ErrNo != NoErr {
				t.Errorf("Problem in %v test. In Set func (Get not exist key) ErrCode != NoErr %v", n, reply.ErrNo)
			}
			if reply.Data != nil {
				t.Errorf("Problem in %v test. In Set func  Data!=nil %v ", n, reply.Data)
			}
		}
		if err := KeyValueStorage.Get(&test.args, &reply); err != nil {
			t.Errorf("Problem in %v test. In Get func (Get exist key)", n)
		} else {
			if reply.ErrNo != NoErr {
				t.Errorf("Problem in %v test. In Get func (Get exist key) ErrCode != NoErr %v", n, reply.ErrNo)
			}
			if reply.Data != test.reply.Data {
				t.Errorf("Problem in %v test. In Get func  Data %v!=%v  ", n, reply.Data, test.reply.Data)
			}
		}
		if err := KeyValueStorage.Delete(&test.args, &reply); err != nil {
			t.Errorf("Problem in %v test. In Delete func (Delete exist key)", n)
		} else {
			if reply.ErrNo != NoErr {
				t.Errorf("Problem in %v test. In Delete func (Get exist key) ErrCode != NoErr %v", n, reply.ErrNo)
			}
		}
	}
}

func TestIMKV(t *testing.T) {
	fmt.Printf("\nTestIMKV\n")
	Storage := IMKV{safemap.New(1)}

	tests := []testPairKV{{"", ""}, {"", "value"}, {"1", "value"},
		{"my_key", "value"}, {"key", new(interface{})}, {"key2", testPairKV{"1", "2"}}}

	for n, test := range tests {
		if value, PKVE := Storage.Get(test.key); PKVE == nil {
			t.Errorf("Problem in %v test. In Get func (Get not exist key)", n)
		} else {
			if PKVE.ErrCode != NotFoundKey {
				t.Errorf("ErrCode %v != NotFoundKey  ", PKVE.ErrCode)
			}
			if PKVE.Err != nil {
				t.Errorf("Problem in %v test. In Get func (Get not exist key) Err: %v", n, PKVE.Err)
			}
			if value != nil {
				t.Errorf("Problem in %v test. In Get func (Get not exist key) Value not nil %v", n, value)
			}
		}
		if PKVE := Storage.Set(test.key, test.value); PKVE != nil {
			t.Errorf("Problem in %v test: %v in Set func", n, PKVE.ToString())
		}
		if PKVE := Storage.Set(test.key, test.value); PKVE != nil {
			t.Errorf("Problem in %v test: %v in Set func (Set exist key)", n, PKVE.ToString())
		}
		if value, PKVE := Storage.Get(test.key); PKVE != nil {
			t.Errorf("Problem in %v test: %v in Get func (Get exist key)", n, PKVE.ToString())
		} else {
			if value != test.value {
				t.Errorf("value get: %v!= value set: %v", value, test.value)
			}
		}
		if PKVE := Storage.Delete(test.key); PKVE != nil {
			t.Errorf("Problem in %v test: %v in Delete func (Delete exist key))", n, PKVE.ToString())
		}
		// Double delete
		if PKVE := Storage.Delete(test.key); PKVE != nil {
			t.Errorf("Problem in %v test: %v in Delete func (Delete not exist key))", n, PKVE.ToString())
		}
	}
}
