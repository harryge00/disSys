package main 
import "net" 
import "fmt" 
import "bufio" 
import "strings" // only needed below for sample processing 
func serve(conn net.Conn){
	fmt.Print("Accept conn :", conn) 
	for {     // will listen for message to process ending in newline (\n)     
		message, err := bufio.NewReader(conn).ReadString('\n')     // output message received 
		if err != nil {
			fmt.Print("conn error/close ", err)
			return 
		}
		fmt.Print("Message received ", string(message))     // sample process for string received     
		newmessage := strings.ToUpper(message)     // send new string back to client     
		conn.Write([]byte(newmessage + "\n"))   
	}
}
func main() {   
	fmt.Println("Launching server...")   // listen on all interfaces   
	ln, _ := net.Listen("tcp", ":8081")   // accept connection on port   
	for {
		conn, err := ln.Accept()   // run loop forever (or until ctrl-c)
		if err != nil {
			fmt.Print("error when ln accepting ", err)
			
		} else {
			go serve(conn)
		}
	}
} 