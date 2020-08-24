package grpc

import (
	"context"

	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
)

// IsGroupMember ...
func (s *LDAPManagerServer) IsGroupMember(ctx context.Context, in *pb.IsGroupMemberRequest) (*pb.GroupMemberStatus, error) {
	return s.LDAPManagerServer.IsGroupMember(ctx, in)
}

// GetGroup ...
func (s *LDAPManagerServer) GetGroup(ctx context.Context, in *pb.GetGroupRequest) (*pb.Group, error) {
	return s.LDAPManagerServer.GetGroup(ctx, in)
}

// AddGroupMember ...
func (s *LDAPManagerServer) AddGroupMember(ctx context.Context, in *pb.GroupMember) (*pb.Empty, error) {
	return s.LDAPManagerServer.AddGroupMember(ctx, in)
}

// DeleteGroupMember ...
func (s *LDAPManagerServer) DeleteGroupMember(ctx context.Context, in *pb.GroupMember) (*pb.Empty, error) {
	return s.LDAPManagerServer.DeleteGroupMember(ctx, in)
}
