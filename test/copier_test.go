package test

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	v1 "struct-copy/example/api/material-group/v1"
	"struct-copy/example/domain"
	"struct-copy/pkg/struct-copy/copier"
	"testing"
	"time"
)

func Test_Copier(t *testing.T) {

	objectID, err := primitive.ObjectIDFromHex("5dbba1e31fd96208db5a00a1")

	if err != nil {
		panic(err)
	}

	inputMaterial := &domain.MaterialGroup{
		Id:         &objectID,
		OrgId:      "5dbba1e31fd96208db5a00a1",
		Name:       "test",
		Type:       "org",
		UserId:     "wxbba1e31fd96208db5a00a1",
		Scope:      "group",
		StoryPoint: 3.5,
		Point:      5.5,
		IsValid:    true,
		It:         55,
		Ut32:       1,
		Ut64:       2,
		Order:      66,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	expectedRequest := &v1.SaveMaterialGroupRequest{
		Id:         &wrapperspb.StringValue{Value: "5dbba1e31fd96208db5a00a1"},
		OrgId:      "5dbba1e31fd96208db5a00a1",
		IsValid:    &wrapperspb.BoolValue{Value: true},
		It:         &wrapperspb.Int32Value{Value: 55},
		Ut32:       &wrapperspb.UInt32Value{Value: 1},
		Ut64:       &wrapperspb.UInt64Value{Value: 2},
		StoryPoint: &wrapperspb.DoubleValue{Value: 3.5},
		Point:      &wrapperspb.FloatValue{Value: 5.5},
		Name:       "test1",
		Type:       &wrapperspb.StringValue{Value: "org"},
		Scope:      &wrapperspb.StringValue{Value: "group"},
		Order:      66,
		UpdateTime: timestamppb.New(time.Now()),
		CreateTime: timestamppb.New(time.Now()),
	}

	inputRequest := &v1.SaveMaterialGroupRequest{}

	expectedMaterial := &domain.MaterialGroup{}

	testCases := []struct {
		Name       string
		Req        []interface{}
		ErrorOccur bool
	}{
		{"copyByValidStruct", []interface{}{inputRequest, inputMaterial}, false},
		{"copyByValidStruct", []interface{}{expectedMaterial, expectedRequest}, false},
		{"copyByInValidStruct", []interface{}{map[string]string{"test": "test"}, inputMaterial}, true},
		{"copyByInValidStruct", []interface{}{expectedMaterial, map[string]string{"test": "test"}}, true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			err := copier.Copy(testCase.Req[0], testCases[1])
			if testCase.ErrorOccur {
				if err == nil {
					//zap.L().Error("expect occur exception, but not found")
					t.FailNow()
				}
			} else {
				if err != nil {
					//zap.L().Error("occur error:", zap.Error(err))
					t.FailNow()
				}
			}
		})
	}

}
