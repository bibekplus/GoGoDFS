package main

import (
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"net/rpc"
)


const BLOCK_SIZE = 100
const REPLICATION_FACTOR = 2

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


type Master int


var fileMap = make(map[string][]string)//Filename to blocks

var blockMinions = make(map[string][]string)

//var minion1= MinionAddr{"120.0.0.1", 8000}
//var minion2= MinionAddr{"120.0.0.1", 9000}


var minions = map[string]MinionAddr{
	"1" : MinionAddr{"120.0.0.1", 9000},
	"2" : MinionAddr{"120.0.0.1", 9001},
	"3" : MinionAddr{"120.0.0.1", 9002},

}





func (m * Master) Write (File FileBlock, Blocks * []Block) error{

	_, ok := fileMap[File.Name]
	if ok {
		delete(fileMap, File.Name)
	}

	numOfBlocks := int(math.Ceil(float64(File.Size / BLOCK_SIZE)))


	tempBlocks, err := allocateBlocks(File, numOfBlocks)


	*Blocks = tempBlocks
	fmt.Println("Blocks:", Blocks)


	if err != nil{
		fmt.Println(err)
	}
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
		//minionIds := reflect.ValueOf(minions).Interface().(string)


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
