package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
)

///离线设备
func offline(offlineMac string){
  fmt.Println("离线设备：" + offlineMac)
}

///在线设备
func online(onlineMac string){
  fmt.Println("上线设备：" + onlineMac)
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
  fmt.Println("正在扫描： " + ip)
  cmd:=exec.Command("/bin/sh","-c","nmap -sP " + ip)
  if whoami,err = cmd.Output(); err != nil{
    fmt.Println(err)
    return
  }
  user,uerr:=user.Current()
  var homeDir string
  if nil != uerr{
   fmt.Println(uerr)
   return
  }

  homeDir = user.HomeDir
  file,e := os.OpenFile(homeDir+"/.maclist",os.O_CREATE|os.O_RDWR,0666)
  if e != nil {
    fmt.Println("读取文件：文件打开失败",e)
    return
  }
  defer file.Close()
  reader := bufio.NewReader(file)
  buf,err := reader.ReadBytes('\n')
  if err != nil {
    if err != io.EOF{
     fmt.Println(" error = ",err)
     return
    }
  }
  //文件中mac 地址
  fileMacstring := string(buf)
  fileMacList := strings.Split(fileMacstring,",")
  reg:=regexp.MustCompile("[A-F\\d]{2}:[A-F\\d]{2}:[A-F\\d]{2}:[A-F\\d]{2}:[A-F\\d]{2}:[A-F\\d]{2}")
  if reg==nil {
    fmt.Println("正则异常")
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
   fmt.Println("")
   fmt.Println("")
   fmt.Println("")
  for ii := 0; ii < len(fileMacList); ii++ {
    if !contain(MACList,fileMacList[ii]){
       offline(fileMacList[ii])
    }
  }

  maclist := listString(MACList)
  fmt.Println("列表文件已保存")

   file,e = os.OpenFile(homeDir+"/.maclist",os.O_TRUNC|os.O_CREATE|os.O_RDWR,0666)
   if e != nil {
    fmt.Println("写入文件：文件打开失败",e)
    return
  }
  defer file.Close()

  write := bufio.NewWriter(file)
  write.WriteString(maclist)
  write.Flush()
}
