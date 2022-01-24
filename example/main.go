package main

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"reflect"
	v1 "struct-copy/example/api/material-group/v1"
	"struct-copy/example/domain"
)

var (
	stringValue = reflect.TypeOf(&wrapperspb.StringValue{}).Kind()
	int64Value  = reflect.TypeOf(&wrapperspb.Int64Value{}).Kind()
	boolValue   = reflect.TypeOf(&wrapperspb.BoolValue{}).Kind()
	bytesValue  = reflect.TypeOf(&wrapperspb.BytesValue{}).Kind()
	objectID    = reflect.TypeOf(primitive.ObjectID{}).Kind()
)

func copy(dst, src interface{}) (err error) {
	// 防止意外panic
	defer func() {
		if err := recover(); err != nil {
			err = errors.New(fmt.Sprintf("%v", err))
		}
	}()
	dstType, dstValue := reflect.TypeOf(dst), reflect.ValueOf(dst)
	srcType, srcValue := reflect.TypeOf(src), reflect.ValueOf(src)

	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
		return errors.New("dest type should be a struct pointer")
	}

	if srcType.Kind() == reflect.Ptr {
		srcType, srcValue = srcType.Elem(), srcValue.Elem()
	}

	if srcType.Kind() != reflect.Struct {
		return errors.New("src type should be a struct pointer")
	}

	// 取具体的值
	dstType, dstValue = dstType.Elem(), dstValue.Elem()

	fieldNum := dstType.NumField()

	for i := 3; i < fieldNum; i++ {
		field := dstType.Field(i)
		// 找到src中与dst相同字段的value
		fieldValue := srcValue.FieldByName(field.Name)

		// 无效, 说明src没有这个属性
		if !fieldValue.IsValid() {
			continue
		}

		switch fieldValue.Kind() {
		case stringValue:
			dstValue.Field(i).SetString(fieldValue.Elem().FieldByName("Value").String())
		case int64Value:
			dstValue.Field(i).SetInt(fieldValue.Elem().FieldByName("Value").Int())
		case boolValue:
			dstValue.Field(i).SetBool(fieldValue.Elem().FieldByName("Value").Bool())
		case bytesValue:
			dstValue.Field(i).SetBytes(fieldValue.Elem().FieldByName("Value").Bytes())
		// 其它的基础类型
		default:
			dstValue.Field(i).Set(fieldValue)
		}
	}
	return err
}

func antiCopy(dst, src interface{}) (err error) {
	// 防止意外panic
	defer func() {
		if err := recover(); err != nil {
			err = errors.New(fmt.Sprintf("%v", err))
		}
	}()
	dstType, dstValue := reflect.TypeOf(dst), reflect.ValueOf(dst)
	srcType, srcValue := reflect.TypeOf(src), reflect.ValueOf(src)

	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
		return errors.New("dest type should be a struct pointer")
	}

	if srcType.Kind() == reflect.Ptr {
		srcType, srcValue = srcType.Elem(), srcValue.Elem()
	}

	if srcType.Kind() != reflect.Struct {
		return errors.New("src type should be a struct pointer")
	}

	// 取具体的值
	dstType, dstValue = dstType.Elem(), dstValue.Elem()

	fieldNum := dstType.NumField()

	for i := 3; i < fieldNum; i++ {
		field := dstType.Field(i)
		// 找到src中与dst相同字段的value
		fieldValue := srcValue.FieldByName(field.Name)

		// 无效, 说明src没有这个属性
		if !fieldValue.IsValid() {
			continue
		}

		fmt.Printf("%v", dstValue)
	}
	return err
}

func main() {
	var id *wrapperspb.StringValue
	id = &wrapperspb.StringValue{
		Value: "225jk25g523j5gjh45g",
	}
	objId, err := primitive.ObjectIDFromHex("225jk25g523j5gjh45g")
	if err != nil {
		panic(err)
	}

	saveMaterialGroupRequest := &v1.SaveMaterialGroupRequest{}
	materialGroup := &domain.MaterialGroup{Id: objId, Name: "test"}
	if err := copy(materialGroup, saveMaterialGroupRequest); err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%+v\n", saveMaterialGroupRequest)
	fmt.Printf("%+v\n", id)
}
