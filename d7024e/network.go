package d7024e

//https://github.com/holwech/UDP-module/blob/master/UDPmodule.go

import(
	"net"
	"fmt"
	//"time"
	//"net/http"
	//"os"
	"encoding/json"

	//"github.com/labstack/echo/v4"
	//"github.com/labstack/echo/v4/middleware"
	
)
type Network struct {
	contact *Contact
	rt *RoutingTable
}

type Message struct {
	SenderIP string
	ReceiverIP string
	RPC string
}

// NewRoutingTable returns a new instance of a RoutingTable
func NewNetwork(contact *Contact, rt *RoutingTable) *Network {
	network := &Network{}
	network.contact = contact
	network.rt = rt
	return network
}

func Listen(ip string, port string) {
	/*e := echo.New()
	// TODO
	fmt.Println("We are in the listening")
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})
	e.Start(":8080")
	ln, err := net.Listen("tcp", "localhost:8080")
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
	}*/

	fmt.Println("11111111")

	localAddress, err1 := net.ResolveUDPAddr("udp", port)
	if err1 != nil {
		// handle error
	}

	connection, err2 := net.ListenUDP("udp", localAddress)
	if err2 != nil {
		// handle error
	}

	fmt.Println("22222222")

	defer connection.Close()

	for {
		fmt.Println("3333333")
		var message Message
		buffer := make([]byte, 4096)
		length, _, err := connection.ReadFromUDP(buffer)
		if err != nil {
			// handle error
		}
		buffer = buffer[:length]
		err = json.Unmarshal(buffer, &message)
		if message.RPC == "ping" {
			fmt.Println("ping pong")
		}
	}

}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
	
	/* conn, err1 := net.Dial("udp", contact.Address)
	if err1 != nil {
	}
	fmt.Println("YES WE ARE PINGING")
	fmt.Fprintf(conn, "Hi UDP Server, How are you doing?") */
	
	//status, err := bufio.NewReader(conn).ReadString('\n')
	//status, err := bufio.NewReader(conn).ReadString('\n')

    //port := "4000"

    /*timeout := time.Duration(1 * time.Second)
    _, err := net.DialTimeout("udp", contact.Address + "/ping", timeout)
    if err != nil {
        fmt.Printf("%s %s %s\n", contact.Address, "not responding", err.Error())
    } else {
        fmt.Printf("%s %s \n", "responding on adress:", contact.Address)
    }*/

	fmt.Println("send ping message")

	contactAddress, _ := net.ResolveUDPAddr("udp", contact.Address)
	localAddress, _ := net.ResolveUDPAddr("udp", GetLocalIP())
	connection, _ := net.DialUDP("udp", localAddress, contactAddress)
	defer connection.Close()

	for {
		message := &Message{}
		message.SenderIP = network.contact.Address
		message.ReceiverIP = contact.Address
		message.RPC = "ping"
		convMsg, err := json.Marshal(message)
		connection.Write(convMsg)
		if err != nil {
			// handle error
		}
	}


}

func GetLocalIP() string {
	var localIP string
	addr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("GetLocalIP in communication failed")
		return "localhost"
	}
	for _, val := range addr {
		if ip, ok := val.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				localIP = ip.IP.String()
			}
		}
	}
	return localIP
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO

	fmt.Println("Sending find contact message not implemented :( ")
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
