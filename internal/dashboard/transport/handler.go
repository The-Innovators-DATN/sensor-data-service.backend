package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"sensor-data-service.backend/api/pb/commonpb"
	"sensor-data-service.backend/api/pb/dashboardpb"
	"sensor-data-service.backend/internal/common"
	"sensor-data-service.backend/internal/common/castutil"
	"sensor-data-service.backend/internal/dashboard/domain"

	"sensor-data-service.backend/internal/dashboard/service"
)

type DashboardHandler struct {
	dashboardpb.UnimplementedDashboardServiceServer
	svc service.DashboardService
}

func NewDashboardHandler(svc service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) GetDashboard(ctx context.Context, req *dashboardpb.GetDashboardRequest) (*commonpb.StandardResponse, error) {
	db, err := h.svc.GetDashboardByID(ctx, req.Uid, req.CreatedBy)
	if err != nil {
		return common.WrapError(fmt.Sprintf("Failed to get dashboard: %v", err)), nil
	}
	return common.WrapSuccess("GetDashboard OK", toProto(db))
}

func (h *DashboardHandler) ListDashboards(ctx context.Context, req *dashboardpb.CreatedBy) (*commonpb.StandardResponse, error) {

	list, err := h.svc.ListDashboardsByUser(ctx, req.CreatedBy)
	if err != nil {
		return common.WrapError(fmt.Sprintf("Failed to list dashboards: %v", err)), nil
	}
	var protoList []*dashboardpb.Dashboard
	for _, d := range list {
		protoList = append(protoList, toProto(d))
	}
	return common.WrapSuccess("ListDashboards OK", &dashboardpb.DashboardList{Dashboards: protoList})
}

func (h *DashboardHandler) CreateDashboard(ctx context.Context, req *dashboardpb.CreateDashboardRequest) (*commonpb.StandardResponse, error) {

	if err := h.svc.CreateDashboard(ctx, fromProto(req.Dashboard)); err != nil {
		return common.WrapError(fmt.Sprintf("Failed to create dashboard: %v", err)), nil
	}
	return common.WrapSuccess("CreateDashboard OK", &emptypb.Empty{})
}

func (h *DashboardHandler) UpdateDashboard(ctx context.Context, req *dashboardpb.UpdateDashboardRequest) (*commonpb.StandardResponse, error) {
	if err := h.svc.UpdateDashboard(ctx, fromProto(req.Dashboard), req.CreatedBy); err != nil {
		return common.WrapError(fmt.Sprintf("Failed to update dashboard: %v", err)), nil
	}
	return common.WrapSuccess("UpdateDashboard OK", &emptypb.Empty{})
}

func (h *DashboardHandler) PatchDashboard(ctx context.Context, req *dashboardpb.PatchDashboardRequest) (*commonpb.StandardResponse, error) {
	d := &domain.Dashboard{
		UID:                 castutil.MustParseUUID(req.Uid),
		Name:                req.Name,
		Description:         req.Description,
		LayoutConfiguration: req.LayoutConfiguration,
		Status:              req.Status,
	}
	if err := h.svc.PatchDashboard(ctx, d, req.CreatedBy); err != nil {
		return common.WrapError(fmt.Sprintf("Failed to patch dashboard: %v", err)), nil
	}
	return common.WrapSuccess("PatchDashboard OK", &emptypb.Empty{})
}

func (h *DashboardHandler) DeleteDashboard(ctx context.Context, req *dashboardpb.DeleteDashboardRequest) (*commonpb.StandardResponse, error) {
	if err := h.svc.DeleteDashboard(ctx, req.Uid, req.CreatedBy); err != nil {
		return common.WrapError(fmt.Sprintf("Failed to delete dashboard: %v", err)), nil
	}
	return common.WrapSuccess("DeleteDashboard OK", &emptypb.Empty{})
}

// Mapping functions

func toProto(d *domain.Dashboard) *dashboardpb.Dashboard {
	var layoutMap map[string]interface{}
	_ = json.Unmarshal([]byte(d.LayoutConfiguration), &layoutMap)
	layoutStruct, _ := structpb.NewStruct(layoutMap)

	return &dashboardpb.Dashboard{
		Uid:                 d.UID.String(),
		Name:                d.Name,
		Description:         d.Description,
		LayoutConfiguration: layoutStruct,
		CreatedBy:           d.CreatedBy,
		CreatedAt:           d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           d.UpdatedAt.Format(time.RFC3339),
		Version:             d.Version,
		Status:              d.Status,
	}
}

func fromProto(d *dashboardpb.Dashboard) *domain.Dashboard {
	layoutJSON, _ := json.Marshal(d.LayoutConfiguration) // []byte
	layoutStr := string(layoutJSON)
	// string
	var uid uuid.UUID
	if d.Uid != "" {
		uid = castutil.MustParseUUID(d.Uid)
	}
	return &domain.Dashboard{
		UID:                 uid,
		Name:                d.Name,
		Description:         d.Description,
		LayoutConfiguration: layoutStr, // <-- string JSON
		CreatedBy:           d.CreatedBy,
		CreatedAt:           castutil.ToTime(d.CreatedAt),
		UpdatedAt:           castutil.ToTime(d.UpdatedAt),
		Version:             d.Version,
		Status:              d.Status,
	}
}
