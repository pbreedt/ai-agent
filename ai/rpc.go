package ai

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// StartRPCServer starts an RPC server on port 1234
// and blocks until the server is shut down
func StartRPCServer(a *Agent) {
	err := rpc.Register(a)
	if err != nil {
		log.Fatal("failed to register:", err)
	}

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening on :1234")

	var wg sync.WaitGroup
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting:", err)
				return
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer conn.Close()
				rpc.ServeConn(conn)
			}()
		}
	}()

	// Handle interrupts
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	fmt.Printf("Received signal: %v\n", sig)

	// Shutdown
	fmt.Println("Shutting down server...")
	listener.Close()

	// Wait for connections to complete with a timeout
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-timeoutCtx.Done():
		fmt.Println("Timeout waiting for connections to close")
	case <-done:
		fmt.Println("All connections closed")
	}

	fmt.Println("Server shutdown complete")
}

func StartRPCServerHTTP(a *Agent) {
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
	return rpc.Dial("tcp", ":1234")
}

func GetRPCClientHTTP() (*rpc.Client, error) {
	return rpc.DialHTTP("tcp", ":1234")
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
