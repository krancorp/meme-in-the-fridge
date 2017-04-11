package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"io"
	"strings"
	"strconv"
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
		fmt.Println(m)		
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
			_, ok := m[key]
			fmt.Println(ok)
			fmt.Println(m[key])
			if (ok) {
				m[key] = value
			} else {
				fmt.Println("recieved bad msg 1")	
				fmt.Println(key)	
				continue	
			}
		} else {
			fmt.Println("recieved bad msg 2")
			fmt.Println(msg)		
			continue
		}
	}
}

func main() {

	port, configPath := ":8080","./config.json"
	
	fmt.Println("Reading Config...")
	fridgeStock := readConfig(configPath)
	//fridgeStock["Krombacher"]=11
	fmt.Println(len(fridgeStock))
	fmt.Println("Preparing Server...")
	/* Lets prepare a address at any address at port 	*/   
	ServerAddr,err 	:= net.ResolveUDPAddr("udp", port)
	CheckError(err)
	
	
	fmt.Println("Listening on port %s...", port)
	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()
	
	messages := make(chan string)
	go msgDigest(messages,fridgeStock)

	buf := make([]byte, 1024)
	
	for {
		n,addr,err := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ",string(buf[0:n]), " from ",addr)
		messages <- string(buf[0:n])
		CheckError(err)
	}
}

