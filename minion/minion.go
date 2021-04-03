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
	if len(message.Minions) > 0{
		msg := MessageToMinion{message.BlockId, message.Data, message.Minions}
		forward(msg)
	}

	*reply = 1
	return err
}

func (m * Minion) Get (blockId string, block * []byte) error{
	fileSource := filepath.Join(DataDir, string(blockId))
	if _, err := os.Stat(fileSource);os.IsNotExist(err){
		return err
	}
	f, err := os.Open(fileSource)
	var reply []byte
	_, err = f.Read(reply)
	if err != nil{
		return err
	}
	*block = reply
	return err
}

func forward (message MessageToMinion) error {
	nextMin := message.Minions[0]
	minions := message.Minions[1:]
	minAddress := fmt.Sprintf("%s:%d", nextMin.Host, nextMin.Port)
	conn, err := rpc.DialHTTP("tcp", minAddress)
	if err!= nil{
		return err
	}
	reply := 0
	msg := MessageToMinion{message.BlockId, message.Data, minions}
	err = conn.Call("Minion.Put", msg, &reply)
	fmt.Println(reply)
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

