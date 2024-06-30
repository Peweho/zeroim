package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"zeroim/common/libnet"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:15678")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server. Enter text to send:")

	protocol := libnet.NewIMProtocol()
	codec := protocol.NewCodec(conn)

	go readServerResponse(codec)

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		msg := libnet.Message{
			Body: []byte(text),
		}
		err = codec.Send(msg)
		if err != nil {
			fmt.Printf("send error: %v\n", err)
		}
	}
}

func readServerResponse(codec libnet.Codec) {
	for {

		msg, err := codec.Receive()
		if err != nil {
			fmt.Println("Error reading from server:", err)
			break
		}

		fmt.Println("Server response: " + string(msg.Body))
	}
}
