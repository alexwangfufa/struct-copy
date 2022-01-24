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
	//defer func() {
	//	if err := recover(); err != nil {
	//		err = errors.New(fmt.Sprintf("%v", err))
	//	}
	//}()

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
	fmt.Printf("the struct have %v fields\n", fieldNum)

	for i := 0; i < fieldNum; i++ {
		field := dstType.Field(i)
		// 找到src中与dst相同字段的value
		fieldValue := srcValue.FieldByName(field.Name)

		// 无效, 说明src没有这个属性
		if !fieldValue.IsValid() {
			fmt.Println(field.Name)
			continue
		}

		fmt.Printf("field name is %v\n", field.Name)
		if field.Name == "Order" {
			fmt.Printf("value is %v\n", fieldValue)
		}

		//switch field.Name {
		//case "Id":
		//	println("step id")
		//	dstValue.Field(i).Set(reflect.ValueOf(fieldValue))
		//	println("=======")
		//case "OrgId":
		//	println("step orgId")
		//	dstValue.Field(i).Set(dstValue)
		//case "Order":
		//	println("step order")
		//	dstValue.Field(i).SetInt(reflect.ValueOf(fieldValue).Int())
		//default:
		//	println("step default")
		//}

		if field.Name == "Id" {
			println("step id")
			value := fieldValue.Elem().FieldByName("Value").String()
			id, err := primitive.ObjectIDFromHex(value)
			if err != nil {
				return err
			}
			dstValue.Field(i).Set(reflect.ValueOf(id))
			println("=======")
		} else {
			continue
		}

	}
	return err
}

func main() {
	objId, err := primitive.ObjectIDFromHex("5dbba1e31fd96208db5a00a1")
	if err != nil {
		panic(err)
	}

	saveMaterialGroupRequest := &v1.SaveMaterialGroupRequest{Name: "test", Order: 66}
	materialGroup := &domain.MaterialGroup{Id: objId}
	if err := copy(materialGroup, saveMaterialGroupRequest); err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%+v\n", materialGroup)
}
