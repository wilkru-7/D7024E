package d7024e

//https://github.com/holwech/UDP-module/blob/e03eccee9bfb5585d2c27c7e153fef273285099c/communication.go#L15

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
	ContactID string
}

// NewRoutingTable returns a new instance of a RoutingTable
func NewNetwork(contact *Contact, rt *RoutingTable) *Network {
	network := &Network{}
	network.contact = contact
	network.rt = rt
	return network
}

func (network *Network) Listen(ip string, port string) {
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


	localAddress, err1 := net.ResolveUDPAddr("udp", GetLocalIP()+":"+port)
	if err1 != nil {
		fmt.Println(err1)
	}
	//fmt.Println("local address in listen : " , localAddress)

	connection, err2 := net.ListenUDP("udp", localAddress)
	if err2 != nil {
		fmt.Println("error connection")
	}
	//fmt.Println("connection in listen : " , connection)

	defer connection.Close()

	for {
		var message Message
		buffer := make([]byte, 4096)
		length, _, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("error in listen")
		}
		buffer = buffer[:length]
		err = json.Unmarshal(buffer, &message)
		if message.RPC == "ping" {
			fmt.Println("received ping from "+ message.SenderIP)
			newContact := NewContact(NewKademliaID(message.ContactID), message.SenderIP)
			network.rt.AddContact(newContact)
			fmt.Println("new contact address: " + newContact.Address)
			network.SendPongMessage(&newContact)
		} else if message.RPC == "pong" {
			fmt.Println("received pong from "+ message.SenderIP)
			newContact := NewContact(NewKademliaID(message.ContactID), message.SenderIP)
			network.rt.AddContact(newContact)
		} else if message.RPC == "FIND_NODE" {
			fmt.Println("received FIND_NODE from "+ message.SenderIP)
			k_contacts := network.rt.FindClosestContacts(NewKademliaID(message.ContactID), 3)
			//These contacts need to be returned to 
			//kademlia file (maybe through channel??)

		}
		else if message.RPC == "FIND_NODE_RETURN" {
			fmt.Println("received FIND_NODE_RETURN from "+ message.SenderIP)
			k_contacts := network.rt.FindClosestContacts(NewKademliaID(message.ContactID), 3)
		}
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	contactAddress, _ := net.ResolveUDPAddr("udp", contact.Address)
	fmt.Println("Sending ping to: " , contactAddress)
	localAddress, _ := net.ResolveUDPAddr("udp", GetLocalIP()+":80")
	fmt.Println("Sending ping from: " , localAddress)
	connection, _ := net.DialUDP("udp", localAddress, contactAddress)
	defer connection.Close()

	convMsg := network.createRPC("ping", contact)
	connection.Write(convMsg)
}


func (network *Network) SendPongMessage(contact *Contact) {
	contactAddress, _ := net.ResolveUDPAddr("udp", contact.Address)
	fmt.Println("Sending pong to: " , contactAddress)
	localAddress, _ := net.ResolveUDPAddr("udp", GetLocalIP()+":80")
	fmt.Println("Sending pong from: " , localAddress)
	connection, _ := net.DialUDP("udp", localAddress, contactAddress)
	defer connection.Close()

	convMsg := network.createRPC("pong", contact)
	connection.Write(convMsg)
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
	contactAddress, _ := net.ResolveUDPAddr("udp", contact.Address)
	fmt.Println("Sending FIND_NODE to: " , contactAddress)
	localAddress, _ := net.ResolveUDPAddr("udp", GetLocalIP()+":80")
	fmt.Println("Sending FIND_NODE from: " , localAddress)
	connection, _ := net.DialUDP("udp", localAddress, contactAddress)
	defer connection.Close()

	convMsg := network.createRPC("FIND_NODE", contact)
	connection.Write(convMsg)
	//fmt.Println("Sending find contact message not implemented :( ")
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func (network *Network) createRPC(rpc string, contact *Contact) ([]byte){
	/* SenderIP string
	ReceiverIP string
	RPC string
	ContactID string */
	message := &Message{}
	message.SenderIP = network.contact.Address
	message.ReceiverIP = contact.Address
	message.RPC = rpc
	message.ContactID = contact.ID.String()
	convMsg, err := json.Marshal(message)
	if err != nil{
		fmt.Println("We got an error in createRPC")
	}
	return convMsg
}