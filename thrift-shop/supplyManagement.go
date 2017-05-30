package main
import(
	"strconv"
	"strings"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)
var mqttClient MQTT.Client
func modifyPrefProd(product string, price float64, prefProd string,){
	supplies[product].Lock()
	defer supplies[product].Unlock()
	supplies[product].prefProd = prefProd
	supplies[product].price = price
}
func modifySupply(product string, amount int32){
	supplies[product].Lock()
	defer supplies[product].Unlock()
	supplies[product].stock += amount
}
func initMQTT(){
	opts := MQTT.NewClientOptions()
	opts.AddBroker(brokerIP)
	opts.SetClientID("ThriftShop@" + GetLocalIP())
	mqttClient = MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	mqttClient.Subscribe("bargains", 0, func(client MQTT.Client, msg MQTT.Message){
			tmp := strings.Split(string(msg.Payload()), ";")
			if(len(tmp) < 3){
				return
			}
			price, err := strconv.ParseFloat(tmp[1], 64)
			if(err!=nil){
				fmt.Println(err)
			}
			fmt.Println("Received offer for " + tmp[0])
			if _, ok := supplies[tmp[0]]; ok{
				if (price < supplies[tmp[0]].price){
					modifyPrefProd(tmp[0], price, tmp[2])
					evalSupplies(tmp[0])
				}
			}
		})
}
func evalSupplies(product string){
	if(supplies[product].stock <= 20){
		
		token := mqttClient.Publish("requests/" + product, 0, false, "Need;61")
		token.Wait()
		fmt.Println("Published request for " + product)
	}
	if((supplies[product].stock <= 10 )&& (supplies[product].prefProd != "default")){
		token := mqttClient.Publish("orders/" + product, 0, false, supplies[product].prefProd+";61")
		token.Wait()
		fmt.Println("Ordered" + product + "from " + supplies[product].prefProd)
		modifySupply(product, 61)
	}
}
