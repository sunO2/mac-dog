package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
)

//入职输出
var loger *log.Logger

func init(){
  user,uerr:=user.Current()
  var homeDir string
  if nil != uerr{
    loger.Println(uerr)
   return
  }

  homeDir = user.HomeDir

  log.SetPrefix("MAC_DOG: ")
  log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
  file := homeDir + "/develop/mac_dog/mac_dog.log"
  logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 6)
  if nil != err{
    fmt.Println(err)
  }
  loger = log.New(logFile,"[MAC_DOG]",log.LstdFlags|log.Lshortfile)
}

///离线设备
func offline(offlineMac string){
  loger.Println("离线设备：" + offlineMac)
}

///在线设备
func online(onlineMac string){
  loger.Println("上线设备：" + onlineMac)
}

func listString(list []string) string{
  var stringlist string
  for i:=0;i<len(list);i++{
     if i == 0{
       stringlist = fmt.Sprintf("%s",list[i])
     }else{
       stringlist = fmt.Sprintf("%s,%s",stringlist,list[i])
     }
  }
  return stringlist
}

func contain(maclist []string,contain string) bool {
  for i := 0; i< len(maclist); i++{
    if maclist[i] == contain{
      return true
    }
  }
  return false
}

func containList(list list.List,contain string) bool {
   for e :=list.Front(); e!=nil; e=e.Next(){
      if contain == e.Value{
        return true
      }
   }
   return false
}

func main(){
  var err error
  var whoami []byte
  var ip string
  if len(os.Args) <= 1{
     ip = "192.168.0.*"
  }else {
     ip = os.Args[1]
  }
  loger.Println("正在扫描： " + ip)
  cmd:=exec.Command("/bin/sh","-c","sudo nmap -sP " + ip)
  if whoami,err = cmd.Output(); err != nil{
    loger.Println(err)
    return
  }
  user,uerr:=user.Current()
  var homeDir string
  if nil != uerr{
    loger.Println(uerr)
   return
  }

  homeDir = user.HomeDir
  file,e := os.OpenFile(homeDir+"/.maclist",os.O_CREATE|os.O_RDWR,0666)
  if e != nil {
    loger.Println("读取文件：文件打开失败",e)
    return
  }
  defer file.Close()
  reader := bufio.NewReader(file)
  buf,err := reader.ReadBytes('\n')
  if err != nil {
    if err != io.EOF{
      loger.Println(" error = ",err)
     return
    }
  }
  //文件中mac 地址
  fileMacstring := string(buf)
  fileMacList := strings.Split(fileMacstring,",")
  reg:=regexp.MustCompile("[A-F\\d]{2}:[A-F\\d]{2}:[A-F\\d]{2}:[A-F\\d]{2}:[A-F\\d]{2}:[A-F\\d]{2}")
  if reg==nil {
    loger.Println("正则异常")
    return
  }
  result:=reg.FindAllStringSubmatch(string(whoami),-1)
  var MACList []string
  for i:=0;i<len(result);i++{
    deviceMac := result[i][0]
    MACList = append(MACList,deviceMac)
    if !contain(fileMacList,deviceMac){
          online(deviceMac)
    }
  }
  loger.Println("")
  loger.Println("")
  loger.Println("")
  for ii := 0; ii < len(fileMacList); ii++ {
    if !contain(MACList,fileMacList[ii]){
       offline(fileMacList[ii])
    }
  }

  maclist := listString(MACList)
  loger.Println("列表文件已保存")

   file,e = os.OpenFile(homeDir+"/.maclist",os.O_TRUNC|os.O_CREATE|os.O_RDWR,0666)
   if e != nil {
    loger.Println("写入文件：文件打开失败",e)
    return
  }
  defer file.Close()

  write := bufio.NewWriter(file)
  write.WriteString(maclist)
  write.Flush()
}
