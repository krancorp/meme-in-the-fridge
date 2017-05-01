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
var lastEntries []string
func init(){
	
}
/* A Simple function to verify error */
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }

}

func readConfig (path string) (m map[string]int, tableHeader string)  {
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
	tableHeader = "<tr> <th>Zeitstempel</th>"
	for i := range s{
		if(len(s[i])>0){
			tableHeader += "<th>"+s[i]+"</th>"
		}	
	}
	tableHeader += "</tr>"
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
func printStock(m map[string]int){
	var keys []string
	for k := range m {
		v := strconv.Itoa(m[k])
		keys = append(keys, k +" : "+ v)
	}
	sort.Strings(keys)
	fmt.Println(keys)
}

const http404 string = "HTTP/1.1 400 Bad Request \r\nContent-Length: 50\r\nContent-Type: text/html\r\n\r\n<html><body><h1>400 Bad Request</h1></body></html>"
const http408 string = "HTTP/1.1 408 Request Time-out \r\nContent-Length: 55\r\nContent-Type: text/html\r\n\r\n<html><body><h1>408 Request Time-out</h1></body></html>"


func handleWebRequest(conn net.Conn, tableHeader string, method string, subUrl string){
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	request := string(buf)
	//check if message is a complete http request, else start timeout
	retry, timeout := 0, 5
	for !strings.Contains(request,"\r\n\r\n"){
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		request = string(buf)
		if(retry >= timeout){
			conn.Write([]byte(http408))
			conn.Close()
			return
		}
		time.Sleep(time.Second * 1)
		retry++					
	}
	requestLines := strings.Split(request,"\r\n")
	//check if the request is valid

	if(!strings.Contains(strings.ToUpper(requestLines[0]), method +" "+ subUrl +" HTTP/1.1")){
		conn.Write([]byte(http404))
		conn.Close()
		return
	}
	// Combine the HTML fragments
	fh, _ := ioutil.ReadFile("./stock_head")
	ff, _ := ioutil.ReadFile("./stock_foot")
	bth := []byte(tableHeader)
	f := append(fh, bth...)
	for i := len(lastEntries)-1; i>=0; i--{
		btc := []byte(lastEntries[i])
		f = append(f, btc...)
	}
	f = append(f, ff...)
	// build the http Header
	header := "HTTP/1.1 200 OK \r\nContent-Length: "+strconv.Itoa(len(f))+"\r\nContent-Type: text/html\r\n\r\n"
	bhh := []byte(header)
	f = append(bhh, f...)
	conn.Write(f)
	// Close the connection when you're done with it.
	conn.Close()
}


func startHttpServer(fridgeStock map[string]int, tableHeader string){
	//Web-interface in the making
	fmt.Println("Starting Http-Server")
	ln, err := net.Listen("tcp", ":80")
	
	CheckError(err)
	for {
		conn, err := ln.Accept()
		CheckError(err)
		go handleWebRequest(conn, tableHeader, "GET", "/STOCK")
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

func main() {
	configPath := "./config.json"
	fridgeStock, tableHeader := readConfig(configPath)
	block := make(chan bool)
	go startUdpServer(fridgeStock)
	go startHttpServer(fridgeStock, tableHeader)
	<- block
}

