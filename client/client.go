
package main

import (
	"fmt"
	"log"
	"net/rpc"
	"reflect"
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




/*def get(master, file):
file_table = master.read(file)
if not file_table:
logging.info("file not found")
return

for block in file_table:
for host, port in block['block_addr']:
try:
con = rpyc.connect(host, port=port).root
data = con.get(block['block_id'])
if data:
sys.stdout.write(data)
break
except Exception as e:
continue
else:
logging.error("No blocks found. Possibly a corrupt file")
*/

func get(client *rpc.Client, fileName string) ([]Block, error){
	var fileBlocks []Block
	xType := reflect.TypeOf(client)
	fmt.Println(xType)
	err := client.Call("Master.Read", fileName, &fileBlocks)
	if err != nil{
		fmt.Println("File not found")
		return  nil, err
	}
	return fileBlocks, err
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

	xType := reflect.TypeOf(client)
	fmt.Println(xType)

	//fmt.Println(reply)
	fmt.Println("_________________________________\n")
	var reply2 []Block
	//err = client.Call("Master.Read", "/abc.txt", &reply2)
	reply2, err = get(client, "/abc.txt")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	fmt.Println(reply2)
}