package d7024e

import (
"testing"
)

func TestKademliaID_CalcDistance(t *testing.T) {
	id := NewKademliaID("0000000000000000000000000000000000000000")
	xor := id.CalcDistance(NewKademliaID("0000000000000000000000000000000000000001"))
	if xor.String() != "0000000000000000000000000000000000000001" {
		t.Error("Error")
	}
}

func TestKademliaID_Equals(t *testing.T) {
	id := NewKademliaID("0000000000000000000000000000000000000000")
	id2 := NewKademliaID("0000000000000000000000000000000000000001")
	if !id.Equals(id) {
		t.Error("Error")
	}
	if id.Equals(id2) {
		t.Error("Error")
	}
}

func TestKademliaID_Less(t *testing.T) {
	id := NewKademliaID("0000000000000000000000000000000000000000")
	id2 := NewKademliaID("0000000000000000000000000000000000000001")
	if !id.Less(id2) {
		t.Error("Error")
	}
	if id.Less(id) {
		t.Error("Error")
	}
}

func TestKademliaID_NewRandomKademliaID(t *testing.T) {
	id := NewRandomKademliaID()
	id2 := NewRandomKademliaID()
	if id.Equals(id2) {
		t.Error("Error")
	}
}

func TestKademliaID_NewKademliaID(t *testing.T) {
	id := NewKademliaID("00000000000000000000000000000000000")
	id2 := NewKademliaID("0000000000000000000000000000000000000000")
	if !id.Equals(id2) {
		t.Error("Error")
	}
}

func TestNewKademliaID_String(t *testing.T) {
	id := NewKademliaID("00000000000000000000000000000000000")
	string := id.String()
	if string != "0000000000000000000000000000000000000000"{
		t.Error("Error")
	}
}