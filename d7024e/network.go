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
	pongChannel chan string
	data []Data
	storeChannel chan string
	findValueChannel chan string
	
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
	key string
	value string
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
	network.pongChannel = make(chan string)
	network.storeChannel = make(chan string)
	network.findValueChannel = make(chan string)
	return network
}

func (network *Network) Listen(ip string, port string) {
	localAddress, err1 := net.ResolveUDPAddr("udp", GetLocalIP()+":"+port)
	if err1 != nil {
		fmt.Println(err1)
	}


	connection, err2 := net.ListenUDP("udp", localAddress)
	if err2 != nil {
		fmt.Println("error connection")
	}
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
			network.pongChannel <- "pong"
		} else if message.RPC == "FIND_NODE" {
			fmt.Println("received FIND_NODE from "+ message.Sender.Address)
			k_contacts := network.rt.FindClosestContacts(message.Target.ID, 3)
			/* fmt.Println("Found these contacts: ", k_contacts) */
			network.SendFindContactMessageReturn(&message.Sender, k_contacts)
		}else if message.RPC == "FIND_NODE_RETURN" {
			fmt.Println("received FIND_NODE_RETURN from "+ message.Sender.Address)
			network.c <- message.Contacts
		}else if message.RPC == "STORE" {
			fmt.Println("received STORE from "+ message.Sender.Address)
			exist := false
			for _, s := range network.data{
				if s.key == message.key {
					exist = true
				}
			}
			if !exist{
				data := Data{}
				data.key = message.key
				data.value = message.value
				network.data = append(network.data, data)
				network.SendStoreReturnMessage(&message.Sender)
			}
		}else if message.RPC =="STORE_RETURN"{
			network.storeChannel <- "Store Completed"
		}else if message.RPC == "FIND_VALUE" {
			fmt.Println("received FIND_VALUE from "+ message.Sender.Address)
			for _, s := range network.data{
				if s.key == message.key {
					network.SendFindDataMessageReturn(&message.Sender, s.value)
				} else {
					k_contacts := network.rt.FindClosestContacts(NewKademliaID(message.key), 3)
					network.SendFindContactMessageReturn(&message.Sender, k_contacts)
				}
			}
		}else if message.RPC == "FIND_VALUE_RETURN" {
			fmt.Println("received FIND_VALUE_RETURN from "+ message.Sender.Address)
			network.findValueChannel <- message.key
		}
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	contactAddress, _ := net.ResolveUDPAddr("udp", contact.Address)
	fmt.Println("Sending ping to: " , contactAddress)
	localAddress, _ := net.ResolveUDPAddr("udp", GetLocalIP()+":80")
	/* fmt.Println("Sending ping from: " , localAddress) */
	connection, _ := net.DialUDP("udp", localAddress, contactAddress)
	defer connection.Close()

	convMsg := network.createRPC("ping", contact)
	connection.Write(convMsg)
}


func (network *Network) SendPongMessage(contact *Contact) {
	contactAddress, _ := net.ResolveUDPAddr("udp", contact.Address)
	fmt.Println("Sending pong to: " , contactAddress)
	localAddress, _ := net.ResolveUDPAddr("udp", GetLocalIP()+":80")
	/* fmt.Println("Sending pong from: " , localAddress) */
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
	/* fmt.Println("Sending FIND_NODE from: " , localAddress) */
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
	fmt.Println("Sending FIND_NODE_RETURN to: " , contactAddress)
	localAddress, _ := net.ResolveUDPAddr("udp", GetLocalIP()+":80")
	/* fmt.Println("Sending FIND_NODE from: " , localAddress) */
	connection, _ := net.DialUDP("udp", localAddress, contactAddress)
	defer connection.Close()

	convMsg := network.createRPCFindNodeReturn("FIND_NODE_RETURN", contact, result)
	connection.Write(convMsg)

	// contacts := <- network.c
	// fmt.Println("Received contacts on channel: ", contacts)
	//fmt.Println("Sending find contact message not implemented :( ")
}

func (network *Network) SendFindDataMessage(contact *Contact, hash string) {
	contactAddress, _ := net.ResolveUDPAddr("udp", contact.Address)
	fmt.Println("Sending FIND_VALUE to: " , contactAddress)
	localAddress, _ := net.ResolveUDPAddr("udp", GetLocalIP()+":80")
	/* fmt.Println("Sending FIND_NODE from: " , localAddress) */
	connection, _ := net.DialUDP("udp", localAddress, contactAddress)
	defer connection.Close()

	convMsg := network.createRPCFindValue("FIND_VALUE", contact, hash)
	connection.Write(convMsg)
}

func (network *Network) SendFindDataMessageReturn(contact *Contact, hash string) {
	contactAddress, _ := net.ResolveUDPAddr("udp", contact.Address)
	fmt.Println("Sending FIND_VALUE to: " , contactAddress)
	localAddress, _ := net.ResolveUDPAddr("udp", GetLocalIP()+":80")
	/* fmt.Println("Sending FIND_NODE from: " , localAddress) */
	connection, _ := net.DialUDP("udp", localAddress, contactAddress)
	defer connection.Close()

	convMsg := network.createRPCFindValue("FIND_VALUE_RETURN", contact, hash)
	connection.Write(convMsg)
}

//data []byte
func (network *Network) SendStoreMessage(contact *Contact, key string, value string) {
	contactAddress, _ := net.ResolveUDPAddr("udp", contact.Address)
	fmt.Println("Sending STORE to: " , contactAddress)
	localAddress, _ := net.ResolveUDPAddr("udp", GetLocalIP()+":80")
	/* fmt.Println("Sending FIND_NODE from: " , localAddress) */
	connection, _ := net.DialUDP("udp", localAddress, contactAddress)
	defer connection.Close()

	convMsg := network.createRPCStore("STORE", contact, key, value)
	connection.Write(convMsg)
}

func (network *Network) SendStoreReturnMessage(contact *Contact) {
	// TODO
	contactAddress, _ := net.ResolveUDPAddr("udp", contact.Address)
	fmt.Println("Sending STORE to: " , contactAddress)
	localAddress, _ := net.ResolveUDPAddr("udp", GetLocalIP()+":80")
	/* fmt.Println("Sending FIND_NODE from: " , localAddress) */
	connection, _ := net.DialUDP("udp", localAddress, contactAddress)
	defer connection.Close()

	convMsg := network.createRPC("STORE_RETURN", contact)
	connection.Write(convMsg)
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

func (network *Network) createRPCStore(rpc string, receiver *Contact, key string, value string) ([]byte){
	/* SenderIP string
	ReceiverIP string
	RPC string
	ContactID string */
	message := &Message{}
	message.Sender = *network.contact
	message.Receiver = *receiver
	message.RPC = rpc
	message.key = key
	message.value = value
	convMsg, err := json.Marshal(message)
	if err != nil{
		fmt.Println("We got an error in createRPC")
	}
	return convMsg
}

func (network *Network) createRPCFindValue(rpc string, receiver *Contact, key string) ([]byte){

	message := &Message{}
	message.Sender = *network.contact
	message.Receiver = *receiver
	message.RPC = rpc
	message.key = key
	convMsg, err := json.Marshal(message)
	if err != nil{
		fmt.Println("We got an error in createRPC")
	}
	return convMsg
}
