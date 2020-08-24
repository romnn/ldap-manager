package grpc

import (
	"context"

	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
)

// NewGroup ...
func (s *LDAPManagerServer) NewGroup(ctx context.Context, in *pb.NewGroupRequest) (*pb.Empty, error) {
	return s.LDAPManagerServer.NewGroup(ctx, in)
}

// DeleteGroup ...
func (s *LDAPManagerServer) DeleteGroup(ctx context.Context, in *pb.DeleteGroupRequest) (*pb.Empty, error) {
	return s.LDAPManagerServer.DeleteGroup(ctx, in)
}

// RenameGroup ...
func (s *LDAPManagerServer) RenameGroup(ctx context.Context, in *pb.RenameGroupRequest) (*pb.Empty, error) {
	return s.LDAPManagerServer.RenameGroup(ctx, in)
}

// GetGroupList ...
func (s *LDAPManagerServer) GetGroupList(ctx context.Context, in *pb.GetGroupListRequest) (*pb.GroupList, error) {
	return s.LDAPManagerServer.GetGroupList(ctx, in)
}
