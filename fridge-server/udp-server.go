package main

import(
	"fmt"
	"os"
	"strconv"
	"math"
	"time"
	"strings"
	"sort"
	"net"
	"meme-in-the-fridge/thrift-shop/gen-go/store"
	"git.apache.org/thrift.git/lib/go/thrift"
)
var clients[]*store.StoreClient
var bills[] string
func ShoppingWrapper(m map[string]int){
	sm := readStores("./stores")
	for k, v := range sm{
		fmt.Println(k,v)
		addClient(k, v)
	}
	for {
		shopping(m)
	}
		
}
func addClient(host, port string){
	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	trans, err := thrift.NewTSocket(net.JoinHostPort(host, port))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		os.Exit(1)
	}
	client := store.NewStoreClientFactory(trans, protocolFactory)
	if err := trans.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
		os.Exit(1)	
	}
	clients = append(clients, client)
}
func buy(k string){
	var cheapest *store.StoreClient
	price := math.MaxFloat64
	for _, store := range clients{
		tmp, _ := store.GetPrice(k)
		if(tmp < price){
			price = tmp
			cheapest = store
		}
	}			
	if(price != math.MaxFloat64){
		//fmt.Println("Ordered "+k+" for " , 10.0*price,  " $")
		tmp := "Ordered "+k+" for " + strconv.FormatFloat(10.0*price,'f',2,64) + "$"			
		if(len(bills) >= 40){
			//throw away eldest entry
			bills = bills[1:]
		}
		bills = append(bills, tmp)
		cheapest.Order(k, 10)
	}
}
func shopping(m map[string]int){
	time.Sleep(time.Second*10)
	for k, v := range m{
		if(v < 6){
			buy(k)	
		}
	}	
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
	
	fmt.Println("Listening on ", GetLocalIP()+":", port)
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
