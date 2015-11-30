package main 
import "net" 
import "fmt" 
import "bufio" 
import "strconv"
type myMsg struct {
	seqNum int
	message string
}
func main() {   // connect to this socket   
	conn, _ := net.Dial("tcp", "127.0.0.1:9999")   
	i :=0
	for ; i < 99; i++ {     // read in input from stdin     
		//reader := bufio.NewReader(os.Stdin)     
		  
		fmt.Fprintf(conn, strconv.Itoa(i) + "\n")     // listen for reply     
		fmt.Printf("sending  %s message\n", strconv.Itoa(i))
	}
	fmt.Printf("sent  %d message\n", i)
	for {     // read in input from stdin     
		   
		message, _ := bufio.NewReader(conn).ReadString('\n')     
		fmt.Print("Message from server: "+message)
		   
	}  
} 