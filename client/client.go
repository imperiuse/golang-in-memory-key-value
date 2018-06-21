package main

import (
	"fmt"
	"net"
	"net/rpc/jsonrpc"
	"../kv_storage"
	"flag"
	"os"
	"bufio"
	"strings"
)

const BASE_PART_NAME_METHOD = "KeyValue"

const (
	RED_COLOR    = "\x1b[31m"
	GREEN_COLOR  = "\x1b[32m"
	YELLOW_COLOR = "\x1b[33m"
	BLUE_COLOR   = "\x1b[34m"
	RESET_COLOR  = "\x1b[0m"
)

func hellowText() {
	fmt.Printf("%vHi All!\nIt's a simple client for Simple In-Memort Key-Value storage based on Golang!\n%v", YELLOW_COLOR, RESET_COLOR)
}

func helpText() {
	fmt.Printf("Availible command: %vSet, %vGet, %vDelete%v!\n", GREEN_COLOR, BLUE_COLOR, RED_COLOR, RESET_COLOR)
	fmt.Printf("Use the following syntax:%v METHOD_NAME[SPACE]ARG1[SPACE]ARG2 %v\n", GREEN_COLOR, RESET_COLOR)
	fmt.Printf("Example: %vget key1%v", BLUE_COLOR, RESET_COLOR)
}

func main() {
	// Разбор аргументов командной строки
	IP := flag.String("addr", "127.0.0.1", "Addr default localHost 127.0.0.1 ")
	PORT := flag.Int("port", 9008, "Port number default = 9008")
	flag.Parse()

	// коннектимся
	client, err := net.Dial("tcp", fmt.Sprintf("%v:%v", *IP, *PORT))
	if err != nil {
		fmt.Printf("Dialing problem: %v\n", err)
		os.Exit(1)
	}

	// Создаем новый экземпляр клиент JSON-RPC
	RpcClient := jsonrpc.NewClient(client)
	var reply kv_storage.Reply // для результатов (ответов)
	var args kv_storage.Args   // для передачи параметров
	hellowText()
	helpText()
	for {

		fmt.Printf("\n# ")

		// Чтение байтов с консоли для отправки
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		strings.ToLower(text)
		s_mas := strings.Split(text, " ")
		n := len(s_mas)

		if s_mas[0] == "help\n" || s_mas[0] == "help" {
			helpText()
			continue
		}

		if n < 2 {
			fmt.Printf("%vToo low param! %v %v", RED_COLOR, n, RESET_COLOR)
			continue
		}

		if n > 3 {
			fmt.Printf("%vToo many param! %v %v", RED_COLOR, n, RESET_COLOR)
			continue
		}

		method := ""
		args = kv_storage.Args{"", ""}
		switch s_mas[0] {
		case "set":
			if n != 3 {
				fmt.Printf("%vBad number param! %v %v", RED_COLOR, n, RESET_COLOR)
				continue
			} else {
				method = "Set"
				args.Key = s_mas[1]
				args.Data = strings.Split(s_mas[2], "\n")[0]
			}
		case "get":
			fallthrough
		case "delete":
			method = strings.ToUpper(string(s_mas[0][0])) + s_mas[0][1:]
			args.Key = strings.Split(s_mas[1], "\n")[0]
		default:
			fmt.Printf("%vBad mathod! %v", RED_COLOR, RESET_COLOR)
			continue
		}

		fmt.Printf("Method: %v\nParams:[%v]\n", method, args.ToString())

		// Synchronous call, висим ждем ответа от сервера
		if RpcClient.Call(BASE_PART_NAME_METHOD+"."+method, args, &reply); err != nil {
			fmt.Printf("Receive problem: %v\n", err)
			//os.Exit(2)
		}
		fmt.Printf("Result: %v \n", reply)
	}
}
