package main

import (
	"fmt"
	"os"
	"bufio"
	"net"
)

type task struct{
	User    string
	Importance  int
	Task    string
	Duedate string
	Duetime string
	Completiontime  int
}

func main(){
	args := os.Args
	config, err := os.Open(".config")
	ip := ""
	port := ":"
	if (err != nil){
		fmt.Println("Config file not found beggining intitial set up")
		fmt.Println("Enter Server Ip")
		ip,_ = bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Println("Enter Port")
		port,_ = bufio.NewReader(os.Stdin).ReadString('\n')
		fiOut, _ := os.Create(".config")
		fmt.Fprint(fiOut,ip)
		fmt.Fprint(fiOut,port)
	} else {
		configReader := bufio.NewReader(config)
		ip,_ = configReader.ReadString('\n')
		port,_ = configReader.ReadString('\n')
	}
	port = ":" + port
	server,errs := net.Dial("tcp","localhost:6666")
	//defer server.Close()
	if (errs != nil){
		fmt.Println(errs)
		return
	}
	if (len(args) < 2){
		fmt.Fprint(server,"3\n")
		strs,_ := bufio.NewReader(server).ReadString('\n')
		server.Close()
		fmt.Println(strs)
	}
}
