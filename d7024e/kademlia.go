package d7024e

import(
	"fmt"
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
type Triple struct{
	Address  string
	Port int
	ID *KademliaID
}
func (kademlia *Kademlia) LookupContact(target *Contact) Contact{
	// TODO
	shortlist := kademlia.rt.FindClosestContacts(target.ID, kademlia.alpha)
	closestNode := shortlist[0]
	var visited []Contact
	fmt.Println("Closest node: ", closestNode)
	for _, contact := range shortlist{
		if(!contains(visited, contact)){
			kademlia.net.SendPingMessage(&contact)
			kademlia.net.SendFindContactMessage(&contact, *target)
			visited = append(visited, contact)
			var k_triples []Contact
			k_triples = <- kademlia.net.c
			for _, s := range k_triples{
				s.CalcDistance(target.ID)
				shortlist = append(shortlist, s)
				if(s.Less(&closestNode)){
					closestNode = s
				}
			}
			/* fmt.Println("Contacts in LookupContact: ",k_triples) */
			//Check if response
				//then resend find_node to nodes learned about from previous RPC
			// else remove from consideration until they respond
		}
		
	}
	return closestNode
	//fmt.Println("LookupContact is not implemented yet :(")
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func contains(visited []Contact, contact Contact) bool {
	for _, a := range visited {
	   if a.ID == contact.ID {
		  return true
	   }
	}
	return false
 }
