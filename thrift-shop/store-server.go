package main
import (
	"fmt"
	"meme-in-the-fridge/thrift-shop/gen-go/store"

	"git.apache.org/thrift.git/lib/go/thrift"
)

func runServer(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, addr string, secure bool) error {
	var transport thrift.TServerTransport
	var err error
	
	transport, err = thrift.NewTServerSocket(addr)
	
	if err != nil {
		return err
	}
	fmt.Printf("%T\n", transport)
	handler := NewStoreHandler()
	processor := store.NewStoreProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	fmt.Println("Starting the simple server... on ", addr)
	return server.Serve()
}
