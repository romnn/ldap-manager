package base

import (
	"context"

	"github.com/neko-neko/echo-logrus/v2/log"
	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
)

// IsGroupMember ...
func (s *LDAPManagerServer) IsGroupMember(ctx context.Context, in *pb.IsGroupMemberRequest) (*pb.GroupMemberStatus, error) {
	var result pb.GroupMemberStatus
	log.Info("IsGroupMember")
	return &result, nil
}

// GetGroup ...
func (s *LDAPManagerServer) GetGroup(ctx context.Context, in *pb.GetGroupRequest) (*pb.Group, error) {
	var result pb.Group
	log.Info("GetGroup")
	return &result, nil
}

// AddGroupMember ...
func (s *LDAPManagerServer) AddGroupMember(ctx context.Context, in *pb.GroupMember) (*pb.Empty, error) {
	var result pb.Empty
	log.Info("AddGroupMember")
	return &result, nil
}

// DeleteGroupMember ...
func (s *LDAPManagerServer) DeleteGroupMember(ctx context.Context, in *pb.GroupMember) (*pb.Empty, error) {
	var result pb.Empty
	log.Info("DeleteGroupMember")
	return &result, nil
}
