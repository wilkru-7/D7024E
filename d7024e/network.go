package d7024e

/* Some of the functionality in the listen and createRPC functions are inspired by this site:
https://github.com/holwech/UDP-module/blob/e03eccee9bfb5585d2c27c7e153fef273285099c/communication.go#L15 */

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

/*
 * Creates and return a new instance of a network.
 * 
 */
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

/*
 * Listening function that takes a port as argument. It sets up a
 * listening connection on the local IP address from the port.
 * Handles the incoming messages and sends the corresponding RPC:s.
 * 
 */
func (network *Network) Listen(port string) {
	localAddress, err := net.ResolveUDPAddr("udp", GetLocalIP()+":"+port)
	if err != nil {
		fmt.Println(err)
	}

	connection, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		fmt.Println(err)
	}

	defer func(connection *net.UDPConn) {
		err := connection.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(connection)

	for {
		var message Message
		buffer := make([]byte, 4096*2)
		
		length, _, err := connection.ReadFromUDP(buffer)

		if err != nil {
			fmt.Println("error in listen")
		}

		buffer = buffer[:length]
		err = json.Unmarshal(buffer, &message)
		network.PickRPC(message)

	}
}

/*
 * Takes a message as argument and takes different actions based
 * on the wanted RPC type. Creates RPC:s, sends values over the channels, etc.
 * 
 */
func (network *Network) PickRPC(message Message){
	RPC := message.RPC
	switch {
	case RPC == "ping":
		fmt.Println("received ping from "+ message.Sender.Address)
		network.rt.AddContact(message.Sender)
		network.CreateRPC("pong", &message.Sender, "", []Contact{}, "", "")
	case RPC == "pong":
		fmt.Println("received pong from "+ message.Sender.Address)
		network.rt.AddContact(message.Sender)
		network.pongChannel <- "pong"
	case RPC == "FIND_NODE":
		fmt.Println("received FIND_NODE from "+ message.Sender.Address)
		k_contacts := network.rt.FindClosestContacts(NewKademliaID(message.TargetID), 3)
		network.CreateRPC("FIND_NODE_RETURN", &message.Sender, "", k_contacts, "", "")
	case RPC == "FIND_NODE_RETURN":
		fmt.Println("received FIND_NODE_RETURN from "+ message.Sender.Address)
		network.c <- message.Contacts
	case RPC == "STORE":
		fmt.Println("received STORE from "+ message.Sender.Address)
		if !network.Contains(message.Key){
			network.AddData(message.Key, message.Value)
			network.CreateRPC("STORE_RETURN", &message.Sender, "", []Contact{}, message.Key, message.Value)
		}else{
			network.CreateRPC("STORE_RETURN_FAIL", &message.Sender, "", []Contact{}, "", "")
		}
	case RPC == "STORE_RETURN":
		network.storeChannel <- true
	case RPC == "STORE_RETURN_FAIL":
		network.storeChannel <- false
	case RPC == "FIND_VALUE":
		fmt.Println("received FIND_VALUE from "+ message.Sender.Address)
		if len(network.data) == 0 {
			network.CreateRPC("FIND_VALUE_RETURN_NIL", &message.Sender, "", []Contact{}, "", "")
			kContacts := network.rt.FindClosestContacts(NewKademliaID(message.Key), 3)
			network.CreateRPC("FIND_NODE_RETURN", &message.Sender, "", kContacts, "", "")
		} else {
			for i, s := range network.data{
				if s.key == message.Key {
					network.data[i].lastAccess = time.Now().Unix()
					fmt.Println("Updating TTL: ", time.Now().Unix())
					network.CreateRPC("FIND_VALUE_RETURN", &message.Sender, "", []Contact{}, "", s.value)
				} else {
					kContacts := network.rt.FindClosestContacts(NewKademliaID(message.Key), 3)
					network.CreateRPC("FIND_VALUE_RETURN_NIL", &message.Sender, "", []Contact{}, "", "")
					network.CreateRPC("FIND_NODE_RETURN", &message.Sender, "", kContacts, "", "")
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

/*
 * Creates a RPC based on the incoming arguments. This process is 
 * incapsuled with a mutex lock to prevent it from trying to
 * create a connection that already exists.
 * 
 */
func (network *Network) CreateRPC(rpc string, receiver *Contact, targetID string, contacts []Contact, key string, value string) {
	network.mu.Lock()
	contactAddress, _ := net.ResolveUDPAddr("udp", receiver.Address)
	fmt.Println("Sending " + rpc + " to: " , contactAddress)
	localAddress, _ := net.ResolveUDPAddr("udp", GetLocalIP()+":80")
	connection, err := net.DialUDP("udp", localAddress, contactAddress)
	
	if err != nil{
		fmt.Println(err)
	}
	
	defer func(connection *net.UDPConn) {
		err := connection.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(connection)

	connection.Write(network.CreateMessage(rpc, receiver, targetID, contacts, key, value))
	network.mu.Unlock()
}

/*
 * Function used by createRPC. Takes the same arguments and returns the
 * message in the form of a []byte.
 * 
 */
func (network *Network) CreateMessage(rpc string, receiver *Contact, targetID string, contacts []Contact, key string, value string) []byte{
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

/*
 * Function that returns the local IP address as a string.
 * This function is copied from the following site:
 * https://github.com/holwech/UDP-module/blob/e03eccee9bfb5585d2c27c7e153fef273285099c/communication.go#L15
 */
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

/*
 * Function called by addData asynchronously. Checks if the 
 * data object is to be removed or not.
 * 
 */
func (network *Network) CheckTTL(data *Data, TTL int){
	for now := range time.Tick(time.Second){
		index := DataGetIndex(network.data, data.key)
		if now.Unix() - network.data[index].lastAccess > int64(TTL){
			if index != -1{
				network.data = Remove(network.data, index)
				fmt.Println("REMOVING OBJECT")
				break
			}
		}
	}
}

/*
 * Takes a key and value as arguments and add these to the data 
 * array of the network. 
 * 
 */
 func (network *Network) AddData(key string, value string) {
	data := Data{}
	data.key = key
	data.value = value
	data.lastAccess = time.Now().Unix()
	network.data = append(network.data, data)
	fmt.Println("Creating data TTL: ", data.lastAccess)
	go network.CheckTTL(&data, TTL)
}

/*
 * Returns the index of a key in a data array.
 * 
 */
func DataGetIndex(data []Data, hash string) int {
	for i, a := range data {
		if a.key == hash {
			return i
		}
	}
	return -1
}

/*
 * Checks if a key is present in the data array of the network.
 * 
 */
func (network *Network) Contains(key string) bool {
	for _, s := range network.data{
		if s.key == key {
			return true
		}
	}
	return false
}

/*
 * Removes the data at index i from the network and returns
 * the new array.
 * 
 */
func Remove(data []Data, i int) []Data{
	if len(data) > i && i > -1 {
		data[i] = data[len(data)-1]
		return data[:len(data)-1]
	}
	return data
}