package main

import (
	"fmt"
	"net"
	"net/rpc/jsonrpc"
	"../key_value_storage"
	"flag"
	"os"
	"bufio"
	"strings"
)

const BaseNameMethod = "KeyValue"

const (
	RedColor    = "\x1b[31m"
	GreenColor  = "\x1b[32m"
	YellowColor = "\x1b[33m"
	BlueColor   = "\x1b[34m"
	ResetColor  = "\x1b[0m"
)

// Приветсвенный текст
func hellowText() {
	fmt.Printf("%vHi All!\nIt's a simple client for Simple In-Memort Key-Value storage based on Golang!\n%v", YellowColor, ResetColor)
}

// Текст помощи, команда help
func helpText() {
	fmt.Printf("Availible command: %vSet, %vGet, %vDelete, %vChange{imkv, mukv}%v!\n", GreenColor, BlueColor, RedColor, YellowColor, ResetColor)
	fmt.Printf("Use the following syntax:%v METHOD_NAME[SPACE]ARG1[SPACE]ARG2 %v\n", GreenColor, ResetColor)
	fmt.Printf("Example: %vget key1%v", BlueColor, ResetColor)
}

func main() {
	// Разбор аргументов командной строки
	IP := flag.String("addr", "127.0.0.1", "Addr default localHost 127.0.0.1 ")
	PORT := flag.Int("port", 9008, "Port number default = 9008")
	flag.Parse()

	// коннектимся к серверу
	client, err := net.Dial("tcp", fmt.Sprintf("%v:%v", *IP, *PORT))
	if err != nil {
		fmt.Printf("Dialing problem: %v\n", err)
		os.Exit(1)
	}

	// Создаем новый экземпляр клиент JSON-RPC
	RpcClient := jsonrpc.NewClient(client)
	var reply key_value_storage.Reply // для результатов (ответов)
	var args key_value_storage.Args   // для передачи параметров

	hellowText()
	helpText()

	for { // бесконечный цикл опроса пользователя, и отправления команд на сервер при правильных введенных командах

		fmt.Printf("\n# ")

		// Чтение байтов с консоли и разбор на части (делиметер пробел)
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		strings.ToLower(text)
		sliceText := strings.Split(text, " ")
		n := len(sliceText)

		// Команда помощи
		if sliceText[0] == "help\n" || sliceText[0] == "help" {
			helpText()
			continue
		}

		// Незатейливые проверки
		if n < 2 {
			fmt.Printf("%vToo low param! %v %v", RedColor, n, ResetColor)
			continue
		}

		if n > 3 {
			fmt.Printf("%vToo many param! %v %v", RedColor, n, ResetColor)
			continue
		}

		// Предварительное создание арщументов для RPC клиента
		method := ""
		args = key_value_storage.Args{}

		// Анализ метода
		switch sliceText[0] {
		case "set":
			if n != 3 {
				fmt.Printf("%vBad number param! %v %v", RedColor, n, ResetColor)
				continue
			} else {
				method = "Set"
				args.Key = sliceText[1]
				args.Data = strings.Split(sliceText[2], "\n")[0] // убираем символ \n
			}
		case "get":
			fallthrough
		case "delete":
			method = strings.ToUpper(string(sliceText[0][0])) + sliceText[0][1:] // делаем первую букву большой
			args.Key = strings.Split(sliceText[1], "\n")[0]                      // убираем символ \n
		case "change":
			fallthrough
		case "Change":
			method = "ChangeBackEnd"
			args.Key = strings.Split(sliceText[1], "\n")[0]
			//args.Data = key_value_storage.IMKV{safemap.New(1)}
		default:
			fmt.Printf("%vBad mathod! %v", RedColor, ResetColor)
			continue
		}

		fmt.Printf("RPC call\nMethod: %v\nParams:[%v]\n", method, args.ToString())

		// Synchronous call, висим ждем ответа от сервера
		if RpcClient.Call(BaseNameMethod+"."+method, args, &reply); err != nil {
			fmt.Printf("Receive problem: %v\n", err)
			//os.Exit(2)
		}
		fmt.Printf("Result: %v \n", reply.ToString())
	}
}
