package d7024e

import(
	"net"
	"fmt"
	"time"
	
)
type Network struct {
	//contact *Contact
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
	/*
	conn, err := net.Dial("tcp", contact.Address)
	if err != nil {
	}
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	 */
	//status, err := bufio.NewReader(conn).ReadString('\n')
	//status, err := bufio.NewReader(conn).ReadString('\n')

    port := "80"
    timeout := time.Duration(1 * time.Second)
    _, err := net.DialTimeout("tcp", contact.Address+":"+port, timeout)
    if err != nil {
        fmt.Printf("%s %s %s\n", contact.Address, "not responding", err.Error())
    } else {
        fmt.Printf("%s %s %s\n", contact.Address, "responding on port:", port)
    }
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
