package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user-service_gc2p3/config"
	"user-service_gc2p3/pb"
)

type AuthServer struct {
	jwtManager *JWTManager
}

func NewAuthServer(jwtManager *JWTManager) *AuthServer {
	return &AuthServer{jwtManager}
}

func (server *AuthServer) Login(ctx context.Context, req *pb.LoginServiceRequest) (*pb.LoginServiceResponse, error) {

	// check service name and password if same from the this service from the env
	serviceName := config.Viper.GetString("USER_SERVICE_NAME")
	servicePassword := config.Viper.GetString("USER_SERVICE_PASSWORD")
	serviceRole := config.Viper.GetString("USER_SERVICE_ROLE")

	if req.GetServiceName() != serviceName || req.GetPassword() != servicePassword {
		return nil, status.Errorf(codes.PermissionDenied, "invalid service name or password")
	}

	token, err := server.jwtManager.Generate(serviceName, serviceRole)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while generating JWT token")
	}

	res := &pb.LoginServiceResponse{
		AccessToken: token,
	}

	return res, nil
}
