package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"user-service_gc2p3/config"
	"user-service_gc2p3/pb"
	"user-service_gc2p3/service"
)

func accessibleRoles() map[string][]string {
	const (
		userServicePath = "/user.UserService/"
	)

	return map[string][]string{
		userServicePath + "Create": {"admin"},
		userServicePath + "Read":   {"admin", "user"},
	}
}

func main() {
	config.InitViper()
	config.InitMongo()

	secretKey := config.Viper.GetString("USER_SERVICE_JWT_SECRET")
	duration := config.Viper.GetDuration("USER_SERVICE_JWT_DURATION")

	jwtManager := service.NewJWTManager(secretKey, duration)
	authServer := service.NewAuthServer(jwtManager)

	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)

	userServer := service.NewUserService(config.DB)

	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterAuthUserServiceServer(grpcServer, userServer)
	reflection.Register(grpcServer)

	gRPCPort := config.Viper.GetString("PORT")
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPort))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server running on port :", gRPCPort)
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatal(err)
	}

}
