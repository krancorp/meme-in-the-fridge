package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"io"
	"strings"
	"strconv"
	"sort"
)
 
/* A Simple function to verify error */
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}

func readConfig (path string) (m map[string]int)  {
	//Open File
	file, err := os.Open(path)
	CheckError(err)
	defer file.Close()

	m = make(map[string]int)

	//Start Reading..
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString(10) // 0x0A separator = newline
		if err == io.EOF {
			break
		} 
		CheckError(err)
		key := strings.TrimSpace(line)
		m[key] = 0
	}
	return
	
}

func msgDigest(c chan string, m map[string]int) { 
	for {
		msg := <-c	
		msgStr := strings.Split(msg, ":")
		if (len(msgStr) == 2) {
			//formatting
			key := msgStr[0]
			value, err := strconv.Atoi(msgStr[1])
			if (err != nil) {
				fmt.Println("recieved bad msg 0")
				fmt.Println(msg)		
				continue
			}
			
			if _, ok := m[key]; ok {
				m[key] = value
				printStock(m)
			} else {
				fmt.Println("recieved bad msg")	
				fmt.Println(msgStr)	
				continue	
			}
		} else {
			fmt.Println("recieved bad msg")
			fmt.Println(msgStr)		
			continue
		}
	}
}

func printStock(m map[string]int){
	var keys []string
	for k := range m {
		v := strconv.Itoa(m[k])
		keys = append(keys, k +" : "+ v)
	}
	sort.Strings(keys)
	fmt.Println(keys)
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

func main() {

	port, configPath := ":8080","./config.json"
	
	fmt.Println("Reading Config...")
	fridgeStock := readConfig(configPath)

	fmt.Println("Preparing Server...")
	/* Lets prepare a address at any address at port 	*/   
	ServerAddr,err 	:= net.ResolveUDPAddr("udp", port)
	CheckError(err)
	
	
	fmt.Println("Listening on port ", GetLocalIP(), port)
	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()
	
	messages := make(chan string)
	go msgDigest(messages,fridgeStock)

	buf := make([]byte, 1024)
	
	for {
		n,_,err := ServerConn.ReadFromUDP(buf)
		CheckError(err)
		messages <- string(buf[0:n])
		
	}
}

