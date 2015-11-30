// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0
import (
"bufio"
"fmt"
"net"
"strconv"
)

type Client struct {
	reader *bufio.Reader
	writer *bufio.Writer
	clientConn net.Conn
	clientChan chan string
}

type multiEchoServer struct {
	msgchan	 chan string
	quitchan chan bool
	ln       net.Listener
	clients  map[int](*Client)
}

func NewCli(conn net.Conn) *Client{
	return &Client {
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
		clientConn: conn,
		clientChan: make(chan string, 100),
	}
}
// Creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	server := multiEchoServer{
		msgchan	:	make(chan string),
		quitchan: 	make(chan bool),
		clients : 	map[int](*Client){},	
	}
	return &server
}
func (mes *multiEchoServer) broadcast() {
	for {
		select {
		case <- mes.quitchan:
			return
		case message := <- mes.msgchan

		}
	}
}
func listenClient(mes *multiEchoServer, cli *Client) {
	for {     // will listen for message to process ending in newline (\n)   
		message, err := cli.reader.ReadString('\n')     // output message received 
		if err != nil {
			// fmt.Print("conn error/close ", err)
			delete(mes.clients, index)
			return 
		}
		fmt.Println("Message received ", string(message))     // sample process for string received     
		// go func() {
		// 	if(len(cli.clientChan) <= 100) {
		// 		cli.clientChan <- string(message) + "\n"
		// 	}
		// }()
		// go func(){
		// 	cli.writer.WriteString(<-cli.clientChan)
		// 	cli.writer.Flush()
		// } ()
		mes.msgchan <- message  
	}
}

func (mes *multiEchoServer) Start(port int) error {

	// fmt.Println("Launching server...")   // listen on all interfaces
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Print("error when ln listenning ", err)
		return err
	}
	mes.ln = ln
	cliNum := 0
	go func() {
		for {
			select {
			case <-mes.quitchan:
				return
			case message := <- mes.msgchan:
				go mes.broadcast()
			default:
				conn, err := mes.ln.Accept()
				if err != nil {
					fmt.Println("Couldn't accept: ", err)
					continue
				}
				cli := NewCli(conn)
				mes.clients[cliNum++] = cli
				go listenClient(mes, cli)
			}
		}
	} ()
	return nil
}

func CloseAllClients(mes *multiEchoServer) {
	for _, client := range mes.clients {
		client.clientConn.Close()
	} 
}
func (mes *multiEchoServer) Close() {
	go func() {
		mes.quitchan <- true
	} ()
	mes.ln.Close()
	go CloseAllClients(mes)
	return
}

func (mes *multiEchoServer) Count() int {
	return len(mes.clients)
}

// TODO: add additional methods/functions below!