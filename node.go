package main

// Importing packages
import (
	"github.com/wilkru-7/D7024E/d7024e"
    //"fmt"
    //"sort"
    //"strings"
    //"time"
)
  
// Main function
func main() {
    
    network := d7024e.NewNetwork()
    //contact := d7024e.NewContact(d7024e.NewRandomKademliaID(), "google.com")
    
    known := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), "172.20.0.2:8000") 
    me := d7024e.NewContact(d7024e.NewRandomKademliaID(), "localhost:8000")
    //network.contact = d7024e.NewContact(d7024e.NewRandomKademliaID(), "localhost:80")
    //network.SendPingMessage(&contact)
    rt := d7024e.NewRoutingTable(me)
    rt.AddContact(known)
    //kademlia := d7024e.NewKademlia(*rt, *network)
   //kademlia.LookupContact(&me)
    d7024e.Listen("127.0.0.1", 8080)
    network.SendPingMessage(&known)
    /* for {
        fmt.Println(me.String())
        //d7024e.Listen("0.0.0.0", 8080)
        if !true {
            break
        }
    } */
    


    

 
  
}