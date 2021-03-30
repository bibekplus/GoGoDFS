package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"net/rpc"
)


const BLOCK_SIZE = 100
const REPLICATION_FACTOR = 2

type Master int

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


var fileMap = make(map[string][]string)//Filename to blocks

var blockMinions = make(map[string][]string) //Block to minionIDs

var minions = map[string]MinionAddr{
	"1" : MinionAddr{"120.0.0.1", 9000},
	"2" : MinionAddr{"120.0.0.1", 9001},
	"3" : MinionAddr{"120.0.0.1", 9002},
}

func (m * Master) Read (fileName string, Blocks * []Block) error{
	var returnBlocks []Block
	var err error

	blcks, exists:= fileMap[fileName]

	if !exists{
		err = errors.New("File does not exist\n")
		return err
	}
	//Loop across all the blocks containing a file
	//Get minion ids hosting each blocks
	//Get minions address from for each block
	//Create a new block object for each block and
	// append on an array of Blocks to be returned

	for i:=0; i < len(blcks); i++{
		minionAddrs := make([]MinionAddr, 0, REPLICATION_FACTOR-1)
		blockMin := blockMinions[blcks[i]]
		for j:=0; j< len(blockMin); j++{
			minionAddrs = append(minionAddrs, minions[blockMin[j]])
		}
		tempBlock := Block{blcks[i], minionAddrs}
		returnBlocks = append(returnBlocks, tempBlock)
	}

	*Blocks = returnBlocks

	return err
}


func (m * Master) Write (File FileBlock, Blocks * []Block) error{

	//check if file already exists, delete if does
	//Probably need to let the user decide
	_, ok := fileMap[File.Name]
	if ok {
		delete(fileMap, File.Name)
	}

	numOfBlocks := int(math.Ceil(float64(File.Size / BLOCK_SIZE)))


	tempBlocks, err := allocateBlocks(File, numOfBlocks)
	if err!= nil{
		panic(err)
	}

	*Blocks = tempBlocks

	return err

}

func allocateBlocks(file FileBlock, numOfBlocks int) ([]Block, error){
	var returnBlocks []Block
	var err error
	for i:= 0; i < numOfBlocks; i++ {
		blockId, err := newUUID()
		if err!= nil{
			panic(err)
		}


		minionKeys := make([]string, len(minions))

		i := 0
		for k := range minions {
			minionKeys[i] = k
			i++
		}

		minionIDs := sample(minionKeys, REPLICATION_FACTOR)

		minionAddrs := make([]MinionAddr, 0, REPLICATION_FACTOR-1)




		for j:=0; j < len(minionIDs); j++ {
			if value, ok := minions[minionIDs[j]]; ok {
				minionAddrs = append(minionAddrs, value)
			} else {
				fmt.Println("key not found")
			}
		}

		blockMinions[blockId] = minionIDs
		fileMap[file.Name] = append(fileMap[file.Name], blockId)

		tempBlock := Block{blockId, minionAddrs}

		returnBlocks = append(returnBlocks, tempBlock)

	}

	return returnBlocks, err

}

func main() {


	master := new(Master)
	// Publish the receivers methods
	err := rpc.Register(master)
	if err != nil {
		log.Fatal("Format of service Task isn't correct. ", err)
	}
	// Register a HTTP handler
	rpc.HandleHTTP()
	// Listen to TPC connections on port 1234
	listener, e := net.Listen("tcp", "localhost:8000")
	if e != nil {
		log.Fatal("Listen error: ", e)
	}
	log.Printf("Serving RPC server on port %d", 8000)
	// Start accept incoming HTTP connections
	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("Error serving: ", err)
	}
}
