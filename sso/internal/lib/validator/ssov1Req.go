package validator

import (
	ssov1 "ChatService/protos/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyValue = 0
)

func LoginValid(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password required")
	}
	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app id required")
	}
	return nil
}

func LogoutValid(req *ssov1.LogoutRequest) error {
	return nil
}

func RegisterValid(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password required")
	}
	return nil
}

func IsAdminValid(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user id required")
	}
	return nil
}

func IsModValid(req *ssov1.IsModeratorRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user id required")
	}
	return nil
}

func ChangePasswordValid(req *ssov1.ChangePasswordRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user id required")
	}
	if req.GetNewPassword() == "" {
		return status.Error(codes.InvalidArgument, "new password can not be empty")
	}
	if req.GetOldPassword() == "" {
		return status.Error(codes.InvalidArgument, "old password required")
	}
	return nil
}

func ChangeNameValid(req *ssov1.ChangeNameRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user id required")
	}
	if req.GetNewName() == "" {
		return status.Error(codes.InvalidArgument, "new name can not be empty")
	}
	return nil
}

func ChangeRoleValid(req *ssov1.ChangeRoleRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user id required")
	}
	if req.GetAdminId() == emptyValue {
		return status.Error(codes.InvalidArgument, "admin id required")
	}
	if !RoleValid(req.NewRole) {
		return status.Error(codes.InvalidArgument, "new role must have a valid role")
	}
	return nil
}

func RoleValid(role int32) bool {
	if role > 4 || role < 1 {
		return false
	}
	return true
}
