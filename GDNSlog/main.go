package main

import (
	"GDNSlog/httpserver"
	Sniffer "GDNSlog/source"
	"fmt"
	"net/http"
	"sync"
)



var wg sync.WaitGroup

func HttpServer()  {

	httpPort:=":443"
	fmt.Println("[*] listen http port:"+httpPort)
	http.HandleFunc("/v1/dnslog", httpserver.DnsLogApi)
	http.HandleFunc("/v1/dnslog/login", httpserver.Login)
	http.HandleFunc("/v1/dnslog/getAllDnslog", httpserver.GetAllDnslog)
	http.HandleFunc("/v1/dnslog/delAll", httpserver.DelAll)
	http.HandleFunc("/v1/dnslog/resetPassword", httpserver.ResetPassword)
	http.HandleFunc("/v1/dnslog/checkLogin", httpserver.CheckLogin)
	http.HandleFunc("/v1/dnslog/ModifyAPIToken", httpserver.ModifyAPIToken)


	fmt.Println("---------------------------------------------------------------------")
	err := http.ListenAndServe(httpPort, nil)
	if err != nil {
		fmt.Println("[*] port alread be used"+httpPort)
		panic(err)
	}

	wg.Done()
}


func main() {

	wg.Add(1)
	go HttpServer()
	wg.Add(1)
	go Sniffer.RunSniffer(wg)

	wg.Wait()



}
