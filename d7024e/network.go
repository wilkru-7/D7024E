package d7024e

import(
	"net"
	"fmt"
	"time"
	"net/http"
	//"os"

	"github.com/labstack/echo/v4"
	//"github.com/labstack/echo/v4/middleware"
	
)
type Network struct {
	//contact *Contact
}

// NewRoutingTable returns a new instance of a RoutingTable
func NewNetwork() *Network {
	network := &Network{}
	return network
}

func Listen(ip string, port int) {
	e := echo.New()
	// TODO
	fmt.Println("We are in the listening")
	e.GET("/ping", func(c echo.Context) error {

		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})
	e.Start(":8080")
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		// handle error
	}
	for {
		//conn, err := ln.Accept()
		_, err := ln.Accept()
		if err != nil {
			// handle error
		}
		//go handleConnection(conn)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
	
	/* conn, err1 := net.Dial("udp", contact.Address)
	if err1 != nil {
	}
	fmt.Println("YES WE ARE PINGING")
	fmt.Fprintf(conn, "Hi UDP Server, How are you doing?") */
	
	//status, err := bufio.NewReader(conn).ReadString('\n')
	//status, err := bufio.NewReader(conn).ReadString('\n')

    //port := "4000"
    timeout := time.Duration(1 * time.Second)
    _, err := net.DialTimeout("udp", contact.Address, timeout)
    if err != nil {
        fmt.Printf("%s %s %s\n", contact.Address, "not responding", err.Error())
    } else {
        fmt.Printf("%s %s \n", "responding on adress:", contact.Address)
    }
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO

	fmt.Println("Sending find contact message not implemented :( ")
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
