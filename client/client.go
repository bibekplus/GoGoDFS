
package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type MinionAddr struct{
	Host string
	Port int
}
type Block struct{
	BlockId string
	Minions []MinionAddr
}

type FileBlock struct{
	Name string
	Size int
}

func main() {

	var reply []Block

	// Create a TCP connection to localhost on port 8000
	client, err := rpc.DialHTTP("tcp", "localhost:8000")
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	file := FileBlock{"/abc.txt", 2000}
	client.Call("Master.Write", file , &reply)
	fmt.Println(reply)
	fmt.Println("_________________________________\n")
	var reply2 []Block
	client.Call("Master.Read", "/abc.txt", &reply2)
	fmt.Println(reply2)
}