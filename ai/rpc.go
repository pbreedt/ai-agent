package ai

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func StartRPCServer(a *Agent) {
	err := rpc.Register(a)
	if err != nil {
		log.Fatal("failed to register:", err)
	}
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	http.Serve(l, nil)
}

/*
// Synchronous call
args := &server.Args{7,8}
var reply int
err = client.Call("Arith.Multiply", args, &reply)

	if err != nil {
		log.Fatal("arith error:", err)
	}

fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)
or

// Asynchronous call
quotient := new(Quotient)
divCall := client.Go("Arith.Divide", args, quotient, nil)
replyCall := <-divCall.Done	// will be equal to divCall
// check errors, print, etc.
*/
func GetRPCClient() (*rpc.Client, error) {
	return rpc.DialHTTP("tcp", ":1234")
}

func (a *Agent) RPCGetAgent(ignored string, agent *Agent) error {
	*agent = *a
	return nil
}

func (a *Agent) RPCRespondToPrompt(prompt string, response *string) error {

	resp, err := a.RespondToPrompt(context.Background(), prompt)
	if err != nil {
		log.Println("AI returned error: ", err.Error())
		*response = err.Error()
		return err
	}

	*response = resp

	return nil
}
