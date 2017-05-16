package main

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"fmt"
	"meme-in-the-fridge/thrift-shop/gen-go/shared"
	"os"
	"strconv"
	"math"
	"bufio"
	"io"
	"strings"
	"errors"
	"net"
)

type StoreHandler struct {
	log map[int]*shared.SharedStruct
}
var m map[string]float64
var mp map[string]int64
var ip string
func init(){
	//Open File
	file, err := os.Open("./config")
	if(err!=nil){
		fmt.Println(err)
	}
	defer file.Close()
	m = make(map[string]float64)
	mp = make(map[string]int64)
	//Start Reading..
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString(10) // 0x0A separator = newline
		if err == io.EOF {
			break
		} 
		tmp := strings.Split(line, " ")
		if(tmp[0] == "SENSORIP"){
			ip = strings.TrimSpace(tmp[1])	
		} else {
		m[tmp[0]], err = strconv.ParseFloat(tmp[1], 64)
		if(err!=nil){
				fmt.Println(err)
			}
		mp[tmp[0]], err = strconv.ParseInt(strings.TrimSpace(tmp[2]), 10, 64)
		if(err!=nil){
				fmt.Println(err)
			}
		}
	}
}

func NewStoreHandler() *StoreHandler {
	return &StoreHandler{log: make(map[int]*shared.SharedStruct)}
}

func (p* StoreHandler) GetPrice(product string) (price float64, err error){
	fmt.Println("Got Request for Price of " + product)
	if(product == "Pfungstaedter"){
		fmt.Println("I'm sorry Dave, I'm afraid I can't do that, your taste in beverages is just too horrrible")	
	}
	if val, ok := m[product]; ok{
		return val, nil
	}
	return math.MaxFloat64, errors.New("not in stock")
}

func (p* StoreHandler) Order(product string, amount int32) (err error){
	fmt.Println(amount, " " + product + " was ordered")
	if _, ok := m[product]; ok{
		s:=ip+":" +strconv.Itoa(int(mp[product]))
		fmt.Println(s)
		ServerAddr, err := net.ResolveUDPAddr("udp", s)
		if(err!=nil){
			fmt.Println(err)
		}
		LocalAddr, err := net.ResolveUDPAddr("udp", GetLocalIP()+":0")
		if(err!=nil){
			fmt.Println(err)
		}
		Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
		if(err!=nil){
			fmt.Println(err)
		}
		buf := []byte(strconv.Itoa(int(amount)))
		Conn.Write(buf)
		return  nil
	}
	return errors.New("not in stock")
}

func (p *StoreHandler) GetStruct(key int32) (*shared.SharedStruct, error) {
	fmt.Print("getStruct(", key, ")\n")
	v, _ := p.log[int(key)]
	return v, nil
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
/*
func (p *StoreHandler) Ping() (err error) {
	fmt.Print("ping()\n")
	return nil
}

func (p *StoreHandler) Add(num1 int32, num2 int32) (retval17 int32, err error) {
	fmt.Print("add(", num1, ",", num2, ")\n")
	return num1 + num2, nil
}

*/

