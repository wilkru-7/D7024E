package d7024e

import (
"testing"
)

func TestKademlia_Forget(t *testing.T) {
	me := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	rt := NewRoutingTable(me)
	network := NewNetwork(&me, rt)
	kademlia := NewKademlia(*rt, network)
	kademlia.keys = append(kademlia.keys, "test", "test1", "test2")
	kademlia.Forget("test")
	if len(kademlia.keys) != 2 {
		t.Error("Error")
	}
	kademlia.Forget("test3")
	if len(kademlia.keys) != 2 {
		t.Error("Error")
	}
	if kademlia.alpha != 3{
		t.Error("Error")
	}
	if kademlia.rt != *rt {
		t.Error("Error")
	}
	if kademlia.net != network{
		t.Error("Error")
	}
	if NewKademlia(*rt, network) == kademlia{
		t.Error("Error")
	}
	/*kademlia.rt = rt
	kademlia.alpha = 3
	kademlia.net = net
	kademlia.keys = []string{}*/
}

func TestKademlia_Contains(t *testing.T) {
	contact := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	contact2 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "")
	contacts := []Contact{contact}
	if !contains(contacts, contact) {
		t.Error("Error")
	}
	if contains(contacts, contact2) {
		t.Error("Error")
	}
}

func TestKademlia_ContainsString(t *testing.T) {
	strings := []string{"test"}
	if !containsString(strings, "test") {
		t.Error("Error")
	}
	if containsString(strings, "test2") {
		t.Error("Error")
	}
}

func TestKademlia_updateShortlist(t *testing.T) {
	contact := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	contact2 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "")
	contact3 := NewContact(NewKademliaID("0000000000000000000000000000000000000011"), "")
	contact4 := NewContact(NewKademliaID("0000000000000000000000000000000000000010"), "")
	contact5 := NewContact(NewKademliaID("0000000000000000000000000000000000100001"), "")
	target := NewContact(NewKademliaID("0000000000000000000000000000000000000111"), "")
	var shortlist ContactCandidates
	var visited ContactCandidates
	shortlist.contacts = []Contact{contact}
	visited.contacts = []Contact{contact2, contact3}
	k_triples := []Contact{contact2, contact3}
	updateShortlist(k_triples, &shortlist, &visited, &target)
	if len(visited.contacts) != 3 && len(shortlist.contacts) != 0 {
		t.Error("Error")
	}
	shortlist.contacts = []Contact{contact}
	visited.contacts = []Contact{contact2, contact3}
	k_triples = []Contact{contact4, contact5}
	updateShortlist(k_triples, &shortlist, &visited, &target)
	if len(visited.contacts) != 3 && len(shortlist.contacts) != 2 {
		t.Error("Error")
	}

}