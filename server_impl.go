// Implementation of a MultiEchoServer. 
//author: Harry Ge
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
	quitchan chan bool
	clientChan chan string
}

type multiEchoServer struct {
	msgchan	 chan string
	quitchan chan bool
	stopchan	chan bool
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
		stopchan: 	make(chan bool),
		clients : 	map[int](*Client){},	
	}
	return &server
}
func (mes *multiEchoServer) broadcast(message string) {
	// fmt.Println("broadcast ", message)
	for _, client := range mes.clients {
		if(len(client.clientChan)<=100) {
			client.clientChan <- message
		}
	}
}
func runClient(mes *multiEchoServer, id int) {
	cli := mes.clients[id]
	go func() {	
		for {     // will write message into client 
			select {
			case <- cli.quitchan:
				delete(mes.clients, id)
				return
			case msg := <- cli.clientChan:
				// fmt.Println("writing ", msg)
				cli.writer.WriteString(msg)
				cli.writer.Flush()

			}  
		}
	}()
	
	//will listen for message to process ending in newline (\n) 
		for{	
			message, err := cli.reader.ReadString('\n')     // output message received 
			if err != nil {
				// fmt.Print("conn error/close ", err)
				delete(mes.clients, id)
				return 
			}
			//fmt.Println("Message received ", string(message))     // sample process for string received     
			go func() {
				mes.msgchan <- string(message) 
			}()
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
				case <- mes.stopchan:
				return
			default:
			conn, err := mes.ln.Accept()
				if err != nil {
					fmt.Println("Couldn't accept: ", err)
					continue
				}
				mes.clients[cliNum] = NewCli(conn)
				go runClient(mes, cliNum)
				cliNum++
				
			}
		}
	}()
	go func() {
		for {
			select {
			case <-mes.quitchan:
				return
			case message := <- mes.msgchan:
				// fmt.Println("need broadcasting ", message)
				go mes.broadcast(message)
			}
		}
	} ()
	return nil
}

func CloseAllClients(mes *multiEchoServer) {
	for _, client := range mes.clients {
		client.quitchan <- true
		client.clientConn.Close()
	} 
}
func (mes *multiEchoServer) Close() {
	go func() {
	mes.quitchan <- true
	mes.stopchan <- true
	} ()
	mes.ln.Close()
	go CloseAllClients(mes)
	return
}

func (mes *multiEchoServer) Count() int {
	return len(mes.clients)
}
