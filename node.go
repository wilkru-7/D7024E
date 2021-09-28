package main

// Importing packages
import (
	"github.com/wilkru-7/D7024E/d7024e"
    "net"
    "fmt"
    "os"
    "bufio"
    "strings"
    "encoding/hex"
   /*  "reflect" */
)
func scanner(kademlia *d7024e.Kademlia) {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        text := strings.Split(scanner.Text(), " ")
        switch{
        case text[0] == "put" && len(text) == 2:
            
            src := []byte(text[1])

            dst := make([]byte, hex.EncodedLen(len(src)))
            hex.Encode(dst, src)
            keyString := string(dst)
            /* fmt.Println("Hash of key is: ", dst) */
            key := d7024e.NewKademliaID(keyString)
            fmt.Println("Hash of key is: ", key)
            kademlia.Store(key, text[1])
            //Expect result??
        case text[0] == "get" && len(text) == 2:
            /* src := []byte(text[1])

            dst := make([]byte, hex.EncodedLen(len(src)))
            hex.Encode(dst, src)
            keyString := string(dst) */
            /* key := d7024e.NewKademliaID(keyString) */
            value := kademlia.LookupData(text[1])
            fmt.Println("Value is: ", value)
            //Print content and node it was retreived from

        case text[0] == "exit":
            break;
        default:
            fmt.Println("Command not supported, try put, get or exit")
        }
        break;
    }
  }
  
// Main function
func main() {
    
    conn,_ := net.Dial("ip:icmp","google.com")
    fmt.Println(conn.LocalAddr())

    /* go scanner() */
    if (conn.LocalAddr().String() == "172.22.0.2"){
        me := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), "172.22.0.2:8080") 
        rt := d7024e.NewRoutingTable(me)
        network := d7024e.NewNetwork(&me, rt)
        kademlia := d7024e.NewKademlia(*rt, network)
        go scanner(kademlia)
        network.Listen("127.0.0.1", "8080")
    } else if (conn.LocalAddr().String() == "172.22.0.3") {
        known := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), "172.22.0.2:8080") 
        me := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000001"), "172.22.0.3:8080")
        rt := d7024e.NewRoutingTable(me)
        rt.AddContact(known)
        network := d7024e.NewNetwork(&me, rt)
        kademlia := d7024e.NewKademlia(*rt, network)
        go kademlia.LookupContact(&me)
        go scanner(kademlia)
        network.Listen("127.0.0.1", "8080")
    } else {
        known := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), "172.22.0.2:8080") 
        me := d7024e.NewContact(d7024e.NewRandomKademliaID(), d7024e.GetLocalIP()+":8080")
        rt := d7024e.NewRoutingTable(me)
        rt.AddContact(known)
        network := d7024e.NewNetwork(&me, rt)
        kademlia := d7024e.NewKademlia(*rt, network)
        go kademlia.LookupContact(&me)
        go scanner(kademlia)
        network.Listen("127.0.0.1", "8080")
    }
  
}