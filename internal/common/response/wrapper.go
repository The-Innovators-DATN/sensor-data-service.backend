package response

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	commonpb "sensor-data-service.backend/api/pb/commonpb"
)

// WrapSuccess returns a standard success response
func WrapSuccess(message string, pb proto.Message) (*commonpb.StandardResponse, error) {
	var data *anypb.Any
	if pb != nil {
		var err error
		data, err = anypb.New(pb)
		if err != nil {
			return nil, err
		}
	}
	return &commonpb.StandardResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}, nil
}

// WrapCreated returns a response similar to HTTP 201 Created (semantically)
func WrapCreated(message string, pb proto.Message) (*commonpb.StandardResponse, error) {
	return WrapSuccess("Created: "+message, pb)
}

// WrapError returns a standard error response, optionally with error detail
func WrapError(message string, errDetail proto.Message) (*commonpb.StandardResponse, error) {
	var detail *anypb.Any
	if errDetail != nil {
		var err error
		detail, err = anypb.New(errDetail)
		if err != nil {
			return nil, err
		}
	}
	return &commonpb.StandardResponse{
		Status:      "error",
		Message:     message,
		ErrorDetail: detail,
	}, nil
}
func Unauthorized(message string) (*commonpb.StandardResponse, error) {
	return WrapError("Unauthorized: "+message, nil)
}

func BadRequest(message string, errDetail proto.Message) (*commonpb.StandardResponse, error) {
	return WrapError("Bad Request: "+message, errDetail)
}

func NotFound(message string) (*commonpb.StandardResponse, error) {
	return WrapError("Not Found: "+message, nil)
}

func InternalServer(message string, errDetail proto.Message) (*commonpb.StandardResponse, error) {
	return WrapError("Internal Server Error: "+message, errDetail)
}
