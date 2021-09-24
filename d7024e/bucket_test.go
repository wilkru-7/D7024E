package d7024e

import (
	"testing"
)

func TestBucket_Len(t *testing.T) {
	bucket := newBucket()
	contact := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	bucket.list.PushFront(contact)
	if bucket.Len() != 1 {
		t.Error("Error")
	}
}

func TestBucket_AddContact(t *testing.T) {
	bucket := newBucket()
	contact := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	bucket.AddContact(contact)
	//Check that it places it in front and not twice
	bucket.AddContact(contact)
	if bucket.list.Len() != 1 {
		t.Error("Error")
	}
}

func TestBucket_GetContactAndCalcDistance(t *testing.T) {
	bucket := newBucket()
	contact := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	bucket.AddContact(contact)
	array := bucket.GetContactAndCalcDistance(NewKademliaID("0000000000000000000000000000000000000001"))
	if array[0].distance.String() != "0000000000000000000000000000000000000001" {
		t.Error("Error")
	}
}