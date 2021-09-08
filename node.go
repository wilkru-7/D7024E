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

    network := new(d7024e.Network)
    contact := d7024e.NewContact(d7024e.NewRandomKademliaID(), "google.com")
    //network.contact = d7024e.NewContact(d7024e.NewRandomKademliaID(), "localhost:80")
    network.SendPingMessage(&contact)
 
  
}