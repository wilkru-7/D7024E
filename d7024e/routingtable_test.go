package d7024e

import (
	"testing"
)

/*func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
}*/

func TestRoutingTable_AddContact(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("0000000000000000000000000000000000000000"), ""))
	contact := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "")
	rt.AddContact(contact)
	if rt.getBucketIndex(contact.ID) != 159 {
		t.Error("Error")
	}
}

func TestRoutingTable_getBucketIndex(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("0000000000000000000000000000000000000000"), ""))
	contact := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	rt.AddContact(contact)
	if rt.getBucketIndex(contact.ID) != 159 {
		t.Error("Error")
	}
}

func TestRoutingTable_FindClosestContacts(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("0000000000000000000000000000000000000000"), ""))
	contact := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "")
	contact2 := NewContact(NewKademliaID("0000000000000000000000000000000000000010"), "")
	contact3 := NewContact(NewKademliaID("0000000000000000000000000000000000000100"), "")
	rt.AddContact(contact)
	rt.AddContact(contact2)
	rt.AddContact(contact3)
	contacts := rt.FindClosestContacts(NewKademliaID("0000000000000000000000000000000000000011"), 1)
	if len(contacts) != 1 && contacts[0] != contact2 {
		t.Error("Error")
	}
	contacts = rt.FindClosestContacts(NewKademliaID("0000000000000000000000000000000000000100"), 2)
	if len(contacts) != 2 && contacts[0] != contact3 && contacts[1] != contact {
		t.Error("Error")
	}
	contacts = rt.FindClosestContacts(NewKademliaID("0000000000000000000000000000000000000100"), 4)
	if len(contacts) != 3 {
		t.Error("Error")
	}
}