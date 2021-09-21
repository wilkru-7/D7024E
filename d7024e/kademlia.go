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
// Hej -> 0002102310023013021
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
			kademlia.net.SendPingMessage(&contact)
			var pong string
			pong = <- kademlia.net.pongChannel 
			if(pong == "pong"){
				kademlia.net.SendFindContactMessage(&contact, *target)
				visited = append(visited, contact)
				var k_triples []Contact
				k_triples = <- kademlia.net.c
				for _, s := range k_triples{
					s.CalcDistance(target.ID) 
					/* fmt.Println("closestNode is: ", closestNode)
					fmt.Println("s is: ", s) */
					shortlist.contacts = append(shortlist.contacts, s)
					if(s.Less(&closestNode)){
						closestNode = s
					}
				}
			}
		}
		
	}
	shortlist.Sort()
	fmt.Println("shortlist: ", shortlist.contacts)
	if(shortlist.Len() < 3){
		return shortlist.contacts
	}else{
		return shortlist.contacts[:3]
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	target := NewContact(NewKademliaID(hash), "")
	contacts := kademlia.LookupContact(&target)
	for _, contact := range contacts{
		kademlia.net.SendFindDataMessage(&contact, hash)
	}
}

//data []bytes
func (kademlia *Kademlia) Store(key string, value string) {
	key2 := NewKademliaID(key)
	target := NewContact(key2, "")
	contacts := kademlia.LookupContact(&target)
	for _, contact := range contacts{
		kademlia.net.SendStoreMessage(&contact, key, value)
		response := <- kademlia.net.storeChannel
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
	   if a.key == &hash {
		  return true
	   }
	}
	return false
 }