package grpc

import (
	"context"

	ldapmanager "github.com/romnnn/ldap-manager"
	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewGroup ...
func (s *LDAPManagerServer) NewGroup(ctx context.Context, in *pb.NewGroupRequest) (*pb.Empty, error) {
	_, err := s.authenticate(ctx)
	if err != nil {
		return &pb.Empty{}, err
	}
	if err := s.Manager.NewGroup(in, false); err != nil {
		if appErr, safe := err.(ldapmanager.Error); safe {
			return &pb.Empty{}, toStatus(appErr)
		}
		log.Error(err)
		return &pb.Empty{}, status.Error(codes.Internal, "error while creating group")
	}
	return &pb.Empty{}, nil
}

// DeleteGroup ...
func (s *LDAPManagerServer) DeleteGroup(ctx context.Context, in *pb.DeleteGroupRequest) (*pb.Empty, error) {
	_, err := s.authenticate(ctx)
	if err != nil {
		return &pb.Empty{}, err
	}
	if err := s.Manager.DeleteGroup(in); err != nil {
		if appErr, safe := err.(ldapmanager.Error); safe {
			return &pb.Empty{}, toStatus(appErr)
		}
		log.Error(err)
		return &pb.Empty{}, status.Error(codes.Internal, "error while deleting group")
	}
	return &pb.Empty{}, nil
}

// UpdateGroup ...
func (s *LDAPManagerServer) UpdateGroup(ctx context.Context, in *pb.UpdateGroupRequest) (*pb.Empty, error) {
	_, err := s.authenticate(ctx)
	if err != nil {
		return &pb.Empty{}, err
	}
	if err := s.Manager.UpdateGroup(in); err != nil {
		if appErr, safe := err.(ldapmanager.Error); safe {
			return &pb.Empty{}, toStatus(appErr)
		}
		log.Error(err)
		return &pb.Empty{}, status.Error(codes.Internal, "error while updating group")
	}
	return &pb.Empty{}, nil
}

// GetGroupList ...
func (s *LDAPManagerServer) GetGroupList(ctx context.Context, in *pb.GetGroupListRequest) (*pb.GroupList, error) {
	_, err := s.authenticate(ctx)
	if err != nil {
		return &pb.GroupList{}, err
	}
	groups, err := s.Manager.GetGroupList(in)
	if err != nil {
		if appErr, safe := err.(ldapmanager.Error); safe {
			return &pb.GroupList{}, toStatus(appErr)
		}
		log.Error(err)
		return &pb.GroupList{}, status.Error(codes.Internal, "error while getting groups")
	}
	return groups, nil
}
