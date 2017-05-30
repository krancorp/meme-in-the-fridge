package main
import(
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"flag"
	"math/rand"
	"strings"
	"net"
	"time"
	"strconv"
)
type product struct{
	name string
	price float64
}
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
func main() {
	var p product
	
	flagName := flag.String("p", "default", "Was soll konstruiert werden, Sire?")
	brokerAddress := flag.String("b", "tcp://localhost:1883", "Wohin gehen wir?")
	flag.Parse()

	p.name = *flagName
	p.price = (rand.Float64() * 5 + 5)
	opts := MQTT.NewClientOptions()
	opts.AddBroker(*brokerAddress)
	opts.SetClientID(p.name + "@" + GetLocalIP())
	mqttClient := MQTT.NewClient(opts)

	var orderCallback MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		tmp := strings.Split(string(msg.Payload()), ";")
		if(len(tmp) < 2){
			return
		}
		if(tmp[0] == GetLocalIP()){
			fmt.Println("Received an order for " + tmp[1] +" " + p.name)
		}
	}
	var requestCallback MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		token := mqttClient.Publish("bargains", 0, false, p.name + ";" + strconv.FormatFloat(p.price, 'f', 2, 64) + ";" + GetLocalIP())
		token.Wait()
		fmt.Println("Published a bargain for " + p.name + "(" + strconv.FormatFloat(p.price, 'f', 2, 64) + " per item)")
	}

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	mqttClient.Subscribe("orders/" + p.name, 0, orderCallback)
	fmt.Println("subscribing to requests/"+p.name)
	mqttClient.Subscribe("requests/" + p.name, 0, requestCallback)
	go updateOffers(mqttClient, p)
	select{}
}
func updateOffers(mqttClient MQTT.Client, p product){
	for{
		time.Sleep(25* time.Second)
		p.price = (rand.Float64() * 5 + 5)
		token := mqttClient.Publish("bargains", 0, false, p.name + ";" + strconv.FormatFloat(p.price, 'f', 2, 64) + ";" + GetLocalIP())
		token.Wait()
		fmt.Println("Published a bargain for " + p.name + "(" + strconv.FormatFloat(p.price, 'f', 2, 64) + " per item)")
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
