package grpc

import (
	"context"

	// ldapmanager "github.com/romnn/ldap-manager"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	// log "github.com/sirupsen/logrus"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

// IsGroupMember checks if an account is member of a group
func (s *LDAPManagerService) IsGroupMember(ctx context.Context, in *pb.IsGroupMemberRequest) (*pb.GroupMemberStatus, error) {
	// claims, err := s.authenticate(ctx)
	// if err != nil {
	// 	return &pb.GroupMemberStatus{}, err
	// }
	// if !claims.IsAdmin && claims.UID != in.GetUsername() {
	// 	return &pb.GroupMemberStatus{}, status.Error(codes.PermissionDenied, "requires admin privileges")
	// }
	// memberStatus, err := s.Manager.IsGroupMember(in)
	// if err != nil {
	// 	if appErr, safe := err.(ldapmanager.Error); safe {
	// 		return &pb.GroupMemberStatus{}, toStatus(appErr)
	// 	}
	// 	log.Error(err)
	// 	return &pb.GroupMemberStatus{}, status.Error(codes.Internal, "error while checking if user is member")
	// }
	// return memberStatus, nil
	return nil, nil
}

// GetGroup gets a group
func (s *LDAPManagerService) GetGroup(ctx context.Context, in *pb.GetGroupRequest) (*pb.Group, error) {
	// _, err := s.authenticate(ctx)
	// if err != nil {
	// 	return &pb.Group{}, err
	// }
	// group, err := s.Manager.GetGroup(in)
	// if err != nil {
	// 	if appErr, safe := err.(ldapmanager.Error); safe {
	// 		return &pb.Group{}, toStatus(appErr)
	// 	}
	// 	log.Error(err)
	// 	return &pb.Group{}, status.Error(codes.Internal, "error while getting group")
	// }
	// return group, nil
	return nil, nil
}

// GetUserGroups gets the groups an account is member of
func (s *LDAPManagerService) GetUserGroups(ctx context.Context, in *pb.GetUserGroupsRequest) (*pb.GroupList, error) {
	// claims, err := s.authenticate(ctx)
	// if err != nil {
	// 	return &pb.GroupList{}, err
	// }
	// if !claims.IsAdmin && claims.UID != in.GetUsername() {
	// 	return &pb.GroupList{}, status.Error(codes.PermissionDenied, "requires admin privileges")
	// }
	// groups, err := s.Manager.GetUserGroups(in)
	// if err != nil {
	// 	if appErr, safe := err.(ldapmanager.Error); safe {
	// 		return &pb.GroupList{}, toStatus(appErr)
	// 	}
	// 	log.Error(err)
	// 	return &pb.GroupList{}, status.Error(codes.Internal, "error while getting groups")
	// }
	// return groups, nil
	return nil, nil
}

// AddGroupMember adds a new member to a group
func (s *LDAPManagerService) AddGroupMember(ctx context.Context, in *pb.GroupMember) (*pb.Empty, error) {
	// _, err := s.authenticate(ctx)
	// if err != nil {
	// 	return &pb.Empty{}, err
	// }
	// if err := s.Manager.AddGroupMember(in, false); err != nil {
	// 	if appErr, safe := err.(ldapmanager.Error); safe {
	// 		return &pb.Empty{}, toStatus(appErr)
	// 	}
	// 	log.Error(err)
	// 	return &pb.Empty{}, status.Error(codes.Internal, "error while adding group member")
	// }
	// return &pb.Empty{}, nil
	return nil, nil
}

// DeleteGroupMember removes a member of a group
func (s *LDAPManagerService) DeleteGroupMember(ctx context.Context, in *pb.GroupMember) (*pb.Empty, error) {
	// claims, err := s.authenticate(ctx)
	// if err != nil {
	// 	return &pb.Empty{}, err
	// }
	// if !claims.IsAdmin && claims.UID != in.GetUsername() {
	// 	return &pb.Empty{}, status.Error(codes.PermissionDenied, "requires admin privileges")
	// }
	// allowDeleteOfDefaultGroups := claims.IsAdmin
	// if err := s.Manager.DeleteGroupMember(in, allowDeleteOfDefaultGroups); err != nil {
	// 	if appErr, safe := err.(ldapmanager.Error); safe {
	// 		return &pb.Empty{}, toStatus(appErr)
	// 	}
	// 	log.Error(err)
	// 	return &pb.Empty{}, status.Error(codes.Internal, "error while deleting group member")
	// }
	return &pb.Empty{}, nil
}
