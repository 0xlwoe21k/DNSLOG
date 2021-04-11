package httpserver

import (
	"encoding/json"
	_"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"strings"
)

func DelAll(w http.ResponseWriter, r *http.Request)  {
	var admin string
	var err error

	fmt.Println("[*] receive http request /v1/dnslog/DelAll")
	_redis, err := redis.Dial("tcp", "127.0.0.1:6379")

	token, err := redis.String(_redis.Do("GET", "token"))
	fmt.Println("token:",token)
	if err != nil {
		fmt.Println("DelAll get token error", err)
		w.Write([]byte("{\"code\":\"302\",\"msg\":\"no login!\"}"))
		return
	}

	apiToken, err := redis.String(_redis.Do("GET", "APITOKEN"))
	fmt.Println("apiToken:",apiToken)
	if err != nil {
		fmt.Println("DelAll get apiToken error", err)
		w.Write([]byte("{\"code\":\"302\",\"msg\":\"no login!\"}"))
		return
	}
	//如果admin也在也要一并存起来
	flag ,err := redis.Int(_redis.Do("EXISTS", "admin"))
	fmt.Println("admin:",admin)
	if err != nil{
		fmt.Println("[!] DelAll EXISTS admin error", err)
		return
	}
	if flag == 1 {
		admin, err = redis.String(_redis.Do("GET", "admin"))
		if err != nil{
			fmt.Println("[!]DelAll get admin error", err)
			return
		}
	}

		_, err = _redis.Do("flushall")
	if err != nil {
		fmt.Println("[*] delete error", err)
		w.Write([]byte("{\"code\":\"100\",\"msg\":\"error!\"}"))
		return
	}else {
		if flag == 1{
			//删除所有后把admin写进去
			_, err = _redis.Do("SET", "admin",admin)
			if err != nil {
				fmt.Println("DelAll set admin error", err)
				return
			}
		}
		//删除所有后把token写进去
		_, err = _redis.Do("SET", "token",token)
		if err != nil {
			fmt.Println("DelAll set error", err)
			return
		}

		//删除所有后把token写进去
		_, err = _redis.Do("SET", "APITOKEN",apiToken)
		if err != nil {
			fmt.Println("DelAll set error", err)
			return
		}
		w.Write([]byte("{\"code\":\"200\",\"msg\":\"delete all keys success!\"}"))
	}
}


func GetAllDnslog(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("[*] receive http request /v1/dnslog/GetAllDnslog")

	var clientToken string

	for name, values := range r.Header {
		for _, value := range values {
			if strings.Contains(name,"Token"){
				clientToken = value
			}
		}
	}

	_redis, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("[GetAllDnslog] Connect to redis error", err)
		w.Write([]byte("{\"code\":\"100\",\"msg\":\"error!\"}"))
		return
	}

	token, err := redis.String(_redis.Do("GET", "token"))
	if err != nil {
		fmt.Println("[GetAllDnslog] get token error", err)
		w.Write([]byte("{\"code\":\"302\",\"msg\":\"no login!\"}"))
		return
	}else{

		if clientToken != token {
			w.Write([]byte("{\"code\":\"302\",\"location\":\"\"}"))
			return
		}
	}

	keys, err := redis.Strings(_redis.Do("KEYS", "*"))
	if err != nil {
		fmt.Println("[GetAllDnslog] KEYS error", err)
		w.Write([]byte("{\"code\":\"100\",\"msg\":\"error!\"}"))
		return
	}

	exitsDataFlag := false

	dnslogdata :="["
	for _, key := range keys {
		var tmp string
		value ,err:= redis.String(_redis.Do("GET", key))
		if err != nil{
			fmt.Println("[*] [GetAllDnslog] read key error from redis:",err)
		}
		if value != "" && key !="token"&& key!="admin" &&key!="APITOKEN"{
			value = strings.Replace(value,"\"","\\\"",-1)
			tmp = fmt.Sprintf("{\"key\":\"%s\",\"value\":\"%s\"}",key,value)
			tmp += ","
			dnslogdata += tmp
			exitsDataFlag = true
		}
	}
	dnslogdata = dnslogdata[:len(dnslogdata)-1]
	dnslogdata +="]"

	if exitsDataFlag{
		w.Write([]byte(dnslogdata))
	}else {
		w.Write([]byte("{\"code\":\"100\",\"msg\":\"无数据\"}"))
	}

	defer _redis.Close()
}

func ModifyAPIToken(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("[*] receive http request /v1/dnslog/ModifyAPIToken")
	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)

	_APITOKEN:=params["APITOKEN"]

	var errorFlag bool = false
	_redis, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		errorFlag = true
	}
	fmt.Println(_APITOKEN)

	_, err = _redis.Do("SET", "APITOKEN",_APITOKEN)
	if err != nil {
		fmt.Println("DelAll set error", err)
		errorFlag = true
	}

	//如果有错误
	if errorFlag {
		w.Write([]byte("{\"code\":\"100\",\"msg\":\"error\"}"))

	}else {
		w.Write([]byte("{\"code\":\"200\",\"msg\":\"修改成功！\"}"))

	}
}

func DnsLogApi(w http.ResponseWriter, r *http.Request)  {

	fmt.Println("[*] receive http request /v1/dnslog")
	//_token := "qweqeqeqeqweqwe12313123213"



	_redis, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}

	APITOKEN, err := redis.String(_redis.Do("GET", "APITOKEN"))
	if err != nil {
		fmt.Println("[*] [GetAllDnslog] read key error from redis:", err)
	}

	apiToken := r.FormValue("token")
	rkey := r.FormValue("key")

	if apiToken == APITOKEN {
		keys, err := redis.Strings(_redis.Do("KEYS", "*"))
		if err != nil {
			fmt.Println("[GetAllDnslog] KEYS error", err)
			w.Write([]byte("{\"code\":\"100\",\"msg\":\"error!\"}"))
			return
		}

		for _, key := range keys {
			if strings.Contains(key, rkey){
				value, err := redis.String(_redis.Do("GET", key))
				if err != nil {
					fmt.Println("[*] [GetAllDnslog] read key error from redis:", err)
				}
				if value != ""{
					w.Write([]byte(value))
				}else{
					w.Write([]byte(""))
				}
			}
		}
	}else {
		w.Write([]byte(""))
	}

	//if token == _token  {
	//	value ,err:= redis.String(_redis.Do("GET", key))
	//	if err != nil{
	//		fmt.Println("[*] read key error from redis:",err)
	//	}
	//
	//	if value != ""{
	//		w.Write([]byte(value))
	//	}else{
	//		w.Write([]byte(""))
	//	}
	//}else {
	//	w.Write([]byte("key error!"))
	//}
	defer _redis.Close()
}