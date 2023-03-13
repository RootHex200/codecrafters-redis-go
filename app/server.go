package main

import (
	"fmt"
	"strings"
	// "strings"

	"regexp"

	// Uncomment this block to pass the first stage
	"net"
)


func nilcheck(err error){
	if err != nil {
		fmt.Printf("Failed to bind to port 6379%s",err.Error())
		return
	}
}


func parsedata(data string)([]string){
	re := regexp.MustCompile(`\$\d+\r\n([^\r\n]+)\r\n`)
    input := data
	results:=make([]string,0)
	
    matches := re.FindAllStringSubmatch(input, -1)
    for _, match := range matches {
		results = append(results, match[1])
    }
	fmt.Printf("results:%s",results)
	// joined := strings.Join(results, " ")
	// results = results[:0]
	return results
}

type command struct {
	arguments map[string]string
	options map[string]string
}
type Strorage struct {
	data map[string]string
	createAt map[string]int64
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	storage:=Strorage{data: make(map[string]string),createAt: make(map[string]int64)}
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	
	nilcheck(err)
	defer l.Close()
	for {
		conn, err := l.Accept()
		nilcheck(err)
		go func (conn net.Conn){
			buf:=make([]byte, 1024)
			for{
				n, err := conn.Read(buf)
				if err != nil {
					fmt.Printf("Failed to read from connection%s",err.Error())
					return
				}
				cmd:=parsedata(string(buf[:n]))
				fmt.Println("this is cmd:",len(cmd))
				fmt.Println("this is cmd[0]:",strings.ToUpper(cmd[0]))

				switch strings.ToUpper(cmd[0]) {
					case "PING":
						if(len(cmd)==1){
							conn.Write([]byte("+PONG\r\n"))
							
						}else{
							conn.Write([]byte("+wrong command\r\n"))
						}
						
						
					case "ECHO":
						if(len(cmd)==2){
							conn.Write([]byte("+"+cmd[1]+"\r\n"))							
						}else{
							conn.Write([]byte("+wrong command\r\n"))
						}
						
						
					case "SET":
						if(len(cmd)==3){
							storage.data[cmd[1]]=cmd[2]
							conn.Write([]byte("+OK\r\n"))
							fmt.Println(storage.data)
						}else{
							conn.Write([]byte("+wrong command\r\n"))
						}
						
						
					case "GET":
						findvalue:=make([]string,0)
						fmt.Println(storage.data)
						if(len(cmd)==2){
							for key, value := range storage.data {
								if(key==cmd[1]){
									// conn.Write([]byte("+"+value+"\r\n"))
									findvalue = append(findvalue, value)
									break
								}
								
						}
						if(len(findvalue)==0){
							conn.Write([]byte("$-1\r\n"))
						}else{
							conn.Write([]byte("+"+findvalue[0]+"\r\n"))
						}
					}else{
						conn.Write([]byte("+wrong command\r\n"))
					}
					default:
					conn.Write([]byte("-ERR unknown command '"+cmd[0]+"'\r\n"))
				}
				
			}
			
			
		}(conn)
	}
}
