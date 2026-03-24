package service

import (
	"context"
	"strings"
	"time"

	"fooder-backend/core/errors"
	detectv1 "fooder-backend/gen/go/detect/v1"
	"fooder-backend/internal/dto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DetectService struct {
	client  detectv1.FoodDetectServiceClient
	conn    *grpc.ClientConn
	timeout time.Duration
}

func NewDetectService(address string) (*DetectService, error) {
	addr := strings.TrimSpace(address)
	if addr == "" {
		addr = "localhost:50051"
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &DetectService{
		client:  detectv1.NewFoodDetectServiceClient(conn),
		conn:    conn,
		timeout: 10 * time.Second,
	}, nil
}

func (s *DetectService) Close() error {
	if s == nil || s.conn == nil {
		return nil
	}
	return s.conn.Close()
}

func (s *DetectService) AnalyzeFood(ctx context.Context, input dto.AnalyzeFoodRequest) (*dto.AnalyzeFoodResponse, *errors.AppError) {
	objectKey := strings.TrimSpace(input.ObjectKey)
	if objectKey == "" {
		return nil, errors.NewAppError(errors.ErrInvalidInput, "object_key is required", nil)
	}

	callCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	result, err := s.client.AnalyzeFood(callCtx, &detectv1.AnalyzeFoodRequest{
		ObjectKey: objectKey,
	})
	if err != nil {
		return nil, errors.NewAppError(errors.ErrThirdParty, "detect service unavailable", err)
	}

	items := make([]dto.FoodItem, 0, len(result.GetItems()))
	for _, item := range result.GetItems() {
		items = append(items, dto.FoodItem{
			Name:        item.GetName(),
			Confidence:  item.GetConfidence(),
			CaloriesEst: item.GetCaloriesEst(),
		})
	}

	return &dto.AnalyzeFoodResponse{
		Items: items,
		Note:  result.GetNote(),
	}, nil
}
