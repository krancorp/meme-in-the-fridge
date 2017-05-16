package main

import(
	"fmt"
	"os"
	"net"
	"io"
	"bufio"
	"strings"
	"sort"
)

/* A Simple function to verify error */
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
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
func readStores (path string) (m map[string]string) {
	//Open File
	file, err := os.Open(path)
	CheckError(err)
	defer file.Close()
	m = make(map[string]string)
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString(10) // 0x0A separator = newline
		if err == io.EOF {
			break
		} 
		CheckError(err)
		tmp := strings.Split(line, " ")
		m[tmp[0]] = strings.TrimSpace(tmp[1])
	}
	return m
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
