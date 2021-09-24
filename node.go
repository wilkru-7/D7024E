package main

// Importing packages
import (
	"github.com/wilkru-7/D7024E/d7024e"
    "net"
    "fmt"
)
  
// Main function
func main() {
    
    conn,_ := net.Dial("ip:icmp","google.com")
    fmt.Println(conn.LocalAddr())

    if (conn.LocalAddr().String() == "172.21.0.2"){
        me := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), "172.21.0.2:8080") 
        rt := d7024e.NewRoutingTable(me)
        network := d7024e.NewNetwork(&me, rt)
        network.Listen("127.0.0.1", "8080")
    } else if (conn.LocalAddr().String() == "172.21.0.3") {
        known := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), "172.21.0.2:8080") 
        me := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000001"), "172.21.0.3:8080")
        rt := d7024e.NewRoutingTable(me)
        rt.AddContact(known)
        network := d7024e.NewNetwork(&me, rt)
        kademlia := d7024e.NewKademlia(*rt, network)
        go kademlia.LookupContact(&me)
        network.Listen("127.0.0.1", "8080")
    } else {
        known := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), "172.21.0.2:8080") 
        me := d7024e.NewContact(d7024e.NewRandomKademliaID(), d7024e.GetLocalIP()+":8080")
        rt := d7024e.NewRoutingTable(me)
        rt.AddContact(known)
        network := d7024e.NewNetwork(&me, rt)
        kademlia := d7024e.NewKademlia(*rt, network)
        go kademlia.LookupContact(&me)
        network.Listen("127.0.0.1", "8080")
    }
  
}