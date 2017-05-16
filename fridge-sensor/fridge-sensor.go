package main
 
import (
	"os"
	"fmt"
	"net"
	"time"
	"strconv"
	"math/rand"
)
var stock int
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
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
func startUdpServer(port string){
	fmt.Println("Preparing UDP-Server...")
	// Lets prepare an address at any address at port  
	ServerAddr,err 	:= net.ResolveUDPAddr("udp", port)
	CheckError(err)
	
	fmt.Println("Listening on ", GetLocalIP(), port)
	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()
	
	buf := make([]byte, 1024)
	
	for {
		n,_,err := ServerConn.ReadFromUDP(buf)
		CheckError(err)
		i, err := strconv.Atoi(string(buf[0:n]))
		CheckError(err)
		stock += i		
	}
}
 //order of cl-arguments: 1. target ip address, 2. monitored item, 3. own upd server port
func main() {
	go startUdpServer(os.Args[3])
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
	
	stock = rand.Intn(42)
	
	defer Conn.Close()
	for {	
		if(stock>0){
			switch(content){
					case "Pfungstaedter":
							growPfungstaedter(&stock)
					case "Grohe" :
							growGrohe(&stock)
					default : 	
						if(rand.Intn(10)>=7){
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
		time.Sleep(time.Second * 5)
	}
}
