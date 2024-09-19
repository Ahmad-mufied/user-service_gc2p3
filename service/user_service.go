package service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"user-service_gc2p3/entity"
	"user-service_gc2p3/pb"
)

type UserService struct {
	db *mongo.Client
}

func NewUserService(db *mongo.Client) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	username := req.GetUsername()
	password := req.GetPassword()

	// check if user exists
	coll := s.db.Database("user_db").Collection("users")
	filter := bson.M{"username": username}
	var user entity.User
	err := coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	// return response
	res := &pb.LoginUserResponse{
		Success: true,
		Message: "login success",
		UserId:  user.Id.String(),
	}

	return res, nil
}

func (s *UserService) Register(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	username := req.GetUsername()
	password := req.GetPassword()

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// insert user
	coll := s.db.Database("user_db").Collection("users")
	user := entity.User{
		Username: username,
		Password: string(hashedPassword),
	}
	_, err = coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	// return response
	res := &pb.RegisterUserResponse{
		Message: "register success",
	}

	return res, nil
}

func (s *UserService) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {
	id := req.GetUserId()

	// get user by id
	coll := s.db.Database("user_db").Collection("users")
	filter := bson.M{"_id": id}
	var user entity.User
	err := coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	// return response
	res := &pb.GetUserByIdResponse{
		Username: user.Username,
	}

	return res, nil
}
