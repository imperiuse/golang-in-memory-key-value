package key_value_storage

// Собственная реализация Безопасной Map
import (
	"../safemap"
	"fmt"
	"errors"
	"sync"
)

// Структура описывающая входные параметры для перадачи в RPC server
type Args struct {
	Key  string
	Data interface{}
}

// Структура описывающая выходные параметры получаемые от RPC server
type Reply struct {
	Data    interface{} // результат - данные
	ErrNo   int         // код ошибки
	ErrDesc string      // описание ошибки
}

// Структура ошибок
type KVError struct {
	Err     error  // ошибка
	ErrCode int    // код ошибки
	ErrDesc string // описание ошибки
}

const (
	NoErr       = iota
	PANIC
	NotFoundKey
	ServerError
)

// ErrCode
// 0 - NO Err, GOOD
//
// 1 - Panic
// 2 - Not found key (exactly it's not a error)
// 3 - Server Error  (Bad interface)

// Интерфейс "Хранилище", его должны реализовывать конретные реализации хранилиш, согласно требованию ТЗ
// " - бэкенды kv должны быть заменяемы, т.е. мы должны иметь возможность с минимальными трудозатратами добавить
// любой другой бэкенд, помимо in-memory "  прим. -er - cуффикс интерфейсов, поэтому такое слово Storager
type Storager interface {
	Set(string, interface{}) *KVError
	Get(string) (interface{}, *KVError)
	Delete(string) *KVError
}

// Главный тип структура Хранилище Ключ-Значение, содержит в себе конкретную реализацию хранилища, удовлетв. интерфейсу. Storager
type KeyValue struct {
	storage Storager // Storager it's Interface ( so have methods Set(key, value), Get(key), Delete(key))
}

// Конструктор
//  @param
//		    StorageBackEnd     interface{}   - "back-end" должен удовлетворять интерфесу Storager
//  @return
//          KeyValue   -  ссылка на новый объект Хранилище KeyValue
//      	error      -  признак ошибки
func CreateKeyValueStorage(storageBackEnd Storager) (*KeyValue, error){
	if storageBackEnd != nil{
		kv := KeyValue{storageBackEnd}
		return &kv, nil
	}else{
		return nil, errors.New("Nil Back-End interface")
	}
	return nil, nil
}


// Метод смены Back-End Key Value
//  @param
//		    storageInterface     interface{}   - новый "back-end" должен удовлетворять интерфесу Storager
//  @return
//      	error  -  признак ошибки
func (kv *KeyValue) changeBackEnd(newStorageBackEnd Storager) error{
	if newStorageBackEnd != nil{
		kv.storage = newStorageBackEnd
	}else{
		return errors.New("Nil Back-End interface")
	}
	return nil
}

// Метод смены Back-End Key Value RPC
//  @param
//     args       *Args   входной параметр - переданные данные запроса (key, value)
//     reply      *Reply  выходной параметр - результат - пустая структура или errDesc опис. ошибки
//  @return
//     error
func (kv *KeyValue) ChangeBackEnd(args *Args, reply *Reply) error{
	if args.Key == "imkv" {
		kv.changeBackEnd(&IMKV{safemap.New(1)})
	}else if args.Key == "mukv" {
		kv.changeBackEnd(&MUKV{make(map[string]interface{},0), *new(sync.Mutex)})
	}else if args.Key == "obj" { // Failed Test send obj over RPC
		var storager Storager = args.Data.(Storager)
		kv.changeBackEnd(storager)
	}else{
		reply.ErrDesc = "BAD COMMAND"
		reply.ErrNo = ServerError
	}
	return nil
}

// Метод записи новой записи типа ключ-значение в хранилище, возращает ошибку
//  @param
//     args       *Args   входной параметр - переданные данные запроса (key, value)
//     reply      *Reply  выходной параметр - результат - пустая структура или errDesc опис. ошибки
//  @return
//                error    ошибка
func (kv *KeyValue) Set(args *Args, reply *Reply) error {
	defer recoveryFunc("(*KeyValue) Set()", "may be interface cast")
	fmt.Printf("\nExecuting Method: %v; ARGS: %v\n", "Set", args.ToString())
	if KVE := kv.storage.Set(args.Key, args.Data); KVE == nil {
		//no error
		reply.ErrNo = NoErr
		reply.ErrDesc = ""
	} else {
		// error was
		reply.ErrNo = KVE.ErrCode
		reply.ErrDesc = KVE.ErrDesc
	}
	fmt.Printf("Result %v\n", reply.ToString())
	return nil
}

// Метод получения значения по ключу из хранилища, возращает ошибку
//  @param
//     args       *Args   входной параметр - переданные данные запроса (key)
//     reply      *Reply  выходной параметр - результат - запрашиваемые данные, errDesc - "",  или errDesc опис. ошибки
//  @return
//                error    ошибка
func (kv *KeyValue) Get(args *Args, reply *Reply) error {
	defer recoveryFunc("(*KeyValue) Get()", "may be interface cast")
	fmt.Printf("\nExecuting Method: %v; ARGS:%v\n", "Get", args.ToString())
	if data, KVE := kv.storage.Get(args.Key); KVE == nil {
		//no error
		reply.ErrNo = NoErr
		reply.ErrDesc = ""
		reply.Data = data
	} else {
		// error was
		reply.ErrNo = KVE.ErrCode
		reply.ErrDesc = KVE.ErrDesc
	}
	fmt.Printf("Result %v\n", reply.ToString())
	return nil
}

// Метод удаления пары ключ-значения из хранилище, возращает ошибку
//  @param
//     args       *Args   входной параметр - переданные данные запроса (key)
//     reply      *Reply  выходной параметр - результат - запрашиваемые данные, errDesc - "",  или errDesc опис. ошибки
//  @return
//                error    ошибка
func (kv *KeyValue) Delete(args *Args, reply *Reply) error {
	defer recoveryFunc("(*KeyValue) Delete()", "may be interface cast")
	fmt.Printf("\nExecuting Method: %v; ARGS: %v\n", "Del", args.ToString())
	if KVE := kv.storage.Delete(args.Key); KVE == nil {
		//no error
		reply.ErrNo = NoErr
		reply.ErrDesc = ""
	} else {
		// error was
		reply.ErrNo = KVE.ErrCode
		reply.ErrDesc = KVE.ErrDesc
	}
	fmt.Printf("Result %v\n", reply.ToString())
	return nil
}

// Конкретная реализация "Хранилища" : "In-memory Key-Value"
type IMKV struct {
	safemap.SafeMap // safe map
}

// Метод запись новой пары ключ-значение, возращает указатель на структуру ошибки
//  @param
//     key       string     входной параметр - ключ
//     value     interface  входной параметр - значение
//  @return
//     err      *KVError   - (nil - все хорошо)
func (s *IMKV) Set(key string, data interface{}) (err *KVError) {
	defer recoveryFuncErr("Set()", "smth bad s.Set()", err)
	s.SafeMap.Set(key, data)
	return nil
}

// Метод запись получение значение, возращает значение и указатель на структуру ошибки
//  @param
//     key       string     входной параметр - ключ
//  @return
//     data     interface{}  - данные хранящиеся по ключу в случае успеха
//     err      *KVError     - (nil - все хорошо)
func (s *IMKV) Get(key string) (data interface{}, err *KVError) {
	defer recoveryFuncErr("Get()", "smth bad in s.Get(key)", err)
	var found bool
	if data, found = s.SafeMap.Get(key); !found {
		return nil, &KVError{nil, NotFoundKey, "Not found Key"}
	}
	return
}

// Метод запись удаление пары ключ-значение, возращает указатель на структуру ошибки
//  @param
//     key       string     входной параметр - ключ
//  @return
//     err      *KVError   - (nil - все хорошо)
func (s *IMKV) Delete(key string) (err *KVError) {
	defer recoveryFuncErr("Delete()", "smth bad in s.Del(key)", err)
	s.SafeMap.Del(key)
	return
}

// Pretty print error info
func (KVE *KVError) ToString() (s string) {
	s = fmt.Sprintf(""+
		"ErrCode:     %v \n "+
		"Err:         %v \n "+
		"Description: %v \n ",
		KVE.ErrCode, KVE.Err, KVE.ErrDesc)
	return
}

// Pretty print Args
func (a *Args) ToString() (s string) {
	s = fmt.Sprintf("{Key: %v, Data:%v}", a.Key, a.Data)
	return
}

// Pretty print Reply info
func (r *Reply) ToString() (s string) {
	s = fmt.Sprintf("{"+
		"Data: %v, "+
		"ErrCode: %v, "+
		"Description: %v}",
		r.Data, r.ErrNo, r.ErrDesc)
	return
}

// Recover функция, для проверки возникновения паники в контролируемой функции
func recoveryFuncErr(f string, reason string, KVE *KVError) {
	if r := recover(); r != nil {
		fmt.Printf("Recovery_func() in %v detect PANIC!. Reason: %v; Err: %v", f, reason, r)
		KVE.Err = errors.New(fmt.Sprintf("%v", r))
		KVE.ErrCode = PANIC // PANIC!
		KVE.ErrDesc = fmt.Sprintf("Panic! was at %v. %v", f, r)
	}
	return
}

// Recover функция, для проверки возникновения паники в контролируемой функции
func recoveryFunc(f string, reason string) {
	if r := recover(); r != nil {
		fmt.Printf("Recovery_func() in %v detect PANIC!. Reason: %v; Err: %v", f, reason, r)
	}
	return
}


// Конкретная реализация "Хранилища" : "Простая мапа с мьютексом"
type MUKV struct {
	m map[string]interface{}
	mu sync.Mutex
}

func CreateMUKV()MUKV{
	return MUKV{make(map[string]interface{},0), *new(sync.Mutex)}
}

func (s *MUKV) Set(key string, data interface{}) (err *KVError) {
	defer recoveryFuncErr("Set()", "smth bad s.Set()", err)
	s.mu.Lock()
	s.m[key] = data
	defer s.mu.Unlock()
	return nil
}


func (s *MUKV) Get(key string) (data interface{}, err *KVError) {
	defer recoveryFuncErr("Get()", "smth bad in s.Get(key)", err)
	s.mu.Lock()
	var found bool
	if data, found = s.m[key]; !found {
		return nil, &KVError{nil, NotFoundKey, "Not found Key"}
	}
	defer s.mu.Unlock()
	return
}

func (s *MUKV) Delete(key string) (err *KVError) {
	defer recoveryFuncErr("Delete()", "smth bad in s.Del(key)", err)
	s.mu.Lock()
	delete(s.m, key)
	defer s.mu.Unlock()
	return
}