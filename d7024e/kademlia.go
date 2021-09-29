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
func (kademlia *Kademlia) LookupData(hash string) []string {
	var shortlist ContactCandidates
	target := NewContact(NewKademliaID(hash), "")
	//shortlist.contacts = kademlia.LookupContact(&target)
	shortlist.contacts = kademlia.rt.FindClosestContacts(target.ID, kademlia.alpha)
	var value, sender string
	var k_triples []Contact
	var visited []Contact
	counter := 0
	//for _, contact := range shortlist.contacts{
	for len(shortlist.contacts) > 0 {
		//fmt.Println("in for loop: ", shortlist.contacts)
		//if(!contains(visited, contact)) {
			kademlia.net.createRPC("FIND_VALUE", &shortlist.contacts[0], "", []Contact{}, hash, "")
			visited = append(visited, shortlist.contacts[0])
			value = <- kademlia.net.findValueChannel
			sender = <- kademlia.net.senderChannel
			if value != "nil" {
				return []string{value, sender}
			}
			k_triples = <- kademlia.net.c
			for _, s := range k_triples{ 
				if(!contains(visited, s) && !contains(shortlist.contacts, s)) {
					s.CalcDistance(target.ID)
					shortlist.contacts = append(shortlist.contacts, s)
					fmt.Println("append contact: ", s)
				}
			}
			if len(shortlist.contacts) == 1 {
				shortlist.contacts = []Contact{}
			} else {
				fmt.Println("removing contact: ", shortlist.contacts[0])
				shortlist.contacts = shortlist.contacts[1:]
			}
			//fmt.Println("in for loop111111: ", shortlist.contacts)
			counter += 1
			fmt.Println("counter is ", counter, "length of shortlist is: ", len(shortlist.contacts))
		//}
	}
	/*shortlist.Sort()
	var result []string
	for i, contact := range shortlist.contacts {
		if i < 3 {
			result = append(result, contact.String())
		}
	}*/
	return []string{"contact test", "test"}
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