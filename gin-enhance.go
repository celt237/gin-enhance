package gin_enhance

import (
	"context"
	"fmt"
	"github.com/celt237/gin-enhance/internal"
	"github.com/gin-gonic/gin"
	"reflect"
)

type ErrorWithCode interface {
	error
	Code() int
}

type ApiHandler interface {
	// WrapContext 从gin.Context中获取context.Context
	WrapContext(c *gin.Context) context.Context

	// Success 成功返回
	// c *gin.Context 上下文
	// produceType string 返回类型
	// data interface{} 返回数据
	Success(c *gin.Context, produceType string, data interface{})

	// CodeError 失败返回
	// c *gin.Context 上下文
	// produceType string 返回类型
	// data interface{} 返回数据
	// code int 错误码
	// err error 错误
	CodeError(c *gin.Context, produceType string, data interface{}, code int, err error)

	// Error 失败返回
	// c *gin.Context 上下文
	// produceType string 返回类型
	// data interface{} 返回数据
	// err error 错误
	Error(c *gin.Context, produceType string, data interface{}, err error)

	// HandleCustomerAnnotation 处理自定义注解
	// c *gin.Context 上下文
	// annotation string 注解名
	// opt ...string 参数
	HandleCustomerAnnotation(c *gin.Context, annotation string, opt ...string) error
}

// GetParamFromContext 从gin.Context中获取参数
// c *gin.Context 上下文
// paramName string 参数名
// dataType string 数据类型
// paramType string 参数类型
// ptr bool 是否指针
// required bool 是否必须
func GetParamFromContext[T any](c *gin.Context, paramName string, dataType string, paramType string, ptr bool, required bool) (value T, err error) {
	value, err = getDefaultValue[T]()
	var v1 interface{}
	if paramType == internal.ParamTypePath {
		str := c.Param(paramName)
		if str == "" {
			err = fmt.Errorf("param %s can not be empty", paramName)
			return
		}
		v1, err = internal.TypeConvert(str, dataType, ptr)
		if err != nil {
			return
		}
		if v1 != nil {
			value = v1.(T)
		}
	} else if paramType == internal.ParamTypeQuery {
		str := c.Query(paramName)
		if str == "" && required {
			err = fmt.Errorf("param %s can not be empty", paramName)
			return
		}
		v1, err = internal.TypeConvert(str, dataType, ptr)
		if err != nil {
			return
		}
		if v1 != nil {
			value = v1.(T)
		}
	} else if paramType == internal.ParamTypeHeader {
		str := c.GetHeader(paramName)
		if str == "" && required {
			err = fmt.Errorf("param %s can not be empty", paramName)
			return
		}
		v1, err = internal.TypeConvert(str, dataType, ptr)
		if err != nil {
			return
		}
		if v1 != nil {
			value = v1.(T)
		}
	} else if paramType == internal.ParamTypeBody {
		err = c.ShouldBindJSON(&value)
		if err != nil {
			return
		}
	}
	return
}

func getDefaultValue[T any]() (result T, err error) {
	var tempValue interface{}
	t := reflect.TypeOf((*T)(nil)).Elem()
	switch t.Kind() {
	case reflect.Struct, reflect.Map:
		tempValue = reflect.Zero(t).Interface()
	case reflect.Slice, reflect.Array:
		tempValue = reflect.MakeSlice(t, 0, 0).Interface()
	case reflect.Ptr:
		subT := t.Elem()
		switch subT.Kind() {
		case reflect.Struct, reflect.Map:
			elem := reflect.New(t.Elem())
			elem.Elem().Set(reflect.Zero(t.Elem()))
			tempValue = elem.Interface()
		case reflect.Slice, reflect.Array:
			elem := reflect.New(t.Elem())
			elem.Elem().Set(reflect.MakeSlice(t.Elem(), 0, 0))
			tempValue = elem.Interface()
		default:
			break
		}
	default:
		break
	}
	if tempValue != nil {
		result = tempValue.(T)
	}
	return result, nil
}
