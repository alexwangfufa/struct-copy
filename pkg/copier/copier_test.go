package copier

import (
	v1 "github.com/alexwangfufa/struct-copy/example/api/material-group/v1"
	"github.com/alexwangfufa/struct-copy/example/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"testing"
	"time"
)

func Test_Copier(t *testing.T) {

	objectID, err := primitive.ObjectIDFromHex("5dbba1e31fd96208db5a00a1")
	if err != nil {
		panic(err)
	}

	stringValue := &wrapperspb.StringValue{Value: "5dbba1e31fd96208db5a00a1"}
	int64Value := &wrapperspb.Int64Value{Value: 66}
	uint64Value := &wrapperspb.UInt64Value{Value: 66}
	int32Value := &wrapperspb.Int32Value{Value: 66}
	uint32Value := &wrapperspb.UInt32Value{Value: 66}
	doubleValue := &wrapperspb.DoubleValue{Value: 5.5}
	floatValue := &wrapperspb.FloatValue{Value: 5.5}
	boolValue := &wrapperspb.BoolValue{Value: true}
	timestamppb := timestamppb.New(time.Now())
	timestamp := time.Now()

	testCases := []struct {
		Name       string
		Req        []interface{}
		ErrorOccur bool
	}{
		{"copyObjectID2String", []interface{}{&v1.MaterialGroupModel{}, &domain.MaterialGroup{Id: &objectID}}, false},
		{"copyString2ObjectID", []interface{}{&domain.MaterialGroup{}, &v1.MaterialGroupModel{Id: "5dbba1e31fd96208db5a00a1"}}, false},
		{"copyStringValue2ObjectID", []interface{}{&domain.MaterialGroup{}, &v1.SaveMaterialGroupRequest{Id: stringValue}}, false},
		{"copyObjectID2StringValue", []interface{}{&v1.SaveMaterialGroupRequest{}, &domain.MaterialGroup{Id: &objectID}}, false},
		{"copyString2String", []interface{}{&domain.MaterialGroup{}, &v1.MaterialGroupModel{Name: "string2string"}}, false},
		{"copyString2StringValue", []interface{}{&v1.SaveMaterialGroupRequest{}, &domain.MaterialGroup{UserId: "stringValue"}}, false},
		{"copyStringValue2String", []interface{}{&domain.MaterialGroup{}, &v1.SaveMaterialGroupRequest{UserId: stringValue}}, false},
		{"copyInt64Value2Int64", []interface{}{&domain.MaterialGroup{}, &v1.SaveMaterialGroupRequest{Order: int64Value}}, false},
		{"copyInt642Int64Value", []interface{}{&v1.SaveMaterialGroupRequest{}, &domain.MaterialGroup{Order: 66}}, false},
		{"copyUInt64Value2UInt64", []interface{}{&domain.MaterialGroup{}, &v1.SaveMaterialGroupRequest{Ut64: uint64Value}}, false},
		{"copyUInt642UInt64Value", []interface{}{&v1.SaveMaterialGroupRequest{}, &domain.MaterialGroup{Ut64: 66}}, false},
		{"copyInt32Value2Int32", []interface{}{&domain.MaterialGroup{}, &v1.SaveMaterialGroupRequest{It: int32Value}}, false},
		{"copyInt322UInt32Value", []interface{}{&v1.SaveMaterialGroupRequest{}, &domain.MaterialGroup{It: 66}}, false},
		{"copyUInt32Value2UInt32", []interface{}{&domain.MaterialGroup{}, &v1.SaveMaterialGroupRequest{Ut32: uint32Value}}, false},
		{"copyUInt322UInt32Value", []interface{}{&v1.SaveMaterialGroupRequest{}, &domain.MaterialGroup{Ut32: 66}}, false},
		{"copyFloat322FloatValue", []interface{}{&v1.SaveMaterialGroupRequest{}, &domain.MaterialGroup{Point: 5.5}}, false},
		{"copyFloatValue2Float32", []interface{}{&domain.MaterialGroup{}, &v1.SaveMaterialGroupRequest{Point: floatValue}}, false},
		{"copyFloat642DoubleValue", []interface{}{&v1.SaveMaterialGroupRequest{}, &domain.MaterialGroup{StoryPoint: 5.5}}, false},
		{"copyDoubleValue2Float64", []interface{}{&domain.MaterialGroup{}, &v1.SaveMaterialGroupRequest{StoryPoint: doubleValue}}, false},
		{"copyBool2BoolValue", []interface{}{&v1.SaveMaterialGroupRequest{}, &domain.MaterialGroup{IsValid: true}}, false},
		{"copyBoolValue2Bool", []interface{}{&domain.MaterialGroup{}, &v1.SaveMaterialGroupRequest{IsValid: boolValue}}, false},
		{"copyTimestamp2Timestamppb", []interface{}{&v1.SaveMaterialGroupRequest{}, &domain.MaterialGroup{UpdateTime: timestamp}}, false},
		{"copyTimestamppb2Timestamp", []interface{}{&domain.MaterialGroup{}, &v1.SaveMaterialGroupRequest{UpdateTime: timestamppb}}, false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			err := Copy(testCase.Req[0], testCases[1])
			if testCase.ErrorOccur {
				if err == nil {
					zap.L().Error("expect occur exception, but not found")
					t.FailNow()
				}
			} else {
				if err != nil {
					zap.L().Error("occur error:", zap.Error(err))
					t.FailNow()
				}
			}
		})
	}

}
