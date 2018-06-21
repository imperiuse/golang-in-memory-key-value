// SERVER in memory key-value based on JSON-RPC v1.0.
// https://ru.wikipedia.org/wiki/JSON-RPC

package main

import (
	"net/rpc"
	"net"
	"net/rpc/jsonrpc"
	"./kv_storage"
	"./safemap"
	"flag"
	"io/ioutil"
	"os"
	"encoding/json"
	"fmt"
	"./pidfile"
	"os/signal"
	"syscall"
)

// Красота
const (
	RED_COLOR    = "\x1b[31m"
	GREEN_COLOR  = "\x1b[32m"
	YELLOW_COLOR = "\x1b[33m"
	BLUE_COLOR   = "\x1b[34m"
	RESET_COLOR  = "\x1b[0m"
)

// Функция для красивой проверки и отображения ошибки.
func CheckErrorFunc(err error, f string) {
	if err != nil {
		fmt.Printf("[CheckErr %v] %v%v%v%v", f, RED_COLOR, "[Error!]\n", err, RESET_COLOR)
	} else {
		fmt.Printf("[CheckErr %v] %v%v%v", f, GREEN_COLOR, "[Successful!]\n", RESET_COLOR)
	}
}

// Структура настроек, по большей части показать умение работать с JSON файлами
type SettingsT struct {
	Addr    string `json:"addr"` // немного магии для распарсивания полей в нижнем регистре
	Port    int    `json:"port"`
	Pidfile string `json:"pidfile"`
}

var Settings SettingsT

// Обработчик событий сигналов
// @param
//    signals  chan os.Signal - канал системных сигналов для данной утилиты
func SignalHandler(signals chan os.Signal) {
	//defer func() { EXIT <- ReasonExit{2, "SignalHandler()", "SignalHandler unexpectedly finished!"} }()
	for {
		sig := <-signals
		fmt.Printf("SIGNAL_HANDLER.  Receive a signal: %v", sig.String())
		switch sig {
		//case syscall.SIGUSR1: //0x0A == 10
		//case syscall.SIGUSR2: //0x0C == 12
		case syscall.SIGINT: //0x02 == 2
			fallthrough
		case syscall.SIGTERM: //0x0f == 15
			fmt.Printf("SIGNAL_HANDLER. Shutdown Server")
			pidfile.Remove(Settings.Pidfile)
			os.Exit(7)
			//EXIT <- ReasonExit{1, "SignalHandler()", fmt.Sprintf("Signal %v catch!", sig)}
		default:
			fmt.Printf("SIGNAL_HANDLER.  Receive a UNKNOWN signal: %v", sig.String())
		}
	}
}

func main() {
	// Разбор аргументов командной строки
	pathConfig := flag.String("json", "./settings.json", "Path to json config file")
	flag.Parse()

	// Загрузка из файла settings
	fileSettings, err := ioutil.ReadFile(*pathConfig)
	CheckErrorFunc(err, "ReadFile()")
	if err != nil {
		os.Exit(1)
	}
	err = json.Unmarshal(fileSettings, &Settings)
	if err != nil {
		fmt.Printf("[Settings unmarshall from %v%v%v] [%vError!%v]\n", BLUE_COLOR, *pathConfig, RESET_COLOR, RED_COLOR, RESET_COLOR)
		os.Exit(2)
	} else {
		fmt.Printf("[Settings unmarshall from %v%v%v] [%vSuccessful!%v]\n", BLUE_COLOR, *pathConfig, RESET_COLOR, GREEN_COLOR, RESET_COLOR)
	}

	//Обработчик сигналов
	signals := make(chan os.Signal, 1)                      // создание канала для приема сигналов
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM) // регистрация канала для сигналов
	go SignalHandler(signals)                               // go-рутина для обработки сигналов

	// PID file create
	if err := pidfile.Write(Settings.Pidfile); err != nil {
		fmt.Printf("%v Can't create pid file! %v %v \n", RED_COLOR, err, RESET_COLOR)
		os.Exit(1)
	}
	defer pidfile.Remove(Settings.Pidfile)

	// Создаем структуру "Хранилища", а также внутри вложенно инициализиурем канал SafeMap
	KeyValueStorage := kv_storage.KeyValue{Storage: kv_storage.IMKV{SM: safemap.New(1)}}

	fmt.Printf("Key-Value Storage created %v[Successful]\n%v", GREEN_COLOR, RESET_COLOR)

	// Создаем новый экземпляр RPC сервера
	server := rpc.NewServer()
	server.Register(&KeyValueStorage) // Регистрация методов (все методы структуры экзмеляра KeyValueStorage)
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	// Вешаем лисенера
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", Settings.Addr, Settings.Port))
	if err != nil {
		fmt.Printf("Listening error: %v\n", err.Error())
		os.Exit(3)
	}

	fmt.Printf("Create listener at %v %v[Successful]\n%v",
		fmt.Sprintf("%v:%v", Settings.Addr, Settings.Port), GREEN_COLOR, RESET_COLOR)

	for {
		// Ждем подключения клиента
		fmt.Printf("Waiting clients ... \n")
		if conn, err := listener.Accept(); err != nil {
			// все плохо
			fmt.Printf("Accept error: %v\n", err.Error())
			os.Exit(4)
		} else { // все хорошо
			fmt.Printf("%vNew connection established! %v\n", GREEN_COLOR, RESET_COLOR)
			go server.ServeCodec(jsonrpc.NewServerCodec(conn)) // вызов зарегистрированного метода, который хочет клиент
		}
	}
}
