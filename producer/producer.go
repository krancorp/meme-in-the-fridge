package main
import(
	"fmt"
	"time"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	//"github.com/eclipse/paho.mqtt.golang/packets"
)

func main() {
	brokerAddress := "tcp://localhost:1883";
	opts := MQTT.NewClientOptions()
	opts.AddBroker(brokerAddress)
	opts.SetClientID("custom-store")
	

	var testCallback MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	}

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	c.Subscribe("orders", 0, testCallback)
	c.Subscribe("requests", 0, testCallback)
/*/
	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		token.Wait()
	}

	for i := 1; i < 5; i++ {
		time.Sleep(1 * time.Second)
	}

	c.Disconnect(250)
*/
	for{
		time.Sleep(1*time.Second)
	}
}

