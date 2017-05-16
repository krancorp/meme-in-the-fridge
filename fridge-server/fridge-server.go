package main

var lastEntries []string

func main() {
	configPath := "./config.json"
	fridgeStock, tableHeader := readConfig(configPath)
	block := make(chan bool)
	go ShoppingWrapper(fridgeStock)
	go startUdpServer(fridgeStock)
	go startHttpServer(fridgeStock, tableHeader)
	<- block
}


