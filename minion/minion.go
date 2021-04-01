package minion

import (
	"log"
	"net/rpc"
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
	BlockId int
	Data []byte
	Minions [] AddrMinion
}

var DATA_DIR string =  "/tmp/files/"

func (m * Minion) PUT (message MessageToMinion, reply * int) error {

}

func main() {
	minion := new(Minion)
	// Publish the receivers methods
	err := rpc.Register(minion)
	if err != nil {
		log.Fatal("Format of service Task isn't correct. ", err)
	}
}