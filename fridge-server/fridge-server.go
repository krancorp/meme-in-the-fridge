package main

var lastEntries []string

func main() {
	configPath := "./config.json"
	fridgeStock, tableHeader := readConfig(configPath)
	
	go ShoppingWrapper(fridgeStock)
	go startUdpServer(fridgeStock)
	startHttpServer(fridgeStock, tableHeader)
}


