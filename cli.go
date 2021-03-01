package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

var device *string = flag.String("d", "mac", "操作的设备 mac/nas")

func init() {
	flag.Usage = usage
}

var usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "tools 使用功能\n")
	fmt.Fprintf(flag.CommandLine.Output(), "action 功能:\n")
	fmt.Fprintf(flag.CommandLine.Output(), "	ssh ssh连接\n")
	fmt.Fprintf(flag.CommandLine.Output(), "	wake 远程开机\n")
	flag.PrintDefaults()
}

func Parse(point int){
	flag.CommandLine.Parse(os.Args[point:])
}

func main() {
	router:=os.Args[1]
	Parse(2)
	fmt.Println("输入的参数：" + router)
	uname,ip,_:=GetInfo(fmt.Sprintf("%s",*device))
    switch {
		case router == "ssh":
			SshLogin(uname,ip,22)
			return
		case router == "wake":
			return
		default:
			usage()
	}
}

func GetInfo(channel string)(string,string,string){
	fmt.Println("获取渠道" + channel)
	switch{
		case channel == "nas":
			return os.Getenv("NAS_UNAME"),os.Getenv("NAS_IP"),os.Getenv("NAS_MAC")
		case channel == "mac":
			fmt.Println("获取渠道" + channel)
			return os.Getenv("MINI_UNAME"),os.Getenv("MINI_IP"),os.Getenv("MINI_MAC")
	}
	return os.Getenv("MINI_UNAME"),os.Getenv("MINI_IP"),os.Getenv("MINI_MAC")
}

/**
 ssh 连接  
 unam: 用户名
 ip: ip地址
 port: 端口号
*/
func SshLogin(uname string,ip string,port int){
	printError:=func(err error,mesg string){
		if nil != err{
			fmt.Println(mesg,err)
		}
	}
	fmt.Println("请输入密码")
	var password string
	fmt.Scanln(&password)
	clinet,err:=ssh.Dial("tcp",fmt.Sprintf("%s:%d",ip,port),&ssh.ClientConfig{
		User: uname,
		Auth: []ssh.AuthMethod{ssh.Password(fmt.Sprintf(password))},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			fmt.Println("ssh 登录：address:" + remote.String())
			return nil
		},
	})
	printError(err,"ssh dial 错误")
	session,err :=clinet.NewSession()
	printError(err,"ssh 创建错误 错误")
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	modes := ssh.TerminalModes{
		ssh.ECHO: 0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("linux",32,160,modes)
	printError(err,"requstsPty")
	err = session.Shell()
    printError(err, "start shell")
	err = session.Wait()
	printError(err, "return")

}
