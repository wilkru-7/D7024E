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
    
    
    //contact := d7024e.NewContact(d7024e.NewRandomKademliaID(), "google.com")
    conn,_ := net.Dial("ip:icmp","google.com")
    fmt.Println(conn.LocalAddr())
    //var me *d7024e.Contact
    //var known *d7024e.Contact
    if (conn.LocalAddr().String() == "172.19.0.2"){
        //fmt.Println("We are in the if")
        me := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), "172.21.0.2:8080") 
        //network.contact = d7024e.NewContact(d7024e.NewRandomKademliaID(), "localhost:80")
        //network.SendPingMessage(&contact)
        rt := d7024e.NewRoutingTable(me)
        /* rt.AddContact(me) */
        network := d7024e.NewNetwork(&me, rt)
        //go network.SendPingMessage(&me)
        network.Listen("127.0.0.1", "8080")
        //kademlia := d7024e.NewKademlia(*rt, *network)
        //kademlia.LookupContact(&me)
        //time.Sleep(1 * time.Second)
    } else if (conn.LocalAddr().String() == "172.19.0.3") {
        known := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), "172.21.0.2:8080") 
        me := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000001"), "172.21.0.3:8080")
        /* third := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000003"), "172.19.0.4:8080") */
        //network.contact = d7024e.NewContact(d7024e.NewRandomKademliaID(), "localhost:80")
        //network.SendPingMessage(&contact)
        rt := d7024e.NewRoutingTable(me)
        /* rt.AddContact(me) */
        rt.AddContact(known)
        network := d7024e.NewNetwork(&me, rt)
        //go network.SendPingMessage(&known)
        //go network.SendFindContactMessage(&known)
        kademlia := d7024e.NewKademlia(*rt, network)
        go kademlia.LookupContact(&me)
        /* go kademlia.LookupContact(&third) */
        network.Listen("127.0.0.1", "8080")
        //time.Sleep(1 * time.Second)
    } else {
        known := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), "172.21.0.2:8080") 
        secondKnown := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000001"), "172.21.0.3:8080")
        /* me := d7024e.NewContact(d7024e.NewRandomKademliaID(), d7024e.GetLocalIP()+":8080") */
        me := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000003"), d7024e.GetLocalIP()+":8080")
        //network.contact = d7024e.NewContact(d7024e.NewRandomKademliaID(), "localhost:80")
        //network.SendPingMessage(&contact)
        rt := d7024e.NewRoutingTable(me)
        /* rt.AddContact(me) */
        rt.AddContact(known)
/*         rt.AddContact(secondKnown) */
        network := d7024e.NewNetwork(&me, rt)
        //go network.SendPingMessage(&known)
        /* go network.SendFindContactMessage(&me, *d7024e.NewKademliaID("0000000000000000000000000000000000000009")) */
        kademlia := d7024e.NewKademlia(*rt, network)
        /* go kademlia.LookupContact(&me) */
        go kademlia.LookupContact(&secondKnown)
        /* contacts := kademlia.LookupContact(&secondKnown)
        fmt.Println("contacts: ",contacts) */
       /*  for _, bucket := range rt.buckets {
            fmt.Println("bucket length: ",bucket.Len())
        } */
       /*  fmt.Println("Routingtable looks like: ", rt) */
        network.Listen("127.0.0.1", "8080")
        //time.Sleep(1 * time.Second)
    }
    //d7024e.Listen("127.0.0.1", "8080")
    //network.contact = d7024e.NewContact(d7024e.NewRandomKademliaID(), "localhost:80")
    //network.SendPingMessage(&contact)
    //rt := d7024e.NewRoutingTable(me)
    //rt.AddContact(known)
    //kademlia := d7024e.NewKademlia(*rt, *network)
    //kademlia.LookupContact(&me)
    
    
   


    

 
  
}