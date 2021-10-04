package d7024e

//https://github.com/holwech/UDP-module/blob/e03eccee9bfb5585d2c27c7e153fef273285099c/communication.go#L15

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)
const TTL = 20

type Network struct {
	contact *Contact
	rt *RoutingTable
	c chan []Contact
	pongChannel chan string
	data []Data
	storeChannel chan bool
	findValueChannel chan string
	senderChannel chan string
	mu sync.Mutex
}

type Message struct {
	Sender Contact
	Receiver Contact
	RPC string
	TargetID string
	Contacts []Contact
	Key string
	Value string
}

type Data struct{
	key string
	value string
	lastAccess int64
}

func dataGetIndex(data []Data, hash string) int {
	for i, a := range data {
	   if a.key == hash {
		  return i
	   }
	}
	return -1
}

func (network *Network) contains(key string) bool {
	for _, s := range network.data{
		if s.key == key {
			return true
		}
	}
	return false
}

func (network *Network) addData(key string, value string) {
	data := Data{}
	data.key = key
	data.value = value
	data.lastAccess = time.Now().Unix()
	network.data = append(network.data, data)
	fmt.Println("Creating data TTL: ", data.lastAccess)
	go network.checkTTL(&data, TTL)
}

func remove(data []Data, i int) []Data{
	if len(data) > i && i > -1 {
		data[i] = data[len(data)-1]
		return data[:len(data)-1]
	}
	return data
}

func NewNetwork(contact *Contact, rt *RoutingTable) *Network {
	network := &Network{}
	network.contact = contact
	network.rt = rt
	network.c = make(chan []Contact)
	network.pongChannel = make(chan string)
	network.storeChannel = make(chan bool)
	network.findValueChannel = make(chan string)
	network.senderChannel = make(chan string)
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
		buffer := make([]byte, 4096*2)
		
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
			if !network.contains(message.Key){
				network.addData(message.Key, message.Value)

				network.createRPC("STORE_RETURN", &message.Sender, "", []Contact{}, message.Key, message.Value)
			}else{
				network.createRPC("STORE_RETURN_FAIL", &message.Sender, "", []Contact{}, "", "")
			}
		case RPC == "STORE_RETURN":
			network.storeChannel <- true
		case RPC == "STORE_RETURN_FAIL":
			network.storeChannel <- false
		case RPC == "FIND_VALUE":
			fmt.Println("received FIND_VALUE from "+ message.Sender.Address)
			if len(network.data) == 0 {
				network.createRPC("FIND_VALUE_RETURN_NIL", &message.Sender, "", []Contact{}, "", "")
				k_contacts := network.rt.FindClosestContacts(NewKademliaID(message.Key), 3)
				network.createRPC("FIND_NODE_RETURN", &message.Sender, "", k_contacts, "", "")
			} else {
				for i, s := range network.data{
					if s.key == message.Key {
						network.data[i].lastAccess = time.Now().Unix()
						fmt.Println("Updating TTL: ", time.Now().Unix())
						network.createRPC("FIND_VALUE_RETURN", &message.Sender, "", []Contact{}, "", s.value)
					} else {
						k_contacts := network.rt.FindClosestContacts(NewKademliaID(message.Key), 3)
						network.createRPC("FIND_VALUE_RETURN_NIL", &message.Sender, "", []Contact{}, "", "")
						network.createRPC("FIND_NODE_RETURN", &message.Sender, "", k_contacts, "", "")
					}
				}
			}
		case RPC == "FIND_VALUE_RETURN":
			fmt.Println("received FIND_VALUE_RETURN from "+ message.Sender.Address)
			network.findValueChannel <- message.Value 
			network.senderChannel <- message.Sender.String()
		case RPC == "FIND_VALUE_RETURN_NIL":
			fmt.Println("received FIND_VALUE_RETURN_NIL from "+ message.Sender.Address)
			network.findValueChannel <- "nil"
			network.senderChannel <- "nil"
		case RPC == "UPDATE_TTL":
			fmt.Println("Recieved UPDATE_TTL from: ", message.Sender.Address)
			for i, s := range network.data{
				if s.key == message.Key{
					network.data[i].lastAccess = time.Now().Unix()
					fmt.Println("Updating TTL: ", time.Now().Unix())
				}
			}

		default:
			fmt.Println("Invalid RPC")
		}
	}
}

func (network *Network) createRPC(rpc string, receiver *Contact, targetID string, contacts []Contact, key string, value string) {
	network.mu.Lock()
	contactAddress, _ := net.ResolveUDPAddr("udp", receiver.Address)
	fmt.Println("Sending " + rpc + " to: " , contactAddress)
	localAddress, _ := net.ResolveUDPAddr("udp", GetLocalIP()+":80")
	connection, err := net.DialUDP("udp", localAddress, contactAddress)
	
	if err != nil{
		fmt.Println(err)
	}
	
	defer connection.Close()

	connection.Write(network.createMessage(rpc, receiver, targetID, contacts, key, value))
	network.mu.Unlock()
}

func (network *Network) createMessage(rpc string, receiver *Contact, targetID string, contacts []Contact, key string, value string) []byte{
	message := &Message{}
	message.Key = key
	message.Value = value
	message.Sender = *network.contact
	message.Receiver = *receiver
	message.RPC = rpc
	message.TargetID = targetID
	message.Contacts = contacts

	convMsg, _ := json.Marshal(message)

	return convMsg
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

func (network *Network) checkTTL(data *Data, TTL int){
	for now := range time.Tick(time.Second){
		index := dataGetIndex(network.data, data.key)
		if now.Unix() - network.data[index].lastAccess > int64(TTL){
			if index != -1{
				network.data = remove(network.data, index)
				fmt.Println("REMOVING OBJECT")
				break
			}
		}
	}
}