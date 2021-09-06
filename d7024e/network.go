package d7024e

import(
	"net"
	"fmt"
	
)
type Network struct {
	contact *Contact
}

func Listen(ip string, port int) {
	// TODO
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		//conn, err := ln.Accept()
		_, err := ln.Accept()
		if err != nil {
			// handle error
		}
		//go handleConnection(conn)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
	conn, err := net.Dial("tcp", contact.Address)
	if err != nil {
		
	}
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	//status, err := bufio.NewReader(conn).ReadString('\n')
	//status, err := bufio.NewReader(conn).ReadString('\n')
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
