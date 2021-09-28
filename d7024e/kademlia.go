package d7024e

import(
	"fmt"
)
type Kademlia struct {
	rt RoutingTable
	net *Network
	alpha int
	data []Data
}

func NewKademlia(rt RoutingTable, net *Network) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.rt = rt
	kademlia.alpha = 3
	kademlia.net = net
	return kademlia
}
type Triple struct{
	Address  string
	Port int
	ID *KademliaID
}
type Data struct{
	key string
	value string
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	var shortlist ContactCandidates
	shortlist.contacts = kademlia.rt.FindClosestContacts(target.ID, kademlia.alpha)
	closestNode := shortlist.contacts[0]
	var visited []Contact
	for _, contact := range shortlist.contacts {
		if(!contains(visited, contact)){
			kademlia.net.createRPC("ping", &contact, "", []Contact{}, "", "")
			var pong string
			pong = <- kademlia.net.pongChannel 
			if(pong == "pong"){
				kademlia.net.createRPC("FIND_NODE", &contact, target.ID.String(), []Contact{}, "", "")
				visited = append(visited, contact)
				var k_triples []Contact
				k_triples = <- kademlia.net.c
				for _, s := range k_triples{
					s.CalcDistance(target.ID) 
					shortlist.contacts = append(shortlist.contacts, s)
					if(s.Less(&closestNode)){
						closestNode = s
					}
				}
			}
		}
	}
	shortlist.Sort()
	if(shortlist.Len() < 3){
		return shortlist.contacts
	}else{
		return shortlist.contacts[:3]
	}
}

// NOT DONE
func (kademlia *Kademlia) LookupData(hash string) string {
	target := NewContact(NewKademliaID(hash), "")
	contacts := kademlia.LookupContact(&target)
	var value string
	var k_triples []Contact
	var visited []Contact
	for _, contact := range contacts{
		if(!contains(visited, contact)) {
			fmt.Println("11111")
			kademlia.net.createRPC("FIND_VALUE", &contact, "", []Contact{}, hash, "")
			fmt.Println("22222")
			value = <- kademlia.net.findValueChannel
			fmt.Println("33333")
			k_triples = <- kademlia.net.c
			fmt.Println("44444")
			if value != "nil" {
				return value
			}
			fmt.Println("55555")
			visited = append(visited, contact)
			for _, s := range k_triples{ 
				contacts = append(contacts, s)
			}
			fmt.Println("66666")
		}
	}
	return "nil"
}

//data []bytes
func (kademlia *Kademlia) Store(key *KademliaID, value string) {
	target := NewContact(key, "")
	contacts := kademlia.LookupContact(&target)
	var response string
	for _, contact := range contacts{
		kademlia.net.createRPC("STORE", &contact, "", []Contact{}, key.String(), value)
		response = <- kademlia.net.storeChannel
		fmt.Println("Store response is: ", response)
	}
	//kademlia.data.value = data
}

func contains(visited []Contact, contact Contact) bool {
	for _, a := range visited {
	   if a.ID == contact.ID {
		  return true
	   }
	}
	return false
}

func dataContains(data []Data, hash KademliaID) bool {
	for _, a := range data {
	   if a.key == hash.String() {
		  return true
	   }
	}
	return false
}
