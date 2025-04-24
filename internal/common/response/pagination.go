package response

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	commonpb "sensor-data-service.backend/api/pb/commonpb"
)

// WrapPagination wraps a paginated response
func WrapPagination(message string, items proto.Message, page, limit, total int32) (*commonpb.StandardResponse, error) {
	itemsAny, err := anypb.New(items)
	if err != nil {
		return nil, err
	}

	paginated := &commonpb.PaginatedData{
		Items: itemsAny,
		Pagination: &commonpb.PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}

	paginatedAny, err := anypb.New(paginated)
	if err != nil {
		return nil, err
	}

	return &commonpb.StandardResponse{
		Status:  "success",
		Message: message,
		Data:    paginatedAny,
	}, nil
}
