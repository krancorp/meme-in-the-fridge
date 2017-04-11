package main
 
import (
	"os"
	"fmt"
	"net"
	"time"
	"strconv"
	"math/rand"
)
 
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
    }
}

 //order of cl-arguments: 1. target ip address, 2. monitored item
func main() {
	ServerAddr,err := net.ResolveUDPAddr("udp",os.Args[1])
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	CheckError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)

	content := "something"
	if(len(os.Args)>2){
		content = os.Args[2]
	}
	seed := rand.NewSource(time.Now().UnixNano())
	randomizer := rand.New(seed)

	stock := randomizer.Intn(42)
	
	defer Conn.Close()
	for {	
		switch(content){
				case "Pfungstaedter":
				default : 
					if(stock > 0){
						if(randomizer.Intn(10)>=7){
							stock--
						}
					}
				}
		msg := content+ ":"+strconv.Itoa(stock)
		buf := []byte(msg)
		_,err := Conn.Write(buf)
		if err != nil {
			fmt.Println(msg, err)
		}
		time.Sleep(time.Second * 1)
	}
}
