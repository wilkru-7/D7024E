package main

// Importing packages
import (
	"github.com/wilkru-7/D7024E/d7024e"
    "net"
    "fmt"
    //"sort"
    //"strings"
    //"time"
)
  
// Main function
func main() {
    
    network := d7024e.NewNetwork()
    //contact := d7024e.NewContact(d7024e.NewRandomKademliaID(), "google.com")
    conn,_ := net.Dial("ip:icmp","google.com")
    fmt.Println(conn.LocalAddr())
    //var me *d7024e.Contact
    //var known *d7024e.Contact
    if (conn.LocalAddr().String() == "172.20.0.2"){
        //fmt.Println("We are in the if")
        me := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), "172.20.0.2:8080") 
        //network.contact = d7024e.NewContact(d7024e.NewRandomKademliaID(), "localhost:80")
        //network.SendPingMessage(&contact)
        rt := d7024e.NewRoutingTable(me)
        rt.AddContact(me)
        //kademlia := d7024e.NewKademlia(*rt, *network)
        //kademlia.LookupContact(&me)
    } else {
        known := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), "172.20.0.2:8080") 
        me := d7024e.NewContact(d7024e.NewRandomKademliaID(), "localhost")
        //network.contact = d7024e.NewContact(d7024e.NewRandomKademliaID(), "localhost:80")
        //network.SendPingMessage(&contact)
        rt := d7024e.NewRoutingTable(me)
        rt.AddContact(known)
        network.SendPingMessage(&known)
        //kademlia := d7024e.NewKademlia(*rt, *network)
        //kademlia.LookupContact(&me)
    }
    d7024e.Listen("127.0.0.1", 8080)
    //network.contact = d7024e.NewContact(d7024e.NewRandomKademliaID(), "localhost:80")
    //network.SendPingMessage(&contact)
    //rt := d7024e.NewRoutingTable(me)
    //rt.AddContact(known)
    //kademlia := d7024e.NewKademlia(*rt, *network)
    //kademlia.LookupContact(&me)
    
    
   


    

 
  
}