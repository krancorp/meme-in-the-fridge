package main
import(
	"fmt"
	"time"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	//"github.com/eclipse/paho.mqtt.golang/packets"
)
type product struct{
	name string
	price float64
}

func main() {
		
	brokerAddress := "tcp://localhost:1883";
	opts := MQTT.NewClientOptions()
	opts.AddBroker(brokerAddress)
	opts.SetClientID("custom-store")
	

	var testCallback MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	}
	var requestCallback MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	}

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	c.Subscribe("orders", 0, testCallback)
	c.Subscribe("requests", 0, requestCallback)
	c.Subscribe("bargains", 0, testCallback)
	time.Sleep(1* time.Second)
	token := c.Publish("bargains", 0 , false, "Pfungsstaedter! Pfungstaedter zum halben Preis! Kauft und sauft die Scheisse!")
	token.Wait()

	for{
		time.Sleep(1*time.Second)
	}
}

