// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0
import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

var numCli chan int
func addCli() {
	x := <- numCli 
	x++
	fmt.Print("numCli added ",x)
	numCli <- x
}

func rmCli() {
	x := <- numCli
	x--
	if(x < 0) {
		x = 0
	}
	fmt.Print("numCli removed ",x)
	numCli <- x
}
func initCount() {
	numCli <- 0
}



type Client struct {
	clientConn net.Conn
	clientChan chan string
}

type multiEchoServer struct {
	msgchan  chan string
	addchan  chan Client
	rmchan   chan net.Conn
	quitchan chan bool
	countReq chan bool
	countRcv chan int
	ln       net.Listener
	clients  map[net.Conn](chan string)
}

// Creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {

	server := multiEchoServer{
		msgchan:  make(chan string),
		addchan:  make(chan Client),
		rmchan:   make(chan net.Conn),
		quitchan: make(chan bool),
		countReq: make(chan bool),
		countRcv: make(chan int),
		clients:  make(map[net.Conn](chan string))}
	return &server
}

func serve(conn net.Conn){
	go addCli()
	var buf = make(chan string, 100)
	fmt.Printf("%s, Accept conn \n", conn) 
	for {     // will listen for message to process ending in newline (\n)     
		message, err := bufio.NewReader(conn).ReadString('\n')     // output message received 
		if err != nil {
			fmt.Print("conn error/close ", err)
			go rmCli()
			return 
		}
		fmt.Print("Message received ", string(message))     // sample process for string received     
		if(len(buf) <= 100) {
			buf <- string(message) + "\n"
		}
		conn.Write([]byte(<-buf))   
	}
}

func listenning (port string) {
	ln, _ := net.Listen("tcp", ":"+port)
	for {
		conn, err := ln.Accept()   // run loop forever (or until ctrl-c)
		if err != nil {
			fmt.Print("error when ln accepting ", err)
			
		} else {
			go serve(conn)
		}
	}
}
func (mes *multiEchoServer) Start(port int) error {
	go initCount()   
	fmt.Println("Launching server...")   // listen on all interfaces
	go listenning(strconv.Itoa(port))
	return nil
}

func (mes *multiEchoServer) Close() {
	close(mes.quitchan)
	mes.ln.Close()

	return
}

func (mes *multiEchoServer) Count() int {
	x := <- numCli
	numCli <- x
	return x
}

// TODO: add additional methods/functions below!