package main

import (
	"fmt"
	"os"
	"os/user"
	"bufio"
	"net"
	"strings"
	"encoding/json"
	"strconv"
	"time"
	"sort"
	"io/ioutil"
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

type byPriority []task
func (a byPriority) Len() int	{return len(a)}
func (a byPriority) Swap(i,j int) {a[i],a[j] = a[j],a[i]}
func (a byPriority) Less(i,j int) bool {return a[i].Priority < a[j].Priority}

func save(a []task){
	usr,systouch := user.Current()
	if(systouch != nil){
		return
	}
	file,_ := os.Create(usr.HomeDir+"/.priority.json")
	for i :=0;i<len(a);i++{
		jss,_:= json.Marshal(a[i])
		fmt.Fprintln(file,string(jss))
	}
}
func sendall(a []task,port string,ip string){
	servs,errs := net.Dial("tcp",ip+port)
	if (errs != nil){
		return
	}
	defer servs.Close()
	sendStr := "1"
	for i:=0; i< len(a);i++{
		jss,_ := json.Marshal(a[i])
		sendStr = sendStr + string(jss)
	}
	sendStr = sendStr + "\n"
	fmt.Fprint(servs,sendStr)
}

func eq(x task,y task) bool{
 return strings.EqualFold(x.User, y.User) && x.Importance == y.Importance && strings.EqualFold(x.Task, y.Task) && strings.EqualFold(x.Duedate,y.Duedate) && strings.EqualFold(x.Duetime,y.Duetime) && x.Completiontime == y.Completiontime
}

func main(){
	args := os.Args
	usr,systouch := user.Current()
	if(systouch != nil){
		return
	}
	config, err := os.Open(usr.HomeDir + "/.priority.conf")
	ip := ""
	port := ":"
	username := ""
	jason := make([]task,0)
	if (err != nil){
		fmt.Println("Config file not found beggining intitial set up")
		fmt.Println("Enter Server Ip")
		ip,_ = bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Println("Enter Port")
		port,_ = bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Println("enter user name")
		username,_ = bufio.NewReader(os.Stdin).ReadString('\n')
		fiOut, _ := os.Create(usr.HomeDir + "/.priority.conf")
		fmt.Fprint(fiOut,ip)
		fmt.Fprint(fiOut,port)
		fmt.Fprint(fiOut,username)
	} else {
		configReader := bufio.NewReader(config)
		ip,_ = configReader.ReadString('\n')
		port,_ = configReader.ReadString('\n')
		username,_ = configReader.ReadString('\n')
	}
	username = username[0:len(username)-1]
	port = ":" + port[0:len(port)-1]
	ip = ip[0:len(ip)-1]
	server,errs := net.Dial("tcp",ip+port)
	strs := make([]string,0)
	if (errs == nil){
		defer server.Close()
		fmt.Fprint(server,"3\n")
		recp,_ := bufio.NewReader(server).ReadString('\n')
		strs = strings.Split(recp,"{")
	}
		datos, _ :=ioutil.ReadFile(usr.HomeDir + "/.priority.json")
		trinkit := strings.Split(string(datos),"{")
		strs = append(strs,trinkit...)
		for i := 0;i<len(strs);i++{
			tmp := task{}
			strs[i] = "{" + strs[i]
			json.Unmarshal([]byte(strs[i]),&tmp)
			//Calculate Priority here
			dues := tmp.Duedate +" "+ tmp.Duetime
			tms,_ := time.Parse("1-2-06 3:04pm MST",dues)
			cur := time.Now()
			durs := tms.Sub(cur.Add(time.Duration(tmp.Completiontime*int(time.Hour)*tmp.Importance)))
			tmp.Priority = int(durs)
			dup := false
			for k:=0;k<len(jason);k++ {
				if(eq(jason[k],tmp)){
					dup = true
					break
				}
			}
			if (!dup){
				jason = append(jason,tmp)
			}
		}
		sort.Sort(byPriority(jason))
		server2,eirrs := net.Dial("tcp",ip+port)
		if (eirrs == nil){
			defer server2.Close()
		}
	if (len(args) < 2){
		trace := 0
		fmt.Println("#:\tTask\tDate\tTime\tImportance\tHours")
		for i :=0;i<len(jason);i++{
			if(strings.EqualFold(jason[i].User,username)){
				fmt.Printf("%d:\t%s\t%s\t%s\t%d\t%d\n",trace,jason[i].Task,jason[i].Duedate,jason[i].Duetime,jason[i].Importance,jason[i].Completiontime)
				trace++
			}
		}
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
		curTask.Duetime = curTask.Duetime[0:len(curTask.Duetime)-1] + " EDT"
		fmt.Println("Completiontime")
		tmpy,_= buffy.ReadString('\n')
		curTask.Completiontime,_ = strconv.Atoi(tmpy[0:len(tmpy)-1])
		curTask.User = username
		jason = append(jason,curTask)
		if(eirrs == nil){
			sje,_ := json.Marshal(curTask)
			sendString := "1"+string(sje)+"\n"
			fmt.Fprint(server2,sendString)
		}
		} else if (args[1] == "-d"){
		for i := 0;i<len(jason);i++{
			fmt.Printf("%d: %s\n",i,jason[i].Task)
		}
		fmt.Println("Which do you want to delete?")
		delS,_ := bufio.NewReader(os.Stdin).ReadString('\n')
		delIndex,_ := strconv.Atoi(delS[0:len(delS)-1])
		tmpr := make([]task,0)
		storss,_ := json.Marshal(jason[delIndex])
		for j := 0;j<len(jason);j++{
			if(j!=delIndex){
				tmpr = append(tmpr,jason[j])
			}
		}
		jason = tmpr
		if(eirrs == nil){
			sendStr:= "2" + string(storss) +"\n"
			fmt.Fprint(server2,sendStr)
		}
		}
	}
	save(jason)
	sendall(jason,port,ip)
}
