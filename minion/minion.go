package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
)

type Minion int

type AddrMinion struct{
	Host string
	Port int
}

type Block struct{
	BlockId string
	Minions []AddrMinion
}

type FileBlock struct{
	Name string
	Size int
}

type MessageToMinion struct{
	BlockId string
	Data []byte
	Minions [] AddrMinion
}

var DataDir string =  "/tmp/files/"

func (m * Minion) Put (message MessageToMinion, reply * int) error {
	storePath := filepath.Join(DataDir, string(message.BlockId))
	f, err := os.Create(storePath)
	if err != nil{
		return err
	}
	defer f.Close()

	data := message.Data
	_, err = f.Write(data)
	if err != nil{
		return err
	}

	*reply = 1
	return err
}

func main() {

	PORT := os.Args[1]
	DataDir = os.Args[2]

	if _, err := os.Stat(DataDir);os.IsNotExist(err){
		os.Mkdir(DataDir, 0777)
	}
	minion := new(Minion)
	err := rpc.Register(minion)
	if err != nil {
		log.Fatal("Format of service Task isn't correct. ", err)
	}
	rpc.HandleHTTP()
	serveAddress:= fmt.Sprintf("localhost:%s", PORT)
	listener, e := net.Listen("tcp", serveAddress)
	if e != nil {
		log.Fatal("Listen error: ", e)
	}
	log.Printf("Serving RPC server on port %s", PORT)
	// Start accept incoming HTTP connections
	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("Error serving: ", err)
	}
}

