package common

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	commonpb "sensor-data-service.backend/api/pb/commonpb"
)

// WrapSuccess wraps any proto message into a StandardResponse with status = "success"
func WrapSuccess(msg string, pb proto.Message) (*commonpb.StandardResponse, error) {
	anyMsg, err := anypb.New(pb)
	if err != nil {
		return nil, err
	}
	return &commonpb.StandardResponse{
		Status:  "success",
		Message: msg,
		Data:    anyMsg,
	}, nil
}

// WrapError returns StandardResponse with error status (no data)
func WrapError(msg string) *commonpb.StandardResponse {
	return &commonpb.StandardResponse{
		Status:  "error",
		Message: msg,
	}
}
