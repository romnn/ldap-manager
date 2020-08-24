package base

import (
	"context"

	"github.com/neko-neko/echo-logrus/v2/log"
	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
)

// NewGroup ...
func (s *LDAPManagerServer) NewGroup(ctx context.Context, in *pb.NewGroupRequest) (*pb.Empty, error) {
	var result pb.Empty
	log.Info("NewGroup")
	return &result, nil
}

// DeleteGroup ...
func (s *LDAPManagerServer) DeleteGroup(ctx context.Context, in *pb.DeleteGroupRequest) (*pb.Empty, error) {
	var result pb.Empty
	log.Info("DeleteGroup")
	return &result, nil
}

// RenameGroup ...
func (s *LDAPManagerServer) RenameGroup(ctx context.Context, in *pb.RenameGroupRequest) (*pb.Empty, error) {
	var result pb.Empty
	log.Info("RenameGroup")
	return &result, nil
}

// GetGroupList ...
func (s *LDAPManagerServer) GetGroupList(ctx context.Context, in *pb.GetGroupListRequest) (*pb.GroupList, error) {
	var result pb.GroupList
	log.Info("GetGroupList")
	return &result, nil
}
