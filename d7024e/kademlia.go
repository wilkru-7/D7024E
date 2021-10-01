package d7024e

import(
	"fmt"
	"time"
)
type Kademlia struct {
	rt RoutingTable
	net *Network
	alpha int
}

func NewKademlia(rt RoutingTable, net *Network) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.rt = rt
	kademlia.alpha = 3
	kademlia.net = net
	return kademlia
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

func (kademlia *Kademlia) LookupData(hash string) []string {
	var shortlist ContactCandidates
	target := NewContact(NewKademliaID(hash), "")
	shortlist.contacts = kademlia.rt.FindClosestContacts(target.ID, kademlia.alpha)
	var value, sender string
	var k_triples []Contact
	var visited ContactCandidates
	counter := 0
	for len(shortlist.contacts) > 0 {
		kademlia.net.createRPC("FIND_VALUE", &shortlist.contacts[0], "", []Contact{}, hash, "")
		visited.contacts = append(visited.contacts, shortlist.contacts[0])
		value = <- kademlia.net.findValueChannel
		sender = <- kademlia.net.senderChannel
		if value != "nil" {
			return []string{value, sender}
		}
		k_triples = <- kademlia.net.c
		for _, s := range k_triples{ 
			if(!contains(visited.contacts, s) && !contains(shortlist.contacts, s)) {
				s.CalcDistance(target.ID)
				shortlist.contacts = append(shortlist.contacts, s)
			}
		}
		if len(shortlist.contacts) == 1 {
			shortlist.contacts = []Contact{}
		} else {	
			shortlist.contacts = shortlist.contacts[1:]
		}
		counter += 1
	}
	visited.Sort()
	var result []string
	for i, contact := range visited.contacts {
		if i < 3 {
			result = append(result, contact.String())
		}
	}
	return result
}

//data []bytes
func (kademlia *Kademlia) Store(key *KademliaID, value string) bool {
	target := NewContact(key, "")
	contacts := kademlia.LookupContact(&target)
	var response bool
	successful := false
	for _, contact := range contacts{
		kademlia.net.createRPC("STORE", &contact, "", []Contact{}, key.String(), value)
		response = <- kademlia.net.storeChannel
		if(response){
			fmt.Println("Store completed")
			go kademlia.updateTTL(contact, key.String())
			successful = true
		}else{
			fmt.Println("Store failed")
		}
	}
	return successful
}

func contains(visited []Contact, contact Contact) bool {
	for _, a := range visited {
	   if a.ID.Equals(contact.ID){
		  return true
	   }
	}
	return false
}
func (kademlia *Kademlia) updateTTL(contact Contact, key string){
	for _ = range time.Tick(time.Second * 5) {
		fmt.Println("Pingin inside UPdateTTL: ", contact)
		kademlia.net.createRPC("ping", &contact, "", []Contact{}, "", "")
		fmt.Println("WE PINGINGG!")
		var pong string
		pong = <- kademlia.net.pongChannel 
		if(pong == "pong"){
			fmt.Println("Updataing TTL inside PONG")
			kademlia.net.createRPC("UPDATE_TTL", &contact, "", []Contact{}, key, "")
		}
		/* kademlia.net.updateTTL(&contact, key) */
	}
}