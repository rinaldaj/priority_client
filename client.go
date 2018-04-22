package main

import (
	"fmt"
	"os"
	"bufio"
	"net"
	"strings"
	"encoding/json"
	"strconv"
)

type task struct{
	User    string
	Importance  int
	Task    string
	Duedate string
	Duetime string
	Completiontime  int
	Priority	int
}


func main(){
	args := os.Args
	config, err := os.Open(".config")
	ip := ""
	port := ":"
	username := ""
	if (err != nil){
		fmt.Println("Config file not found beggining intitial set up")
		fmt.Println("Enter Server Ip")
		ip,_ = bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Println("Enter Port")
		port,_ = bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Println("enter user name")
		username,_ = bufio.NewReader(os.Stdin).ReadString('\n')
		fiOut, _ := os.Create(".config")
		fmt.Fprint(fiOut,ip)
		fmt.Fprint(fiOut,port)
		fmt.Fprint(fiOut,username)
	} else {
		configReader := bufio.NewReader(config)
		ip,_ = configReader.ReadString('\n')
		port,_ = configReader.ReadString('\n')
		username,_ = configReader.ReadString('\n')
	}
	port = ":" + port
	server,errs := net.Dial("tcp","localhost:6666")
	if (errs != nil){
		fmt.Println(errs)
		return
	}
	defer server.Close()
	if (len(args) < 2){
		fmt.Fprint(server,"3\n")
		recp,_ := bufio.NewReader(server).ReadString('\n')
		strs := strings.Split(recp,"{")
		var jason = make([]task,0)
		for i := 0;i<len(strs);i++{
			tmp := task{}
			strs[i] = "{" + strs[i]
			json.Unmarshal([]byte(strs[i]),&tmp)
			//Calculate Priority here
			jason = append(jason,tmp)
		}
		fmt.Println(jason)
	} else {
		if (args[1] == "-a"){
		fmt.Println("Adding assignment")
		buffy := bufio.NewReader(os.Stdin)
		curTask := task{}
		fmt.Println("Assignment Name")
		curTask.Task,_ = buffy.ReadString('\n')
		curTask.Task = curTask.Task[0:len(curTask.Task)-1]
		fmt.Println("Importance")
		tmpy,_:= buffy.ReadString('\n')
		curTask.Importance,_ = strconv.Atoi(tmpy[0:len(tmpy)-1])
		fmt.Println("Duedate")
		curTask.Duedate,_ = buffy.ReadString('\n')
		curTask.Duedate = curTask.Duedate[0:len(curTask.Duedate)-1]
		fmt.Println("Duetime")
		curTask.Duetime,_ = buffy.ReadString('\n')
		curTask.Duetime = curTask.Duetime[0:len(curTask.Duetime)-1]
		fmt.Println("Completiontime")
		tmpy,_= buffy.ReadString('\n')
		curTask.Completiontime,_ = strconv.Atoi(tmpy[0:len(tmpy)-1])
		curTask.User = username[0:len(username)-1]
		sje,_ := json.Marshal(curTask)
		sendString := "1"+string(sje)+"\n"
		fmt.Fprint(server,sendString)
		} else if (args[1] == "-d"){
		fmt.Fprint(server,"3\n")
		recp,_ := bufio.NewReader(server).ReadString('\n')
		strs := strings.Split(recp,"{")
		server2,eirrs := net.Dial("tcp","localhost:6666")
		if (eirrs != nil){
			return
		}
		defer server2.Close()
		for i := 0;i<len(strs);i++{
			tmp := task{}
			strs[i] = "{" + strs[i]
			json.Unmarshal([]byte(strs[i]),&tmp)
			fmt.Printf("%d: %s\n",i,tmp.Task)
		}
		fmt.Println("Which do you want to delete?")
		delS,_ := bufio.NewReader(os.Stdin).ReadString('\n')
		delIndex,_ := strconv.Atoi(delS[0:len(delS)-1])
		sendStr:= "2" + strs[delIndex]+"\n"
		fmt.Fprint(server2,sendStr)
		}
	}
}
