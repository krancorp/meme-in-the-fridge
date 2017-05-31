package main

import(
	"net"
	"fmt"
	"strings"
	"time"
	"strconv"
	"io/ioutil"
)
const http404 string = "HTTP/1.1 400 Bad Request \r\nContent-Length: 50\r\nContent-Type: text/html\r\n\r\n<html><body><h1>400 Bad Request</h1></body></html>"
const http408 string = "HTTP/1.1 408 Request Time-out \r\nContent-Length: 55\r\nContent-Type: text/html\r\n\r\n<html><body><h1>408 Request Time-out</h1></body></html>"

type htmlRenderer func() []byte

func getHttpRequest(conn net.Conn) (route, method string) {
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
	firstLine := strings.Split(requestLines[0], " ")
	method = firstLine[0]
	route = firstLine[1]
	return method, route
}
func rendBills(tableHeader string)[] byte{
	var html string	
	html += "<ul>"
	for _, b := range bills{
		html += "<li>" + b + "</li>"  
	}
	html += "</ul>"
	return []byte(html)
}
func rendStock(tableHeader string) [] byte{
// Combine the HTML fragments
	fh, _ := ioutil.ReadFile("./stock_head")
	ff, _ := ioutil.ReadFile("./stock_foot")
	bth := []byte(tableHeader)
	buttonLine := "<tr><td></td>"
	for _, i := range productNames{
		buttonLine += "<td> <form action = \"/buy/" + i + "\" method=\"post\"><button type=\"submit\"> 10 " + i + " kaufen! </button></form>"	
	}
	bbl := []byte(buttonLine)
	f := append(fh, bth...)
	f = append(f, bbl...)
	for i := len(lastEntries)-1; i>=0; i--{
		btc := []byte(lastEntries[i])
		f = append(f, btc...)
	}
	f = append(f, ff...)
	return f
}
func startHttpServer(fridgeStock map[string]int, tableHeader string){
	//Web-interface in the making
	fmt.Println("Starting Http-Server")
	ln, err := net.Listen("tcp", ":80")
	
	CheckError(err)
	for {
		conn, err := ln.Accept()
		CheckError(err)
		method, route := getHttpRequest(conn)
		fmt.Println(route, method)
		if(route == "/laquenta" && method == "GET"){
			rendHtml := rendBills(tableHeader)
			httpHeader := "HTTP/1.1 200 OK \r\nContent-Length: "+strconv.Itoa(len(rendHtml))+"\r\nContent-Type: text/html\r\n\r\n"
			conn.Write(append([]byte(httpHeader), rendHtml ...))
			conn.Close()
		}else if(route == "/stock" && method == "GET"){
			rendHtml := rendStock(tableHeader)
			httpHeader := "HTTP/1.1 200 OK \r\nContent-Length: "+strconv.Itoa(len(rendHtml))+"\r\nContent-Type: text/html\r\n\r\n"
			conn.Write(append([]byte(httpHeader), rendHtml ...))
			conn.Close()
		}else if(method == "POST"){
			for _, subRoute := range productNames{
				if(route == "/buy/"+subRoute){
					buy(subRoute)
					rendHtml := "<meta http-equiv=\"refresh\" content=\"0; url=../stock\"/>"
					httpHeader := "HTTP/1.1 200 OK \r\nContent-Length: "+strconv.Itoa(len(rendHtml))+"\r\nContent-Type: text/html\r\n\r\n"
					conn.Write(append([]byte(httpHeader), rendHtml ...))
					conn.Close()
				}
			}
		} else{
			conn.Write([]byte(http404))
			conn.Close()
		}
	

	}
}
