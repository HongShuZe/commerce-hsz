package common

import "net/http"

// 声明一个新的数据类型(函数类型)
type FilterHandle func(w http.ResponseWriter, r *http.Request) error

// 拦截器结构体
type Filter struct {
	// 用来存储需要拦截的uri
	filterMap map[string]FilterHandle
}

// Filter初始化函数
func NewFilter() *Filter {
	return &Filter{filterMap: make(map[string]FilterHandle)}
}

// 注册拦截器
func (f *Filter) RegisterFilterUri(uri string, handler FilterHandle) {
	f.filterMap[uri] = handler
}

// 根据uri获取对应的handle
func (f *Filter) GetFilterHandle(uri string) FilterHandle {
	return f.filterMap[uri]
}

// 声明新的函数类型
type WebHandle func(w http.ResponseWriter, r *http.Request)

// 执行拦截器, 返回函数类型
func (f *Filter) Handle(webHandle WebHandle) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for path, handle := range f.filterMap {
			if path == r.RequestURI {
				// 执行拦截业务逻辑
				err := handle(w, r)
				if err != nil {
					w.Write([]byte(err.Error()))
					return
				}
				// 跳出循环
				break
			}
		}
		// 执行正常注册函数
		webHandle(w, r)
	}
}

