package nekoutils

import (
	"context"
	"reflect"
	"unsafe"
)

func In(haystack interface{}, needle interface{}) bool {
	sVal := reflect.ValueOf(haystack)
	kind := sVal.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < sVal.Len(); i++ {
			if sVal.Index(i).Interface() == needle {
				return true
			}
		}

		return false
	}

	return false
}

func getCtxKeys(ctx interface{}) []interface{} {
	var keys = make([]interface{}, 0)

	if a, ok := ctx.(context.Context); ok && ctx != nil {
		ctx = a
	} else {
		return keys
	}

	contextValues := reflect.ValueOf(ctx).Elem()
	contextKeys := reflect.TypeOf(ctx).Elem()

	if contextKeys.Kind() == reflect.Struct {
		for i := 0; i < contextValues.NumField(); i++ {
			reflectValue := contextValues.Field(i)
			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr()))
			reflectValueElem := reflectValue.Elem()
			reflectField := contextKeys.Field(i)

			if reflectField.Name == "key" {
				keys = append(keys, reflectValueElem.Interface())
			} else if reflectValueElem.Kind() == reflect.Struct {
				// timerCtx have a "cancelCtx" struct
				keys = append(keys, getCtxKeys(reflectValue.Interface())...)
			} else if reflectField.Name == "Context" {
				keys = append(keys, getCtxKeys(reflectValueElem.Interface())...)
			}
		}
	}

	return keys
}

func GetCtxKeyValues(ctx context.Context) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	for _, k := range getCtxKeys(ctx) {
		m[k] = ctx.Value(k)
	}
	return m
}

// Return a new context with key value
func CopyCtx(ctx context.Context) context.Context {
	kv := GetCtxKeyValues(ctx)
	ctx2 := context.Background()
	for k, v := range kv {
		ctx2 = context.WithValue(ctx2, k, v)
	}
	return ctx2
}

func CorePtrFromContext(ctx context.Context) uintptr {
	for _, k := range getCtxKeys(ctx) {
		if reflect.TypeOf(k).Name() == "v2rayKeyType" {
			core := ctx.Value(k)
			return reflect.ValueOf(core).Pointer()
		}
	}
	return 0
}
