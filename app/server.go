package main

import (
	"fmt"
	"strings"

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


func parsedata(data string)(string){
	re := regexp.MustCompile(`\$\d+\r\n([^\r\n]+)\r\n`)
    input := data
	var results []string
	
    matches := re.FindAllStringSubmatch(input, -1)

    for _, match := range matches {
		results = append(results, match[1])
    }
	joined := strings.Join(results, " ")
	results = results[:0]
	return joined
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

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
				
				if(strings.ToUpper(cmd)=="PING"){
					conn.Write([]byte("+PONG\r\n"))
				}else{
					conn.Write([]byte(fmt.Sprintf("+%s\r\n",cmd)))
				}
				
			}
			
			
		}(conn)
	}
}
