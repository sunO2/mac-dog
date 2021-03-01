package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"./util"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/wake", wake)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

///远程开机
func wake(write http.ResponseWriter, request *http.Request) {
	mac_address := request.FormValue("mac_address")
	fmt.Println("远程开机：" + mac_address)
	output, err := util.Wake(mac_address)
	if nil != err {
		fmt.Fprintln(write, err.Error())
	} else {
		fmt.Fprintln(write, output)
	}

}

func index(write http.ResponseWriter, request *http.Request) {
	wake := request.FormValue("wake")
	fmt.Println("启动主程序" + wake)
	content, _ := ioutil.ReadFile("./html/index.html")
	write.Write(content)
	list, err := util.GetMacList()
	if nil != err {
		fmt.Fprintln(write, "<p>"+err.Error()+"</></br>")
		// return
	}
	list = append(list, "80:80:80:80:80")
	list = append(list, "90:90:90:90:90")
	if len(list) > 0 {
		for i := 0; i < len(list); i++ {
			// fmt.Fprintln(write, list[i])
			fmt.Fprintln(write, fmt.Sprintf("<a href=javascript:wake('%s')>%s</></br>", list[i], list[i]))
		}
	} else {
		fmt.Fprintln(write, "列表为空")
	}
}
