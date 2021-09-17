package d7024e

import(
	"fmt"
)
type Kademlia struct {
	rt RoutingTable
	net Network
	alpha int
}

func NewKademlia(rt RoutingTable, net Network) *Kademlia {
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
func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
	closestContacts := kademlia.rt.FindClosestContacts(target.ID, kademlia.alpha)
	if len(closestContacts) == 0 {
		return
	}
	for _, contact := range closestContacts {
		kademlia.net.SendFindContactMessage(&contact, *target.ID)
		//Check if response
			//then resend find_node to nodes learned about from previous RPC
		// else remove from consideration until they respond
	}
	fmt.Println("LookupContact is not implemented yet :(")
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
