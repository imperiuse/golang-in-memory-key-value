# golang-in-memory-key-value

### Develop infrastructures:

    OS:    Ubunta 18.04 LTS & Ubunta Server 18.04 LTS
    Lang:  go 1.10.3 linux/amd64
    IDE:   Goland 2018.2 / VSCode

### Standard Library used:

    net
    net/rpc
    net/rpc/jsonrpc
    json

### Custom library used:

    "github.com/imperiuse/golang_lib"

    "./safemap"
    "./pidfile"

### Command to server start:

    go build && ./server

#### Some addition server info:

Server's settings here: `settings.json`
You can choose own settings file by flag -json  ` -json ./setting.json`

Default server port: `9008`


### Command to client start:

    cd client
    go build && ./client

Default destinations addr:port:  `localhost:9008`

#### Some addition client info:
You can choose own Addr and port setting by flags: `addr` and `port`



