package d7024e

//https://github.com/holwech/UDP-module/blob/e03eccee9bfb5585d2c27c7e153fef273285099c/communication.go#L15

import(
	"net"
	"fmt"
	//"time"
	//"net/http"
	//"os"
	"encoding/json"
	//"strings"

	//"github.com/labstack/echo/v4"
	//"github.com/labstack/echo/v4/middleware"
	
)
type Network struct {
	contact *Contact
	rt *RoutingTable
	c chan []Contact
}

type Message struct {
	Sender Contact
	Receiver Contact
	/* SenderIP string
	ReceiverIP string */
	RPC string
	/* ContactID string
	TargetID string */
	Target Contact
	Contacts []Contact
}
type triple struct {
	IP string
	ID string
}
// NewRoutingTable returns a new instance of a RoutingTable
func NewNetwork(contact *Contact, rt *RoutingTable) *Network {
	network := &Network{}
	network.contact = contact
	network.rt = rt
	network.c = make(chan []Contact)
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
			fmt.Println("received ping from "+ message.Sender.Address)
			network.rt.AddContact(message.Sender)
			network.SendPongMessage(&message.Sender)
		} else if message.RPC == "pong" {
			fmt.Println("received pong from "+ message.Sender.Address)
			network.rt.AddContact(message.Sender)
		} else if message.RPC == "FIND_NODE" {
			fmt.Println("received FIND_NODE from "+ message.Sender.Address)
			k_contacts := network.rt.FindClosestContacts(message.Target.ID, 3)
			fmt.Println("Found these contacts: ", k_contacts)
			network.SendFindContactMessageReturn(&message.Sender, k_contacts)
		}else if message.RPC == "FIND_NODE_RETURN" {
			fmt.Println("received FIND_NODE_RETURN from "+ message.Sender.Address)
			network.c <- message.Contacts
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

func (network *Network) SendFindContactMessage(contact *Contact, target Contact) {
	// TODO
	contactAddress, _ := net.ResolveUDPAddr("udp", contact.Address)
	fmt.Println("Sending FIND_NODE to: " , contactAddress)
	localAddress, _ := net.ResolveUDPAddr("udp", GetLocalIP()+":80")
	fmt.Println("Sending FIND_NODE from: " , localAddress)
	connection, _ := net.DialUDP("udp", localAddress, contactAddress)
	defer connection.Close()

	convMsg := network.createRPCFindNode("FIND_NODE", contact, target)
	connection.Write(convMsg)

	/* contacts := <- network.c
	fmt.Println("Received contacts on channel: ", contacts) */
	//fmt.Println("Sending find contact message not implemented :( ")
}

func (network *Network) SendFindContactMessageReturn(contact *Contact, result []Contact) {
	// TODO
	contactAddress, _ := net.ResolveUDPAddr("udp", contact.Address)
	fmt.Println("Sending FIND_NODE to: " , contactAddress)
	localAddress, _ := net.ResolveUDPAddr("udp", GetLocalIP()+":80")
	fmt.Println("Sending FIND_NODE from: " , localAddress)
	connection, _ := net.DialUDP("udp", localAddress, contactAddress)
	defer connection.Close()

	convMsg := network.createRPCFindNodeReturn("FIND_NODE_RETURN", contact, result)
	connection.Write(convMsg)

	// contacts := <- network.c
	// fmt.Println("Received contacts on channel: ", contacts)
	//fmt.Println("Sending find contact message not implemented :( ")
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func (network *Network) createRPC(rpc string, receiver *Contact) ([]byte){
	/* SenderIP string
	ReceiverIP string
	RPC string
	ContactID string */
	message := &Message{}
	message.Sender = *network.contact
	message.Receiver = *receiver
	message.RPC = rpc
	convMsg, err := json.Marshal(message)
	if err != nil{
		fmt.Println("We got an error in createRPC")
	}
	return convMsg
}
func (network *Network) createRPCFindNode(rpc string, receiver *Contact, target Contact) ([]byte){
	/* SenderIP string
	ReceiverIP string
	RPC string
	ContactID string */
	message := &Message{}
	message.Sender = *network.contact
	message.Receiver = *receiver
	message.RPC = rpc
	message.Target = target
	convMsg, err := json.Marshal(message)
	if err != nil{
		fmt.Println("We got an error in createRPC")
	}
	return convMsg
}
func (network *Network) createRPCFindNodeReturn(rpc string, receiver *Contact, contacts []Contact) ([]byte){
	/* SenderIP string
	ReceiverIP string
	RPC string
	ContactID string */
	message := &Message{}
	message.Sender = *network.contact
	message.Receiver = *receiver
	message.RPC = rpc
	message.Contacts = contacts
	convMsg, err := json.Marshal(message)
	if err != nil{
		fmt.Println("We got an error in createRPC")
	}
	return convMsg
}