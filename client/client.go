
package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
)

const BLOCK_SIZE = 100

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

func check(e error){
	if e!= nil{
		panic(e)
	}
}

//
//
//func get(client *rpc.Client, fileName string) ([]Block, error){
//	var fileBlocks []Block
//	xType := reflect.TypeOf(client)
//	fmt.Println(xType)
//	err := client.Call("Master.Read", fileName, &fileBlocks)
//	if err != nil{
//		fmt.Println("File not found")
//		return  nil, err
//	}
//
//	for _, blck := range fileBlocks{
//		//id := blck.BlockId
//		minions:= blck.Minions
//		for _, min := range minions{
//			minAddress := fmt.Sprintf("%s:%d", min.Host, min.Port)
//			conn, err := rpc.DialHTTP("tcp", minAddress)
//			if err != nil {
//				log.Fatal("Connection error: ", err)
//			}
//			//client.Call("Minion.get", blck.BlockId , &reply)
//		fmt.Println(conn)
//		}
//	}
//	fmt.Println("-------------------")
//	return fileBlocks, err
//}


func put(client *rpc.Client, fileSource string, fileName string) error {
	info, err := os.Stat(fileSource)
	if err != nil{
		return err
	}
	size := int(info.Size())
	f, err := os.Open(fileSource)
	check(err)
	file := FileBlock{fileName, size}
	var fileBlocks []Block
	client.Call("Master.Write", file , &fileBlocks)
	for _, blck := range fileBlocks {
		data := make([]byte, BLOCK_SIZE)
		_, err := f.Read(data)
		check(err)
		id := blck.BlockId
		minion := blck.Minions[0]
		minions := blck.Minions[1:]

		message := MessageToMinion{id, data, minions}
		reply := 0
		minAddress := fmt.Sprintf("%s:%d", minion.Host, minion.Port)
		fmt.Println(minAddress)
		conn, err := rpc.DialHTTP("tcp", minAddress)
		fmt.Println(err)
		err = conn.Call("Minion.Put", message, &reply)
		fmt.Println(reply)
		if err != nil{
			return err
		}

	}


	return err
}


func main() {

	//var reply []Block

	// Create a TCP connection to localhost on port 8000
	client, err := rpc.DialHTTP("tcp", "localhost:8000")
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	if os.Args[1] == "put"{
		put(client, os.Args[2], os.Args[3])
	}
	//file := FileBlock{"/abc.txt", 2000}
	//client.Call("Master.Write", file , &reply)

	//xType := reflect.TypeOf(client)
	//fmt.Println(xType)

	//fmt.Println(reply)
	fmt.Println("_________________________________\n")
	//var reply2 []Block
	//err = client.Call("Master.Read", "/abc.txt", &reply2)
	//reply2, err = get(client, "/abc.txt")
	//if err != nil {
	//	log.Fatal("Error: ", err)
	//}
	//fmt.Println(reply2)

	//err = put(client, "./a.txt", "a.txt")
}