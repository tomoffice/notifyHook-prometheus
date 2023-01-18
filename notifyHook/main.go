package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	PkgApi "notifyHook/api"
	PkgJson "notifyHook/jsonWrap"
	PkgLine "notifyHook/line"
)

func alertsHandle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	if r.ContentLength == 0 {
		fmt.Println("got zero bytes")
		return
	}
	var jsonNotify PkgJson.Jsoner = PkgJson.NewNotify()
	err = jsonNotify.Decode(string(body))
	if err != nil {
		fmt.Println(err)
		return
	}
	jPrint, err := json.MarshalIndent(jsonNotify, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	var lineNotify PkgLine.Liner = PkgLine.NewNotify()
	lineNotify.Config("token<-add")
	resp, err := lineNotify.Trigger(string(jPrint))
	if err != nil {
		log.Println(err)
		return
	}
	w.Write([]byte(resp))
}
func loggerHandle(w http.ResponseWriter, r *http.Request) {

}

func main() {
	server := PkgApi.New()
	server.Config("0.0.0.0", 8001)
	server.AddHandle("/notify", alertsHandle)
	server.Run()

}
