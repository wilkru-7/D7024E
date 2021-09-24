package d7024e

//https://github.com/holwech/UDP-module/blob/e03eccee9bfb5585d2c27c7e153fef273285099c/communication.go#L15

import (
	"encoding/json"
	"fmt"
	"net"
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
	RPC string
	TargetID string
	Contacts []Contact
	key string
	value string
}

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
		fmt.Println(err2)
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

		RPC := message.RPC
		switch {
		case RPC == "ping":
			fmt.Println("received ping from "+ message.Sender.Address)
			network.rt.AddContact(message.Sender)
			network.createRPC("pong", &message.Sender, "", []Contact{}, "", "")
		case RPC == "pong":
			fmt.Println("received pong from "+ message.Sender.Address)
			network.rt.AddContact(message.Sender)
			network.pongChannel <- "pong"
		case RPC == "FIND_NODE":
			fmt.Println("received FIND_NODE from "+ message.Sender.Address)
			k_contacts := network.rt.FindClosestContacts(NewKademliaID(message.TargetID), 3)
			network.createRPC("FIND_NODE_RETURN", &message.Sender, "", k_contacts, "", "")
		case RPC == "FIND_NODE_RETURN":
			fmt.Println("received FIND_NODE_RETURN from "+ message.Sender.Address)
			network.c <- message.Contacts
		case RPC == "STORE":
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
				network.createRPC("STORE", &message.Sender, "", []Contact{}, data.key, data.value)
			}
		case RPC == "STORE_RETURN":
			network.storeChannel <- "Store Completed"
		case RPC == "FIND_VALUE":
			fmt.Println("received FIND_VALUE from "+ message.Sender.Address)
			for _, s := range network.data{
				if s.key == message.key {
					network.createRPC("FIND_VALUE_RETURN", &message.Sender, "", []Contact{}, "", s.value)
				} else {
					k_contacts := network.rt.FindClosestContacts(NewKademliaID(message.key), 3)
					network.createRPC("FIND_NODE_RETURN", &message.Sender, "", k_contacts, "", "")
					network.findValueChannel <- "nil"
				}
			}
		case RPC == "FIND_VALUE_RETURN":
			fmt.Println("received FIND_VALUE_RETURN from "+ message.Sender.Address)
			network.findValueChannel <- message.value
		default:
			fmt.Println("Invalid RPC")
		}
	}
}

func (network *Network) createRPC(rpc string, receiver *Contact, targetID string, contacts []Contact, key string, value string) {
	contactAddress, _ := net.ResolveUDPAddr("udp", receiver.Address)
	fmt.Println("Sending " + rpc + " to: " , contactAddress)
	localAddress, _ := net.ResolveUDPAddr("udp", GetLocalIP()+":80")
	connection, _ := net.DialUDP("udp", localAddress, contactAddress)
	defer connection.Close()

	message := &Message{}
	message.Sender = *network.contact
	message.Receiver = *receiver
	message.RPC = rpc
	message.TargetID = targetID
	message.Contacts = contacts
	message.key = key
	message.value = value

	convMsg, err := json.Marshal(message)

	if err != nil{
		fmt.Println("We got an error in createRPC")
	}

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