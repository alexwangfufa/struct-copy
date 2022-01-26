package copier

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"reflect"
	"runtime"
	"time"
)

var (
	stringValue = reflect.TypeOf(wrapperspb.StringValue{}).Kind()
	int64Value  = reflect.TypeOf(wrapperspb.Int64Value{}).Kind()
	int32Value  = reflect.TypeOf(wrapperspb.Int32Value{}).Kind()
	doubleValue = reflect.TypeOf(wrapperspb.DoubleValue{}).Kind()
	floatValue  = reflect.TypeOf(wrapperspb.FloatValue{}).Kind()
	uint32Value = reflect.TypeOf(wrapperspb.UInt32Value{}).Kind()
	uint64Value = reflect.TypeOf(wrapperspb.UInt64Value{}).Kind()
	boolValue   = reflect.TypeOf(wrapperspb.BoolValue{}).Kind()
	objectID    = reflect.TypeOf(primitive.ObjectID{}).Kind()
	pbTimestamp = reflect.TypeOf(timestamppb.Timestamp{}).Kind()
	timestamp   = reflect.TypeOf(time.Time{}).Kind()
)

func Copy(dst, src interface{}) (err error) {
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		err := recover()
		if err != nil {
			switch err.(type) {
			case runtime.Error: // 运行时错误
				err = errors.Cause(err.(runtime.Error))
			default: // 非运行时错误

				err = errors.Cause(err.(error))
			}
		}
	}()
	return copy(dst, src)
}

func copy(dst, src interface{}) error {

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

	for i := 0; i < fieldNum; i++ {
		field := dstType.Field(i)
		// 找到src中与dst相同字段的value
		fieldValue := srcValue.FieldByName(field.Name)

		dstFieldValue := dstValue.Field(i)

		// 无效, 说明src没有这个属性
		if !fieldValue.IsValid() {
			continue
		}

		// 如果是指针类型,并且为nil,不处理
		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				continue
			}
			fieldValue = fieldValue.Elem()
		}

		switch fieldValue.Kind() {
		// StringValue to string, if srcField name is id, StringValue to primitive.ObjectID
		case stringValue:
			value := fieldValue.FieldByName("Value").String()

			if dstFieldValue.Kind() == reflect.String {
				dstFieldValue.SetString(value)
			} else {
				id, err := primitive.ObjectIDFromHex(value)
				if err != nil {
					return err
				}
				dstFieldValue.Set(reflect.ValueOf(&id))
			}

		// Int64Value to int64
		case int64Value:
			value := fieldValue.FieldByName("Value").Int()
			dstFieldValue.Set(reflect.ValueOf(value))

		// Int32Value to int32
		case int32Value:
			value := fieldValue.FieldByName("Value").Int()
			dstFieldValue.Set(reflect.ValueOf(int32(value)))

		// UInt64Value to uint64
		case uint64Value:
			value := fieldValue.FieldByName("Value").Int()
			dstFieldValue.Set(reflect.ValueOf(uint64(value)))

		// UInt32Value to uint32
		case uint32Value:
			value := fieldValue.FieldByName("Value").Int()
			dstFieldValue.Set(reflect.ValueOf(uint32(value)))

		// DoubleValue and FloatValue to float
		case doubleValue, floatValue:
			value := fieldValue.FieldByName("Value").Float()
			dstFieldValue.Set(reflect.ValueOf(value))

		// BoolValue to bool
		case boolValue:
			value := fieldValue.FieldByName("Value").Bool()
			dstFieldValue.Set(reflect.ValueOf(value))

		// timestamppb to time.Time
		case pbTimestamp:
			dstFieldValue.Set(reflect.ValueOf(fieldValue))

		// primitive.ObjectID to StringValue or string
		case objectID:
			value := fieldValue.MethodByName("Hex").Call(nil)
			if dstFieldValue.Kind() == reflect.String {
				dstFieldValue.Set(reflect.ValueOf(value[0].String()))
			} else {
				id := &wrapperspb.StringValue{Value: value[0].String()}
				dstFieldValue.Set(reflect.ValueOf(id))
			}

		// string to StringValue or String
		case reflect.String:
			value := fieldValue.String()
			// if dstField type is StringValue
			if dstFieldValue.Kind() == reflect.String {
				dstFieldValue.SetString(value)
			} else if dstFieldValue.Kind() == objectID {
				id, err := primitive.ObjectIDFromHex(value)
				if err != nil {
					return err
				}
				dstFieldValue.Set(reflect.ValueOf(id))
			} else {
				spbString := &wrapperspb.StringValue{
					Value: value,
				}
				dstFieldValue.Set(reflect.ValueOf(spbString))
			}

		// int64 to int64 or Int64Value
		case reflect.Int64:
			value := fieldValue.Int()
			// if dstFieldValue type is Int64Value
			if dstFieldValue.Kind() == int64Value {
				spbInt64 := &wrapperspb.Int64Value{
					Value: value,
				}
				dstFieldValue.Set(reflect.ValueOf(spbInt64))
			} else {
				dstFieldValue.Set(reflect.ValueOf(value))
			}

		// int32 to int32 or Int32
		case reflect.Int32:
			value := fieldValue.Int()
			// if dstFieldValue type is Int32Value
			if dstFieldValue.Kind() == int32Value {
				spbInt32 := &wrapperspb.Int32Value{
					Value: int32(value),
				}
				dstFieldValue.Set(reflect.ValueOf(spbInt32))
			} else {
				dstFieldValue.Set(reflect.ValueOf(value))
			}

		// uint64 to uint64 or UInt64Value
		case reflect.Uint64:
			value := fieldValue.Int()
			// if dstFieldValue type is UInt64Value
			if dstFieldValue.Kind() == uint64Value {
				spbUint64 := &wrapperspb.UInt64Value{
					Value: uint64(value),
				}
				dstFieldValue.Set(reflect.ValueOf(spbUint64))
			} else {
				dstFieldValue.Set(reflect.ValueOf(value))
			}

		// uint32 to uint32 or UInt32Value
		case reflect.Uint32:
			value := fieldValue.Int()
			// if dstFieldValue type is UInt32Value
			if dstFieldValue.Kind() == uint32Value {
				spbUint32 := &wrapperspb.UInt32Value{
					Value: uint32(value),
				}
				dstFieldValue.Set(reflect.ValueOf(spbUint32))
			} else {
				dstFieldValue.Set(reflect.ValueOf(value))
			}

		// float64 or float32 to DoubleValue or FloatValue
		case reflect.Float64, reflect.Float32:
			value := fieldValue.Float()
			if dstFieldValue.Kind() == doubleValue {
				spbDouble := &wrapperspb.DoubleValue{
					Value: value,
				}
				dstFieldValue.Set(reflect.ValueOf(spbDouble))
			} else if dstFieldValue.Kind() == floatValue {
				spbFloat := &wrapperspb.FloatValue{
					Value: float32(value),
				}
				dstFieldValue.Set(reflect.ValueOf(spbFloat))
			} else {
				dstFieldValue.Set(reflect.ValueOf(value))
			}

		// timestamp to timestamppb
		case timestamp:
			timeValue := fieldValue.Int()
			newTime := timestamppb.New(time.UnixMilli(timeValue))
			dstFieldValue.Set(reflect.ValueOf(newTime))

		// don't have one of these types, should continue
		default:
			continue
		}

	}

	return nil
}
