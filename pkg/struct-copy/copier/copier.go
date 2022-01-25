package copier

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"reflect"
	"time"
)

var (
	stringValue = reflect.TypeOf(&wrapperspb.StringValue{}).Kind()
	int64Value  = reflect.TypeOf(&wrapperspb.Int64Value{}).Kind()
	int32Value  = reflect.TypeOf(&wrapperspb.Int32Value{}).Kind()
	doubleValue = reflect.TypeOf(&wrapperspb.DoubleValue{}).Kind()
	floatValue  = reflect.TypeOf(&wrapperspb.FloatValue{}).Kind()
	uint32Value = reflect.TypeOf(&wrapperspb.UInt32Value{}).Kind()
	uint64Value = reflect.TypeOf(&wrapperspb.UInt64Value{}).Kind()
	boolValue   = reflect.TypeOf(&wrapperspb.BoolValue{}).Kind()
	objectID    = reflect.TypeOf(primitive.ObjectID{}).Kind()
	pbTimestamp = reflect.TypeOf(&timestamppb.Timestamp{}).Kind()
	timestamp   = reflect.TypeOf(&time.Time{}).Kind()
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
			continue
		}

		// 如果是指针类型,并且为nil,不处理
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}

		switch fieldValue.Kind() {
		case stringValue:
			if field.Name == "Id" {
				value := fieldValue.Elem().FieldByName("Value").String()
				id, err := primitive.ObjectIDFromHex(value)
				if err != nil {
					return err
				}
				dstValue.Field(i).Set(reflect.ValueOf(&id))
			} else {
				value := fieldValue.Elem().FieldByName("Value").String()
				dstValue.Field(i).Set(reflect.ValueOf(value))
			}
		case int64Value:
			value := fieldValue.FieldByName("Value").Int()
			dstValue.Field(i).Set(reflect.ValueOf(value))
		case int32Value:
			value := fieldValue.FieldByName("Value").Int()
			dstValue.Field(i).Set(reflect.ValueOf(int32(value)))
		case uint64Value:
			value := fieldValue.FieldByName("Value").Int()
			dstValue.Field(i).Set(reflect.ValueOf(uint64(value)))
		case uint32Value:
			value := fieldValue.FieldByName("Value").Int()
			dstValue.Field(i).Set(reflect.ValueOf(uint32(value)))
		case doubleValue, floatValue:
			value := fieldValue.FieldByName("Value").Float()
			dstValue.Field(i).Set(reflect.ValueOf(value))
		case boolValue:
			value := fieldValue.FieldByName("Value").Bool()
			dstValue.Field(i).Set(reflect.ValueOf(value))
		case pbTimestamp:
			dstValue.Field(i).Set(reflect.ValueOf(floatValue))
		case reflect.String:
			value := fieldValue.String()
			if fieldValue.Kind() == stringValue {
				spbString := &wrapperspb.StringValue{
					Value: value,
				}
				dstValue.Field(i).Set(reflect.ValueOf(spbString))
			} else {
				dstValue.Field(i).Set(reflect.ValueOf(value))
			}

		case reflect.Int64:
			value := fieldValue.Int()
			if fieldValue.Kind() == int64Value {
				spbInt64 := &wrapperspb.Int64Value{
					Value: value,
				}
				dstValue.Field(i).Set(reflect.ValueOf(spbInt64))
			} else {
				dstValue.Field(i).Set(reflect.ValueOf(value))
			}
		case reflect.Int32:
			value := fieldValue.Int()
			if fieldValue.Kind() == int32Value {
				spbInt32 := &wrapperspb.Int32Value{
					Value: int32(value),
				}
				dstValue.Field(i).Set(reflect.ValueOf(spbInt32))
			} else {
				dstValue.Field(i).Set(reflect.ValueOf(value))
			}
		case reflect.Uint64:
			value := fieldValue.Int()
			if fieldValue.Kind() == uint64Value {
				spbUint64 := &wrapperspb.UInt64Value{
					Value: uint64(value),
				}
				dstValue.Field(i).Set(reflect.ValueOf(spbUint64))
			} else {
				dstValue.Field(i).Set(reflect.ValueOf(value))
			}
		case reflect.Uint32:
			value := fieldValue.Int()
			if fieldValue.Kind() == uint32Value {
				spbUint32 := &wrapperspb.UInt32Value{
					Value: uint32(value),
				}
				dstValue.Field(i).Set(reflect.ValueOf(spbUint32))
			} else {
				dstValue.Field(i).Set(reflect.ValueOf(value))
			}
		case reflect.Float64, reflect.Float32:
			value := fieldValue.Float()
			if fieldValue.Kind() == doubleValue {
				spbDouble := &wrapperspb.DoubleValue{
					Value: value,
				}
				dstValue.Field(i).Set(reflect.ValueOf(spbDouble))
			} else if fieldValue.Kind() == floatValue {
				spbFloat := &wrapperspb.FloatValue{
					Value: float32(value),
				}
				dstValue.Field(i).Set(reflect.ValueOf(spbFloat))
			} else {
				dstValue.Field(i).Set(reflect.ValueOf(value))
			}
		case timestamp:
			timeValue := fieldValue.Elem().Int()
			newTime := timestamppb.New(time.UnixMilli(timeValue))
			dstValue.Field(i).Set(reflect.ValueOf(newTime))
		default:
			return errors.New("unsupported data type")
		}

	}
	return err
}
