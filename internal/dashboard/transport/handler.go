package transport

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"
	"sensor-data-service.backend/api/pb/commonpb"
	"sensor-data-service.backend/api/pb/dashboardpb"
	"sensor-data-service.backend/internal/common"
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

func (h *DashboardHandler) GetDashboard(ctx context.Context, req *dashboardpb.DashboardID) (*commonpb.StandardResponse, error) {
	db, err := h.svc.GetDashboardByID(ctx, req.Id)
	if err != nil {
		return common.WrapError(fmt.Sprintf("Failed to get dashboard: %v", err)), nil
	}
	return common.WrapSuccess("GetDashboard OK", toProto(db))
}

func (h *DashboardHandler) ListDashboards(ctx context.Context, _ *emptypb.Empty) (*commonpb.StandardResponse, error) {
	list, err := h.svc.ListDashboards(ctx)
	if err != nil {
		return common.WrapError(fmt.Sprintf("Failed to list dashboards: %v", err)), nil
	}

	var pbList []*dashboardpb.Dashboard
	for _, d := range list {
		pbList = append(pbList, toProto(d))
	}

	return common.WrapSuccess("ListDashboards OK", &dashboardpb.DashboardList{Dashboards: pbList})
}

func (h *DashboardHandler) SaveDashboard(ctx context.Context, req *dashboardpb.SaveDashboardRequest) (*commonpb.StandardResponse, error) {
	d := req.GetDashboard()
	err := h.svc.SaveDashboard(ctx, &domain.Dashboard{
		ID:          d.Id,
		Name:        d.Name,
		Description: d.Description,
		LayoutJSON:  d.LayoutJson,
		CreatedBy:   d.CreatedBy,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		Version:     d.Version,
		Status:      d.Status,
	})
	if err != nil {
		return common.WrapError(fmt.Sprintf("Failed to save dashboard: %v", err)), nil
	}
	return common.WrapSuccess("SaveDashboard OK", nil)
}

func (h *DashboardHandler) DeleteDashboard(ctx context.Context, req *dashboardpb.DashboardID) (*commonpb.StandardResponse, error) {
	err := h.svc.DeleteDashboard(ctx, req.Id)
	if err != nil {
		return common.WrapError(fmt.Sprintf("Failed to delete dashboard: %v", err)), nil
	}
	return common.WrapSuccess("DeleteDashboard OK", nil)
}

func toProto(d *domain.Dashboard) *dashboardpb.Dashboard {
	return &dashboardpb.Dashboard{
		Id:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		LayoutJson:  d.LayoutJSON,
		CreatedBy:   d.CreatedBy,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		Version:     d.Version,
		Status:      d.Status,
	}
}
