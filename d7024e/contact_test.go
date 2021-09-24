package d7024e

import (
	"fmt"
	"testing"
)

func TestContact_CalcDistance(t *testing.T) {
	contact := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	contact2 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "")

	contact.CalcDistance(contact2.ID)
	if contact.distance.String() != "0000000000000000000000000000000000000001"{
		t.Error("Error")
	}
}
func TestContact_Swap(t *testing.T) {
	contact := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	contact2 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "")
	contacts := ContactCandidates{}
	contacts.contacts = append(contacts.contacts, contact)
	contacts.contacts = append(contacts.contacts, contact2)
	contacts.Swap(0,1)
	if contacts.contacts[0] != contact2{
		t.Error("Error")
	}
}
func TestContactCandidates_Append(t *testing.T) {
	contact := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	contact2 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "")
	contacts := ContactCandidates{}
	var contacts2 []Contact
	contacts2 = append(contacts2, contact)
	contacts2 = append(contacts2, contact2)
	contacts.Append(contacts2)
	if contacts.Len() != 2{
		t.Error("Error")
	}
}
func TestContactCandidates_GetContacts(t *testing.T) {
	contact := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	contact2 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "")
	contacts := ContactCandidates{}
	contacts.contacts = append(contacts.contacts, contact)
	contacts.contacts = append(contacts.contacts, contact2)
	testContacts := contacts.GetContacts(1)
	if testContacts[0] != contact{
		t.Error("Error")
	}
}
func TestContactCandidates_Sort(t *testing.T) {
	//0000
	contact := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	//0001
	contact2 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "")
	//10000
	contact3 := NewContact(NewKademliaID("0000000000000000000000000000000000000010"), "")
	contact.CalcDistance(contact3.ID)
	contact2.CalcDistance(contact3.ID)
	contacts := ContactCandidates{}
	contacts.contacts = append(contacts.contacts, contact2)
	contacts.contacts = append(contacts.contacts, contact)
	contacts.Sort()
	fmt.Println("contacts.contacts[0]", contacts.contacts[0])
	fmt.Println("contacts.contacts[1]", contacts.contacts[1])
	if contacts.contacts[0] != contact{
		t.Error("Error")
	}
}
func TestContact_String(t *testing.T) {
	contact := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	if contact.String() != "contact(\"0000000000000000000000000000000000000000\", \"\")" {
		t.Error("Error")
	}
}