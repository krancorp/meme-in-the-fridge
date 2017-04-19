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

func GetLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}

 //order of cl-arguments: 1. target ip address, 2. monitored item
func main() {
	ServerAddr,err := net.ResolveUDPAddr("udp",os.Args[1])
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", GetLocalIP()+":0")
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
		if(stock>0){
			switch(content){
					case "Pfungstaedter":
						if(randomizer.Intn(100)>=98){
							stock--
						}
					default : 	
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
