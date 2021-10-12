package main

// Importing packages
import (
    "bufio"
    "encoding/hex"
    "fmt"
    "github.com/labstack/echo/v4"
    "github.com/wilkru-7/D7024E/d7024e"
    "math/rand"
    "net"
    "net/http"
    "os"
    "strings"
    "time"
)

/*
 * Main function. Creates instances of routing tables, networks, kademlias.
 * 
 */
func main() {
    rand.Seed(time.Now().UnixNano())
    conn,_ := net.Dial("ip:icmp","google.com")
    fmt.Println(conn.LocalAddr())
    
    if (conn.LocalAddr().String() == "172.21.0.2"){
        me := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), "172.21.0.2:8080") 
        rt := d7024e.NewRoutingTable(me)
        network := d7024e.NewNetwork(&me, rt)
        kademlia := d7024e.NewKademlia(*rt, network)
        go startAPI(kademlia)
        go startCLI(kademlia)
        network.Listen("8080")
    } else if (conn.LocalAddr().String() == "172.21.0.3") {
        known := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), "172.21.0.2:8080") 
        me := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000001"), "172.21.0.3:8080")
        rt := d7024e.NewRoutingTable(me)
        rt.AddContact(known)
        network := d7024e.NewNetwork(&me, rt)
        kademlia := d7024e.NewKademlia(*rt, network)
        go startAPI(kademlia)
        go kademlia.LookupContact(&me)
        go startCLI(kademlia)
        network.Listen("8080")
    } else {
        known := d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000000"), "172.21.0.2:8080") 
        me := d7024e.NewContact(d7024e.NewRandomKademliaID(), d7024e.GetLocalIP()+":8080")
        rt := d7024e.NewRoutingTable(me)
        rt.AddContact(known)
        network := d7024e.NewNetwork(&me, rt)
        kademlia := d7024e.NewKademlia(*rt, network)
        go startAPI(kademlia)
        go kademlia.LookupContact(&me)
        go startCLI(kademlia)
        network.Listen("8080")
    }
  
}

/*
 * Starts the command line interface for the nodes in the network.
 * 
 */
func startCLI(kademlia *d7024e.Kademlia) {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        text := strings.Split(scanner.Text(), " ")
        switch{
        case text[0] == "put" && len(text) == 2:
            src := []byte(text[1])
            dst := make([]byte, hex.EncodedLen(len(src)))
            hex.Encode(dst, src)
            keyString := string(dst)
            key := d7024e.NewKademliaID(keyString)
            fmt.Println("Hash of key is: ", key)
            kademlia.Store(key, text[1])
        case text[0] == "get" && len(text) == 2:
            result := kademlia.LookupData(text[1])
            if !strings.HasPrefix(result[0], "contact") {
                fmt.Println("Value is: ", result[0], " retrieved from: ", result[1])
            } else {
                fmt.Println("Value not found but here are the closest contacts: ", result)
            }
        case text[0] == "forget" && len(text) == 2:
            kademlia.Forget(text[1])
        case text[0] == "exit":
            os.Exit(0)
        default:
            fmt.Println("Command not supported, try put, get or exit")
        }
    }
}
  
/*
 * Starts the RESTful HTTP interface for the nodes in the network.
 * 
 */
func startAPI(kademlia *d7024e.Kademlia){
    e := echo.New()

	e.POST("/objects/:hash", func(c echo.Context) error {
        hash := c.Param("hash")
        src := []byte(hash)
        dst := make([]byte, hex.EncodedLen(len(src)))
        hex.Encode(dst, src)
        keyString := string(dst)
        key := d7024e.NewKademliaID(keyString)
        fmt.Println("Hash of key is: ", key)
        if (kademlia.Store(key, hash)){
            return c.String(http.StatusOK, "201 CREATED \n" + hash + "\nLocation: /objects/"+key.String() + "\n")
        }else{
            return c.String(http.StatusOK, "Store Failed\n")
        }
	})

    e.GET("/objects/:hash", func(c echo.Context) error {
        hash := c.Param("hash")
        result := kademlia.LookupData(hash)
        if !strings.HasPrefix(result[0], "contact") {
            return c.String(http.StatusOK, result[0] +"\n")
        } else {
            return c.String(http.StatusOK, "Value not found\n")
        }
	})

	e.Start(":8081")
}