package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"

	"sensor-data-service.backend/api/pb/commonpb"
	"sensor-data-service.backend/api/pb/dashboardpb"
	"sensor-data-service.backend/internal/common/castutil"
	"sensor-data-service.backend/internal/common/response"
	"sensor-data-service.backend/internal/domain/model"
	"sensor-data-service.backend/internal/domain/service"
	"sensor-data-service.backend/internal/transport/utils"
)

type DashboardHandler struct {
	dashboardpb.UnimplementedDashboardServiceServer
	svc service.DashboardService
}

func NewDashboardHandler(svc service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) GetDashboard(ctx context.Context, req *dashboardpb.GetDashboardRequest) (*commonpb.StandardResponse, error) {
	// userID, _ := ctx.Value("user_id").(int32)
	// log.Printf("[info] GetDashboard: %s", userID)
	db, err := h.svc.GetDashboardByID(ctx, req.Uid, req.CreatedBy)
	if err != nil {
		return response.WrapError(fmt.Sprintf("Failed to get dashboard: %v", err), nil)
	}
	return response.WrapSuccess("GetDashboard OK", toProto(db))
}

func (h *DashboardHandler) ListDashboards(ctx context.Context, req *dashboardpb.PaginateDashboardsRequest) (*commonpb.StandardResponse, error) {
	page, limit := utils.SanitizePagination(req.Page, req.Limit)
	list, err := h.svc.ListDashboardsByUser(ctx, req.CreatedBy, page, limit)
	if err != nil {
		return response.WrapError(fmt.Sprintf("Failed to list dashboards: %v", err), nil)
	}
	var protoList []*dashboardpb.Dashboard
	for _, d := range list {
		protoList = append(protoList, toProto(d))
	}
	return response.WrapSuccess("ListDashboards OK", &dashboardpb.DashboardList{Dashboards: protoList})
}

func (h *DashboardHandler) CreateDashboard(ctx context.Context, req *dashboardpb.CreateDashboardRequest) (*commonpb.StandardResponse, error) {
	dashboard := fromProto(req.Dashboard)
	log.Printf("Full request: %v", req)
	dashboardID, err := h.svc.CreateDashboard(ctx, dashboard)
	if err != nil {
		return response.WrapError(fmt.Sprintf("failed to create dashboard: %v", err), nil)
	}

	return response.WrapSuccess("createDashboard OK", &dashboardpb.DashboardID{
		Uid: dashboardID,
	})
}

func (h *DashboardHandler) UpdateDashboard(ctx context.Context, req *dashboardpb.UpdateDashboardRequest) (*commonpb.StandardResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return response.WrapError(fmt.Sprintf("Invalid UID format: %v", err), nil)
	}
	log.Printf("[info] UpdateDashboard: name=%s", req.Dashboard.Name)
	dashboard := fromProto(req.Dashboard)
	dashboard.UID = uid
	dashboard.Version = 0 // Reset version as it should not be input from the proto
	if err := h.svc.UpdateDashboard(ctx, dashboard, req.CreatedBy); err != nil {
		return response.WrapError(fmt.Sprintf("Failed to update dashboard: %v", err), nil)
	}
	return response.WrapSuccess("UpdateDashboard OK", &emptypb.Empty{})
}

func (h *DashboardHandler) PatchDashboard(ctx context.Context, req *dashboardpb.PatchDashboardRequest) (*commonpb.StandardResponse, error) {
	d := &model.Dashboard{
		UID:                 castutil.MustParseUUID(req.Uid),
		Name:                req.Name,
		Description:         req.Description,
		LayoutConfiguration: req.LayoutConfiguration,
		Status:              req.Status,
	}
	if err := h.svc.PatchDashboard(ctx, d, req.CreatedBy); err != nil {
		return response.WrapError(fmt.Sprintf("Failed to patch dashboard: %v", err), nil)
	}
	return response.WrapSuccess("PatchDashboard OK", &emptypb.Empty{})
}

func (h *DashboardHandler) DeleteDashboard(ctx context.Context, req *dashboardpb.DeleteDashboardRequest) (*commonpb.StandardResponse, error) {
	if err := h.svc.DeleteDashboard(ctx, req.Uid, req.CreatedBy); err != nil {
		return response.WrapError(fmt.Sprintf("Failed to delete dashboard: %v", err), nil)
	}
	return response.WrapSuccess("DeleteDashboard OK", &emptypb.Empty{})
}

// Mapping functions

func toProto(d *model.Dashboard) *dashboardpb.Dashboard {
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

func fromProto(d *dashboardpb.Dashboard) *model.Dashboard {
	if d.LayoutConfiguration == nil {
		log.Printf("[info] fromProto: layout_configuration is nil")
		// Print createdBy and type of createdBy
		log.Printf("[info] fromProto: createdBy=%d", d.CreatedBy)
		log.Printf("[info] fromProto: type of createdBy=%T", d.CreatedBy)
		return &model.Dashboard{
			UID:                 uuid.New(),
			Name:                d.Name,
			Description:         d.Description,
			LayoutConfiguration: "{}",
			CreatedBy:           d.CreatedBy,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			Version:             1,
			Status:              "active",
		}
	}
	layoutJSON, _ := json.Marshal(d.LayoutConfiguration)
	layoutStr := string(layoutJSON)
	var uid uuid.UUID
	if d.Uid != "" {
		uid = castutil.MustParseUUID(d.Uid)
	}
	var version int32 = 1 // Default version if null
	if d.Version != 0 {
		version = d.Version
	}
	log.Printf("[info] fromProto: uid=%s", uid)
	log.Printf("[info] fromProto: layout=%s", layoutStr)
	log.Printf("[info] fromProto: created_by=%d", d.CreatedBy)
	return &model.Dashboard{
		UID:                 uid,
		Name:                d.Name,
		Description:         d.Description,
		LayoutConfiguration: layoutStr,
		CreatedBy:           d.CreatedBy,
		CreatedAt:           castutil.ToTime(d.CreatedAt),
		UpdatedAt:           castutil.ToTime(d.UpdatedAt),
		Version:             version,
		Status:              d.Status,
	}
}
