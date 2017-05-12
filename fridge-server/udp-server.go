package main

import(
	"fmt"
	"strconv"
	"time"
	"strings"
	"sort"
	"net"
)

func msgDigest(c chan string, m map[string]int) { 
	for {
		msg := <-c	
		msgStr := strings.Split(msg, ":")
		if (len(msgStr) == 2) {
			//formatting
			key := msgStr[0]
			value, err := strconv.Atoi(msgStr[1])
			if (err != nil) {
				fmt.Println("received bad msg")
				fmt.Println(msg)		
				continue
			}
			
			if _, ok := m[key]; ok {
				m[key] = value
				genHTMLBody(m)
			} else {
				fmt.Println("received bad msg")	
				fmt.Println(msgStr)	
				continue	
			}
		} else {
			fmt.Println("received bad msg")
			fmt.Println(msgStr)		
			continue
		}
	}
}

func genHTMLBody(m map[string]int){
	var keys []string
	keys = append(keys, " %%<tr> <td>" + time.Now().Format("2006-01-02 15:04:05")+"</td>")
	for k := range m {
		v := strconv.Itoa(m[k])
		keys = append(keys, k + " %% <td>"+ v +"</td>")
	}
	sort.Strings(keys)
	tableRow := ""
	for v := range keys {
		splitStr := strings.Split(keys[v], "%%")
		tableRow += splitStr[1]
	}
	tableRow += "</tr>"
	if(len(lastEntries) >= 40){
		//throw away eldest entry
		lastEntries = lastEntries[1:]
	}
	lastEntries = append(lastEntries, tableRow)	
	return
}

func startUdpServer(fridgeStock map[string]int){
	port := ":8080"
	fmt.Println("Preparing UDP-Server...")
	// Lets prepare an address at any address at port  
	ServerAddr,err 	:= net.ResolveUDPAddr("udp", port)
	CheckError(err)
	
	fmt.Println("Listening on ", GetLocalIP(), port)
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
