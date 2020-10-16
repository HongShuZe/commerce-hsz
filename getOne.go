package main

import (
	"sync"
	"fmt"
	"net/http"
	"log"
)

// 已秒杀数量
var sum int64

// 预存商品数量
var productNum int64 = 1000000

// 互斥锁
var mutex sync.Mutex

// 计数(访问数量)
var count int64

// 获取秒杀商品
func GetOneProduct() bool {
	// 加锁
	mutex.Lock()
	defer mutex.Unlock()
	count += 1
	// 判断数据是否超限, 每一百个才有一个秒杀成功
	//if count%100 == 0 {
		if sum < productNum {
			sum += 1
			fmt.Println(sum)
			return true
		}
	//}
	return false
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	if GetOneProduct() {
		w.Write([]byte("true"))
		return
	}
	w.Write([]byte("false"))
	return
}

func main() {
	log.Println("getOne启动")
	http.HandleFunc("/getOne", GetProduct)
	err := http.ListenAndServe(":8084", nil)
	if err != nil {
		log.Fatal("Error:", err)
	}
}
