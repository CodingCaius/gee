package session

import (
	"geeorm/log"
	"reflect"
)

// Hooks constants
const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

// CallMethod 调用注册的钩子函数
func (s *Session) CallMethod(method string, value interface{}) {
	// method：表示要调用的钩子函数的名称, value：表示要调用钩子函数的对象
	// 钩子函数的定义和注册在其他地方完成，而这个代码片段负责根据钩子名称调用相应的函数

	// 使用反射获取了与当前 Session 对象关联的表格模型（Model）的方法集
	// 然后用 MethodByName 获取钩子
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	if value != nil {
		// 如果value不为空，就获取value类型绑定的钩子函数
		fm = reflect.ValueOf(value).MethodByName(method)
	}
	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {
		if v := fm.Call(param); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
}