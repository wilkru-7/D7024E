package d7024e

import(
	"fmt"
	"time"
)
type Kademlia struct {
	rt RoutingTable
	net *Network
	alpha int
	keys []string
}

/*
 * Creates and returns a new instance of a kademlia.
 * 
 */
func NewKademlia(rt RoutingTable, net *Network) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.rt = rt
	kademlia.alpha = 3
	kademlia.net = net
	kademlia.keys = []string{}
	return kademlia
}

/*
 * LookupContact takes a target contact as an argument and finds
 * the closest nodes to it by traversing through the network.
 * 
 */
func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	var shortlist ContactCandidates
	shortlist.contacts = kademlia.rt.FindClosestContacts(target.ID, kademlia.alpha)
	closestNode := shortlist.contacts[0]
	var visited []Contact
	for _, contact := range shortlist.contacts {
		if(!Contains(visited, contact)){
			kademlia.net.CreateRPC("ping", &contact, "", []Contact{}, "", "")
			var pong string
			pong = <- kademlia.net.pongChannel 
			if(pong == "pong"){
				kademlia.net.CreateRPC("FIND_NODE", &contact, target.ID.String(), []Contact{}, "", "")
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

/*
 * LookupData takes a hash as an argument and tries to find the 
 * corresponding value by looking through the contacts.
 * If no value is found the K closest contacts are returned instead.
 * 
 */
func (kademlia *Kademlia) LookupData(hash string) []string {
	var shortlist ContactCandidates
	target := NewContact(NewKademliaID(hash), "")
	shortlist.contacts = kademlia.rt.FindClosestContacts(target.ID, kademlia.alpha)
	var value, sender string
	var k_triples []Contact
	var visited ContactCandidates
	counter := 0
	for len(shortlist.contacts) > 0 {
		kademlia.net.CreateRPC("FIND_VALUE", &shortlist.contacts[0], "", []Contact{}, hash, "")
		value = <- kademlia.net.findValueChannel
		sender = <- kademlia.net.senderChannel
		if value != "nil" {
			return []string{value, sender}
		}
		k_triples = <- kademlia.net.c
		UpdateShortlist(k_triples, &shortlist, &visited, &target)
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

/*
 * Function used by LookupData to update the shortlist used in
 * the search.
 * 
 */
func UpdateShortlist(k_triples []Contact, shortlist *ContactCandidates, visited *ContactCandidates, target *Contact) {
	visited.contacts = append(visited.contacts, shortlist.contacts[0])
	for _, s := range k_triples{
		if(!Contains(visited.contacts, s) && !Contains(shortlist.contacts, s)) {
			s.CalcDistance(target.ID)
			shortlist.contacts = append(shortlist.contacts, s)
		}
	}
	if len(shortlist.contacts) == 1 {
		shortlist.contacts = []Contact{}
	} else {
		shortlist.contacts = shortlist.contacts[1:]
	}
}

/*
 * Tries to store the wanted key and value pair on the contacts closest
 * to the key hash.
 * 
 */
func (kademlia *Kademlia) Store(key *KademliaID, value string) bool {
	target := NewContact(key, "")
	contacts := kademlia.LookupContact(&target)
	var response bool
	successful := false
	for _, contact := range contacts{
		kademlia.net.CreateRPC("STORE", &contact, "", []Contact{}, key.String(), value)
		response = <- kademlia.net.storeChannel
		if(response){
			fmt.Println("Store completed")
			if (!ContainsString(kademlia.keys, key.String())) {
				kademlia.keys = append(kademlia.keys, key.String())
			}
			go kademlia.UpdateTTL(contact, key.String())
			successful = true
		}else{
			fmt.Println("Store failed")
		}
	}
	return successful
}

/*
 * updateTTL is called from Store asynchronously. It updates the TTL
 * on the keys present in the key array of the kademlia.
 * 
 */
 func (kademlia *Kademlia) UpdateTTL(contact Contact, key string){
	for _ = range time.Tick(time.Second * 5) {
		if ContainsString(kademlia.keys, key) {
			kademlia.net.CreateRPC("ping", &contact, "", []Contact{}, "", "")
			var pong string
			pong = <- kademlia.net.pongChannel 
			if(pong == "pong"){
				kademlia.net.CreateRPC("UPDATE_TTL", &contact, "", []Contact{}, key, "")
			}
		}
	}
	
}

/*
 * Function called when the forget CLI command is executed. 
 * Removes the corresponding key from the kademlia key array
 * which stops its TTL from refreshing.
 * 
 */
func (kademlia *Kademlia) Forget(key string){
	for i, a := range kademlia.keys {
		if a == key{
			kademlia.keys[i] = kademlia.keys[len(kademlia.keys)-1]
    		kademlia.keys = kademlia.keys[:len(kademlia.keys)-1]
			fmt.Println("Removing key: " , a)
		}
	}
}

/*
 * Checks if a contact is present in a contact array.
 * 
 */
func Contains(visited []Contact, contact Contact) bool {
	for _, a := range visited {
	   if a.ID.Equals(contact.ID){
		  return true
	   }
	}
	return false
}

/*
 * Checks if a string is present in a string array.
 * 
 */
func ContainsString(visited []string, key string) bool {
	for _, a := range visited {
	   if a == key{
		  return true
	   }
	}
	return false
}