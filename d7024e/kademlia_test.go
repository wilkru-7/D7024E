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
