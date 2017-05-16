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
)

type StoreHandler struct {
	log map[int]*shared.SharedStruct
}

func NewStoreHandler() *StoreHandler {
	return &StoreHandler{log: make(map[int]*shared.SharedStruct)}
}

func (p* StoreHandler) GetPrice(product string) (price int32, err error){
	fmt.Println("getPrice called")
	return 1, nil
}

func (p* StoreHandler) Order(product string, amount int32) (err error){
	fmt.Println("order called")
	return nil
}

func (p *StoreHandler) GetStruct(key int32) (*shared.SharedStruct, error) {
	fmt.Print("getStruct(", key, ")\n")
	v, _ := p.log[int(key)]
	return v, nil
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

