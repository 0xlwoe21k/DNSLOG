package httpserver

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func checkErr(err error)  {
	fmt.Println("[*] error",err)
}


func ResetPassword(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("[*] receive http request /v1/dnslog/ResetPassword")
	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)

	oldPass:=params["oldPass"]
	newPass:=params["newPass"]
	checkPass:=params["checkPass"]

	if newPass != checkPass {
		w.Write([]byte("{\"code\":\"100\",\"msg\":\"两次密码不一样，傻吊!\"}"))
		return
	}

	_redis, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("[!] 32-ResetPassword Connect to redis error", err)
		return
	}
	flag ,err := redis.Int(_redis.Do("EXISTS", "admin"))
	if err != nil{
		fmt.Println("[!] 37-ResetPassword Connect to redis error", err)
		return
	}

	if flag == 1 {
		value ,err:= redis.String(_redis.Do("GET", "admin"))
		if err != nil {
			fmt.Println("[!] 43-ResetPassword Connect to redis error", err)
			return
		}

		if oldPass == value {
			_ ,err:= _redis.Do("SET", "admin",newPass)
			if err != nil {
				fmt.Println("[!] 45-ResetPassword Connect to redis error", err)
				return
			}
			w.Write([]byte("{\"code\":\"200\",\"msg\":\"密码修改成功!\"}"))
		}
	}else{
		_ ,err:= _redis.Do("SET", "admin",newPass)
		if err != nil {
			fmt.Println("[!] 45-ResetPassword Connect to redis error", err)
			return
			w.Write([]byte("{\"code\":\"200\",\"msg\":\"初次设置密码成功!\"}"))
		}
	}


}

func Login(w http.ResponseWriter, r *http.Request) {

	fmt.Println("[*] receive http request /v1/dnslog/login")
	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)
	var _redis redis.Conn
	var err error
	var isExist int


	username:=params["username"]
	password:=params["password"]


	_redis, err = redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("[Login] Connect to redis error", err)
		return
	}

	defer _redis.Close()

	if username == "admin" {
		isExist ,err = redis.Int(_redis.Do("EXISTS", "admin"))
		if err != nil {
			fmt.Println("[Login] EXISTS token error:", err)
			return
		}

		if isExist == 1 {
			pass ,err:= redis.String(_redis.Do("GET", "admin"))
			if err != nil {
				fmt.Println("get key error", err)
				return
			}
			if password == pass{
				timeTemp := strconv.Itoa( DateUnix())
				token := md5V(username+timeTemp)
				_ ,err:= _redis.Do("SET", "token",token)
				if err != nil {
					fmt.Println("set token error", err)
					return
				}
				_ ,err= _redis.Do("EXPIRE", "token","86400") //1800半小时
				if err != nil {
					fmt.Println("set token EXPIRE error", err)
				}

				w.Write([]byte("{\"code\":\"200\",\"token\":\""+token+"\"}"))
			}else{
				w.Write([]byte("{\"code\":\"100\",\"token\":\"\"}"))
			}
		}else {
			if password == "dnslog"{
				timeTemp := strconv.Itoa( DateUnix())
				token := md5V(username+timeTemp)
				_ ,err:= _redis.Do("SET", "token",token)
				if err != nil {
					fmt.Println("set token2  error", err)
				}
				_ ,err= _redis.Do("EXPIRE", "token","86400")
				if err != nil {
					fmt.Println("set token EXPIRE 2 error", err)
				}
				w.Write([]byte("{\"code\":\"200\",\"token\":\""+token+"\"}"))
			}else{
				w.Write([]byte("{\"code\":\"100\",\"token\":\"\"}"))
			}
		}
	}

}

func CheckLogin(w http.ResponseWriter, r *http.Request)  {
	//var pass string

	var token string
	var _redis redis.Conn
	var isExist int
	var hasToken bool
	var err error
	var clientToken string

	fmt.Println("[*] receive http request /v1/dnslog/CheckLogin")


	for name, values := range r.Header {
		for _, value := range values {
			if strings.Contains(name,"Token"){
				clientToken = value
				hasToken = true
			}
		}

	}

	if !hasToken {
		w.Write([]byte("{\"code\":\"302\",\"location\":\"/login\"}"))
		return
	}

	_redis, err = redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("[Login] Connect to redis error", err)
		return
	}
	//用户登录会带有缓存，其中有token ， 到数据库里查这个toke是否存在 ，存在就直接跳转到
	isExist ,err = redis.Int(_redis.Do("EXISTS", "token"))
	if err != nil {
		fmt.Println("[Login] EXISTS token error:", err)
		return
	}
	if isExist == 1 {
		token ,err = redis.String(_redis.Do("GET", "token"))
		if err != nil {
			fmt.Println("[Login] get token error:", err)
			return
		}
		if clientToken != token {
			w.Write([]byte("{\"code\":\"302\",\"location\":\"\"}"))
			return
		}
	}else {
		w.Write([]byte("{\"code\":\"302\",\"location\":\"\"}"))
		return
	}
}

func DateUnix() int {
	t := time.Now().Local().Unix()
	return int(t)
}

func md5V(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}