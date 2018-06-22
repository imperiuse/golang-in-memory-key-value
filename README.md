# golang-in-memory-key-value
Test task


Инфраструктура разработки
    OS: Ubunta 16.04 LTS & Ubunta Server 16.10 LTS (личный администрируемый сервер)
    Lang: go 1.10.3 linux/amd64
    IDE:  Goland 2018.1

Стандартная библиотеки:
    net
    net/rpc
    net/rpc/jsonrpc
    json

Собствнные мини-библиотеки:
    safemap
    pidfile

# Server start

    go build
    ./server

Настройки сервера лежат в settings.json
Путь к нему можно указать с помощью флага path, по-умолчанию путь ``./setting.json`

# Start client

    cd client
    go build
    ./client

Настройки клиента: адресс и порт сервера куда подключаться можно задать через флаги `addr`, `port`
По умолчанию клиент и сервер настроены для быстрого и удобного запуска на одной машине, порт по-умолчанию 9008, адрес естесвенно localhost


