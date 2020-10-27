package main

import (
	"commerce-hsz/common"
	"net/http"
	"fmt"
	"errors"
	"commerce-hsz/encrypt"
	"sync"
	"strconv"
	"io/ioutil"
	"commerce-hsz/rabbitmq"
	"net/url"
	"commerce-hsz/datamodels"
	"encoding/json"
	"log"
	"time"
)

// 设置集群地址, 最好内部IP
var hostArray = []string{"192.168.20.143", "192.168.20.143"}


var localHost = ""
// 数量控制接口服务器内网IP，或者getOne的SLB内网IP
var GetOneIp = "127.0.0.1"
// 47.112.245.134

var GetOnePort = "8084"

var port = "8083"

var hashConsistent *common.Consistent

var rmqValidate *rabbitmq.RabbitMQ

// 用来存放控制信息
type AccessControl struct {
	// 用来存放用户想要的存放的信息
	sourcesArray map[int]time.Time
	sync.RWMutex
}

// 创建全局变量
var accessControl = &AccessControl{
	sourcesArray: make(map[int]time.Time),
}

// 获取定制的数据
func (m *AccessControl)GetNewRecord(uid int) time.Time {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	return m.sourcesArray[uid]
}

// 设置记录
func (m *AccessControl)SetNewRecord(uid int)  {
	m.RWMutex.Lock()
	m.sourcesArray[uid]=time.Now()
	m.RWMutex.Unlock()
}

// 黑名单结构体
type BlackList struct {
	listArray map[int]bool
	sync.RWMutex
}
// 初始化黑名单
var blackList = &BlackList{listArray: make(map[int]bool)}

// 获取黑名单
func (m *BlackList)GetBlackList(uid int) bool {
	m.RLock()
	defer m.RUnlock()
	return m.listArray[uid]
}

// 获取黑名单
func (m *BlackList)SetBlackList(uid int) bool {
	m.Lock()
	defer m.Unlock()
	m.listArray[uid] = true
	return true
}


// 判断服务器所在位置
func (m *AccessControl)GetDistributedRight(req *http.Request) bool {
	// 获取用户id
	uid, err := req.Cookie("uid" )
	if err != nil {
		fmt.Println(11)
		return false
	}

	// 采用一致性hash算法, 根据用户id判断获取具体机器
	hostRequest	, err := hashConsistent.Get(uid.Value)
	fmt.Println("hostRequest:"+hostRequest)
	if err != nil {
		fmt.Println(22)
		return false
	}

	// 判断是否为本机
	if hostRequest == localHost {
		// 执行本机数据读取和校验
		fmt.Println(33)
		return m.GetDataFromMap(uid.Value)
	} else {
		// 不是本机充当代理访问数据返回结果
		fmt.Println(44)
		return GetDataFromOtherMap(hostRequest, req)
	}
}

// 获取本机map, 并且处理业务逻辑, 返回bool
func (m *AccessControl)GetDataFromMap(uid string) bool {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		log.Println(err)
		return false
	}
	// 判断是否被加入黑名单中
	if blackList.GetBlackList(uidInt) {
		log.Println("该用户已被加入黑名单")
		return false
	}

	// 获取记录
	dataRecord := m.GetNewRecord(uidInt)
	// 判断时间是否为零
	if !dataRecord.IsZero() {
		// 业务判断, 是否在指定时间之后(限制抢购时间间隔)
		if dataRecord.Add(time.Duration(20)*time.Second).After(time.Now()) {
			return false
		}
	}
	// 添加记录
	m.SetNewRecord(uidInt)
	return true
}

// 获取其他节点处理结果
func GetDataFromOtherMap(host string, r *http.Request) bool {

	hostUrl := "http://"+host+":"+port+"/checkRight"
	response, body, err := GetCurl(hostUrl, r)
	if err != nil {
		return false
	}

	// 判断状态
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
}

// 模拟请求
func GetCurl(hostUrl string, r *http.Request) (response *http.Response, body []byte, err error) {
	// 获取uid
	uid, err := r.Cookie("uid")
	if err != nil {
		return
	}
	// 获取sign
	sign, err := r.Cookie("sign")
	if err != nil {
		return
	}

	// 模拟接口访问
	client := &http.Client{}
	req, err := http.NewRequest("GET", hostUrl, nil)
	if err != nil {
		return
	}
	// 手动指定, 排查多余cookies
	cookieUid := &http.Cookie{Name: "uid", Value: uid.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: sign.Value, Path: "/"}
	// 添加cookie到模拟的请求中
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	// 获取返回结果
	response, err = client.Do(req)
	//defer response.Body.Close()
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(response.Body)
	/*if err != nil {
		return
	}*/
	return
}

func CheckRight(w http.ResponseWriter, r *http.Request)  {
	right := accessControl.GetDistributedRight(r)
	if !right {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return
}

// 执行正常业务逻辑
func Check(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("执行check!")

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	// todo
	if err != nil || len(queryForm["productID"]) <= 0 {
		w.Write([]byte("false"))
		log.Println(err)
		return
	}

	productString := queryForm["productID"][0]
	// fmt.Println("productString" + productString)

	// 获取用户cookie
	userCookie, err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("false"))
		log.Println(err)
		return
	}

	// 1.分布式权限验证
	right := accessControl.GetDistributedRight(r)
	if right == false {
		log.Println("分布式权限验证err")
		w.Write([]byte("false"))
		return
	}

	// 2.获取数量控制权限, 防止秒杀出现超卖现象
	hostUrl := "http://" + GetOneIp + ":" + GetOnePort + "/getOne"
	// http://172.28.21.91:8084/getOne
	responseValidate, validateBody, err := GetCurl(hostUrl, r)
	if err != nil {
		log.Println(err)
		w.Write([]byte("false"))
		return
	}

	// 判断数量控制接口请求状态
	if responseValidate.StatusCode == 200 {
		if string(validateBody) == "true" {
			// 整合下单
			// 1.获取商品ID
			productID, err := strconv.ParseInt(productString, 10, 64)
			if err != nil {
				log.Println(err)
				w.Write([]byte("false"))
				return
			}
			// 2.获取用户ID
			userID, err := strconv.ParseInt(userCookie.Value, 10, 64)
			if err != nil {
				log.Println(err)
				w.Write([]byte("false"))
				return
			}
			// 3.创建消息体
			message := datamodels.NewMessage(userID, productID)
			// 类型转换
			byteMessage, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				w.Write([]byte("false"))
				return
			}
			fmt.Println("发送到消息队列")
			// 生产消息
			err = rmqValidate.PublishSimple(string(byteMessage))
			if err != nil {
				log.Println(err)
				w.Write([]byte("false"))
				return
			}
			w.Write([]byte("true"))
			return
		}
	}
	log.Println("getOne failed")
	w.Write([]byte("false"))
	return
}
//////////////////////////////////////////////////////////////////////////////////
// 统一验证拦截起
func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("执行验证")
	// 添加基于cookie的权限验证
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

// 身份校验函数
func CheckUserInfo(r *http.Request) error {
	//
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("uid 获取失败")
	}
	//
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("sign 获取失败")
	}

	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("加密串已被篡改")
	}

	if CheckInfo(uidCookie.Value, string(signByte)) {
		return nil
	}

	return errors.New("身份校验失败")
}

// 自定义逻辑判断
func CheckInfo(checkStr string, signStr string) bool {
	if checkStr == signStr {
		return true
	}
	return false
}
//////////////////////////////////////////////////////////////////////
func main()  {
	// 负载均衡器设置
	// 采用一致性hash算法
	hashConsistent = common.NewConsistent()
	// 采用一致性hash算法, 添加节点
	for _, val := range hostArray{
		hashConsistent.Add(val)
	}

	// 获取ip
	localIp, err := common.GetIntranceIp()
	if err != nil {
		fmt.Println(err)
	}
	localHost = localIp
	fmt.Println(localHost)
	// 初始化rmq
	rmqValidate = rabbitmq.NewRabbitMQSimple("order_product")
	defer rmqValidate.Destory()

	// 设置静态文件目录
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("./fronted/web/htmlProductShow"))))
	// 设置资源目录
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./fronted/web/public"))))
	// 1.过滤器
	filter := common.NewFilter()
	// 注册拦截
	filter.RegisterFilterUri("/check", Auth)
	filter.RegisterFilterUri("/checkRight", Auth)
	// 2.启动服务
	http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("/checkRight", filter.Handle(CheckRight))

	http.ListenAndServe(":8083", nil)
}
