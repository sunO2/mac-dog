package util

import (
	"os/exec"
	"regexp"
)

func GetMacList()([]string,error){
	var err error
	var whoami []byte
	cmd:=exec.Command("/bin/sh","-c","sudo nmap -sP 192.168.43.*")
	if whoami,err = cmd.Output(); err != nil{
	  return nil,err
	}
	 reg:=regexp.MustCompile("[A-F\\d]{2}:[A-F\\d]{2}:[A-F\\d]{2}:[A-F\\d]{2}:[A-F\\d]{2}:[A-F\\d]{2}")
	 if reg==nil {
	   return nil,err
	 }
	 result:=reg.FindAllStringSubmatch(string(whoami),-1)
	 var MACList []string
	 for i:=0;i<len(result);i++{
	  deviceMac := result[i][0]
	  MACList = append(MACList,deviceMac)
	 }
	 return MACList,nil
  }

  func Wake(deviceMac string)(string,error){
	cmd:=exec.Command("/bin/sh","-c","sudo etherwake -b -i wlan0 "+deviceMac)
	var output []byte
	var err error
	output,err = cmd.Output()
	if nil != err{
		return "",err
	}
	return string(output),nil
}