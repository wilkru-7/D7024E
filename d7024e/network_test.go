package d7024e

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNetwork_dataGetIndex(t *testing.T) {
	data := Data{}
	data.key = "key"
	data.value = "value"
	data.lastAccess = time.Now().Unix()
	if dataGetIndex([]Data{}, "key") != -1 {
		t.Error("Error")
	}
	array := []Data{data}
	if dataGetIndex(array, "key") != 0 {
		t.Error("Error")
	}
	if dataGetIndex(array, "key2") != -1 {
		t.Error("Error")
	}
}

func TestNetwork_remove(t *testing.T) {
	data := Data{}
	data.key = "key"
	data.value = "value"
	data.lastAccess = time.Now().Unix()
	array := []Data{data}
	if len(remove(array, 3)) != 1{
		t.Error("Error")
	}
	if len(remove(array, -1)) != 1{
		t.Error("Error")
	}
	if len(remove(array, 0)) != 0 {
		t.Error("Error")
	}
}

func TestNetwork_checkTTL(t *testing.T) {
	me := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	rt := NewRoutingTable(me)
	network := NewNetwork(&me, rt)
	data := Data{}
	data.key = "key"
	data.value = "value"
	data.lastAccess = time.Now().Unix()
	network.data = []Data{data}
	go network.checkTTL(&data, 1)
	if len(network.data) != 1{
		t.Error("Error")
	}
	time.Sleep(3*time.Second)
	if len(network.data) != 0{
		t.Error("Error data is: ", network.data)
	}
}

func TestNetwork_createMessage(t *testing.T) {
	me := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	rt := NewRoutingTable(me)
	network := NewNetwork(&me, rt)
	message := network.createMessage("rpc", &me, "targetID", []Contact{}, "key", "value")
	var m Message
	json.Unmarshal(message, &m)

	if m.RPC != "rpc"{
		t.Error("Error rpc")
	}
}

func TestNetwork_contains(t *testing.T) {
	me := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	rt := NewRoutingTable(me)
	network := NewNetwork(&me, rt)
	data := Data{}
	data.key = "key"
	data.value = "value"
	data.lastAccess = time.Now().Unix()
	data1 := Data{}
	data1.key = "key1"
	data1.value = "value"
	data1.lastAccess = time.Now().Unix()
	network.data = []Data{data1, data}
	if !network.contains("key") {
		t.Error("Error no key match")
	}

	if network.contains("test") {
		t.Error("Error wrong key match")
	}
}

func TestNetwork_addData(t *testing.T) {
	me := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	rt := NewRoutingTable(me)
	network := NewNetwork(&me, rt)
	network.addData("key", "value")
	if (len(network.data)) != 1 {
		t.Error("Error")
	}
}

func TestNetwork_ping(t *testing.T) {
	me := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")
	rt := NewRoutingTable(me)
	network := NewNetwork(&me, rt)

	message := network.createMessage("ping", &me, "", []Contact{}, "", "")
	var m Message
	json.Unmarshal(message, &m)

	network.pickRPC(m)
	if (len(network.rt.FindClosestContacts(me.ID, 3))) != 1 {
		t.Error("Error")
	}
}