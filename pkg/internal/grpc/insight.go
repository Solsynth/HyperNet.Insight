package grpc

import (
	"context"

	"git.solsynth.dev/hypernet/insight/pkg/internal/services"
	"git.solsynth.dev/hypernet/insight/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (v *Server) GenerateInsight(ctx context.Context, request *proto.InsightRequest) (*proto.InsightResponse, error) {
	input := request.GetSource()
	if err := services.PlaceOrder(uint(request.GetUserId()), len(input)); err != nil {
		return nil, status.Errorf(codes.ResourceExhausted, "failed to place order: %v", err)
	}

	out, err := services.GenerateInsights(input)
	if err != nil {
		_ = services.MakeRefund(uint(request.GetUserId()), len(input))
		return nil, status.Errorf(codes.Internal, "failed to generate insight: %v", err)
	}

	return &proto.InsightResponse{Response: out}, nil
}
