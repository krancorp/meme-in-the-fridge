package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"io"
	"io/ioutil"
	"strings"
	"strconv"
	"sort"
	"time"
)

func init(){
	os.Create("./stock_content")
}
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
	s := make([]string, 4)
	//Start Reading..
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString(10) // 0x0A separator = newline
		if err == io.EOF {
			break
		} 
		CheckError(err)
		key := strings.TrimSpace(line)
		s = append(s, key)
		m[key] = 0
	}
	//do some work on the html
	sort.Strings(s)
	f, err := os.OpenFile("./stock_content", os.O_APPEND|os.O_WRONLY, 0600)
	CheckError(err)	
	defer f.Close()
	f.WriteString("<tr> <th>Zeitstempel</th>")
	for i := range s{
		if(len(s[i])>0){
			f.WriteString("<th>"+s[i]+"</th>")
		}	
	}
	f.WriteString("</tr>")
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
	for v := range keys {
		splitStr := strings.Split(keys[v], "%%")
		keys[v] = splitStr[1]
	}
	keys = append(keys, "</tr>")

	f, err := os.OpenFile("./stock_content", os.O_APPEND|os.O_WRONLY, 0600)
	CheckError(err)	
	defer f.Close()
	for v := range keys {
		f.WriteString(keys[v])	
	}
	f.Sync()
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
func startTcpServer(){
	//Web-interface in the making
	fmt.Println("Starting tcp-Server")
	ln, err := net.Listen("tcp", ":80")
	
	CheckError(err)
	for {
		conn, err := ln.Accept()
		fmt.Println("Request Incoming")
		CheckError(err)
		go handleWebRequest(conn)
	}
}

func handleWebRequest(conn net.Conn){
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	// Combine the HTML fragments
	fh, _ := ioutil.ReadFile("./stock_head")
	fb, _ := ioutil.ReadFile("./stock_content")
	ff, _ := ioutil.ReadFile("./stock_foot")
	f := append(fh,fb...)
	f = append(f, ff...)
	conn.Write(f)
	// Close the connection when you're done with it.
	conn.Close()
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
	go startTcpServer()
	//Start and handle UDP server
	port, configPath := ":8080","./config.json"
	
	fmt.Println("Reading Config...")
	fridgeStock := readConfig(configPath)

	fmt.Println("Preparing Server...")
	/* Lets prepare a address at any address at port 	*/   
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

