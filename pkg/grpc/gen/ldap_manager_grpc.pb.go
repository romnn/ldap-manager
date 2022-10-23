// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.17.3
// source: ldap_manager.proto

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// LDAPManagerClient is the client API for LDAPManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LDAPManagerClient interface {
	// Authentication
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*Token, error)
	// Users
	GetUserList(ctx context.Context, in *GetUserListRequest, opts ...grpc.CallOption) (*UserList, error)
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*UserData, error)
	NewUser(ctx context.Context, in *NewUserRequest, opts ...grpc.CallOption) (*Empty, error)
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*Token, error)
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*Empty, error)
	ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...grpc.CallOption) (*Empty, error)
	// Groups
	NewGroup(ctx context.Context, in *NewGroupRequest, opts ...grpc.CallOption) (*Empty, error)
	DeleteGroup(ctx context.Context, in *DeleteGroupRequest, opts ...grpc.CallOption) (*Empty, error)
	UpdateGroup(ctx context.Context, in *UpdateGroupRequest, opts ...grpc.CallOption) (*Empty, error)
	GetGroupList(ctx context.Context, in *GetGroupListRequest, opts ...grpc.CallOption) (*GroupList, error)
	GetUserGroups(ctx context.Context, in *GetUserGroupsRequest, opts ...grpc.CallOption) (*GroupList, error)
	// Group members
	IsGroupMember(ctx context.Context, in *IsGroupMemberRequest, opts ...grpc.CallOption) (*GroupMemberStatus, error)
	GetGroup(ctx context.Context, in *GetGroupRequest, opts ...grpc.CallOption) (*Group, error)
	AddGroupMember(ctx context.Context, in *GroupMember, opts ...grpc.CallOption) (*Empty, error)
	DeleteGroupMember(ctx context.Context, in *GroupMember, opts ...grpc.CallOption) (*Empty, error)
}

type lDAPManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewLDAPManagerClient(cc grpc.ClientConnInterface) LDAPManagerClient {
	return &lDAPManagerClient{cc}
}

func (c *lDAPManagerClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*Token, error) {
	out := new(Token)
	err := c.cc.Invoke(ctx, "/ldapmanager.LDAPManager/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lDAPManagerClient) GetUserList(ctx context.Context, in *GetUserListRequest, opts ...grpc.CallOption) (*UserList, error) {
	out := new(UserList)
	err := c.cc.Invoke(ctx, "/ldapmanager.LDAPManager/GetUserList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lDAPManagerClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*UserData, error) {
	out := new(UserData)
	err := c.cc.Invoke(ctx, "/ldapmanager.LDAPManager/GetUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lDAPManagerClient) NewUser(ctx context.Context, in *NewUserRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/ldapmanager.LDAPManager/NewUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lDAPManagerClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*Token, error) {
	out := new(Token)
	err := c.cc.Invoke(ctx, "/ldapmanager.LDAPManager/UpdateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lDAPManagerClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/ldapmanager.LDAPManager/DeleteUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lDAPManagerClient) ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/ldapmanager.LDAPManager/ChangePassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lDAPManagerClient) NewGroup(ctx context.Context, in *NewGroupRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/ldapmanager.LDAPManager/NewGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lDAPManagerClient) DeleteGroup(ctx context.Context, in *DeleteGroupRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/ldapmanager.LDAPManager/DeleteGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lDAPManagerClient) UpdateGroup(ctx context.Context, in *UpdateGroupRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/ldapmanager.LDAPManager/UpdateGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lDAPManagerClient) GetGroupList(ctx context.Context, in *GetGroupListRequest, opts ...grpc.CallOption) (*GroupList, error) {
	out := new(GroupList)
	err := c.cc.Invoke(ctx, "/ldapmanager.LDAPManager/GetGroupList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lDAPManagerClient) GetUserGroups(ctx context.Context, in *GetUserGroupsRequest, opts ...grpc.CallOption) (*GroupList, error) {
	out := new(GroupList)
	err := c.cc.Invoke(ctx, "/ldapmanager.LDAPManager/GetUserGroups", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lDAPManagerClient) IsGroupMember(ctx context.Context, in *IsGroupMemberRequest, opts ...grpc.CallOption) (*GroupMemberStatus, error) {
	out := new(GroupMemberStatus)
	err := c.cc.Invoke(ctx, "/ldapmanager.LDAPManager/IsGroupMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lDAPManagerClient) GetGroup(ctx context.Context, in *GetGroupRequest, opts ...grpc.CallOption) (*Group, error) {
	out := new(Group)
	err := c.cc.Invoke(ctx, "/ldapmanager.LDAPManager/GetGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lDAPManagerClient) AddGroupMember(ctx context.Context, in *GroupMember, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/ldapmanager.LDAPManager/AddGroupMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lDAPManagerClient) DeleteGroupMember(ctx context.Context, in *GroupMember, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/ldapmanager.LDAPManager/DeleteGroupMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LDAPManagerServer is the server API for LDAPManager service.
// All implementations must embed UnimplementedLDAPManagerServer
// for forward compatibility
type LDAPManagerServer interface {
	// Authentication
	Login(context.Context, *LoginRequest) (*Token, error)
	// Users
	GetUserList(context.Context, *GetUserListRequest) (*UserList, error)
	GetUser(context.Context, *GetUserRequest) (*UserData, error)
	NewUser(context.Context, *NewUserRequest) (*Empty, error)
	UpdateUser(context.Context, *UpdateUserRequest) (*Token, error)
	DeleteUser(context.Context, *DeleteUserRequest) (*Empty, error)
	ChangePassword(context.Context, *ChangePasswordRequest) (*Empty, error)
	// Groups
	NewGroup(context.Context, *NewGroupRequest) (*Empty, error)
	DeleteGroup(context.Context, *DeleteGroupRequest) (*Empty, error)
	UpdateGroup(context.Context, *UpdateGroupRequest) (*Empty, error)
	GetGroupList(context.Context, *GetGroupListRequest) (*GroupList, error)
	GetUserGroups(context.Context, *GetUserGroupsRequest) (*GroupList, error)
	// Group members
	IsGroupMember(context.Context, *IsGroupMemberRequest) (*GroupMemberStatus, error)
	GetGroup(context.Context, *GetGroupRequest) (*Group, error)
	AddGroupMember(context.Context, *GroupMember) (*Empty, error)
	DeleteGroupMember(context.Context, *GroupMember) (*Empty, error)
	mustEmbedUnimplementedLDAPManagerServer()
}

// UnimplementedLDAPManagerServer must be embedded to have forward compatible implementations.
type UnimplementedLDAPManagerServer struct {
}

func (UnimplementedLDAPManagerServer) Login(context.Context, *LoginRequest) (*Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedLDAPManagerServer) GetUserList(context.Context, *GetUserListRequest) (*UserList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserList not implemented")
}
func (UnimplementedLDAPManagerServer) GetUser(context.Context, *GetUserRequest) (*UserData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedLDAPManagerServer) NewUser(context.Context, *NewUserRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewUser not implemented")
}
func (UnimplementedLDAPManagerServer) UpdateUser(context.Context, *UpdateUserRequest) (*Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (UnimplementedLDAPManagerServer) DeleteUser(context.Context, *DeleteUserRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (UnimplementedLDAPManagerServer) ChangePassword(context.Context, *ChangePasswordRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePassword not implemented")
}
func (UnimplementedLDAPManagerServer) NewGroup(context.Context, *NewGroupRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewGroup not implemented")
}
func (UnimplementedLDAPManagerServer) DeleteGroup(context.Context, *DeleteGroupRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteGroup not implemented")
}
func (UnimplementedLDAPManagerServer) UpdateGroup(context.Context, *UpdateGroupRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateGroup not implemented")
}
func (UnimplementedLDAPManagerServer) GetGroupList(context.Context, *GetGroupListRequest) (*GroupList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupList not implemented")
}
func (UnimplementedLDAPManagerServer) GetUserGroups(context.Context, *GetUserGroupsRequest) (*GroupList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserGroups not implemented")
}
func (UnimplementedLDAPManagerServer) IsGroupMember(context.Context, *IsGroupMemberRequest) (*GroupMemberStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsGroupMember not implemented")
}
func (UnimplementedLDAPManagerServer) GetGroup(context.Context, *GetGroupRequest) (*Group, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroup not implemented")
}
func (UnimplementedLDAPManagerServer) AddGroupMember(context.Context, *GroupMember) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddGroupMember not implemented")
}
func (UnimplementedLDAPManagerServer) DeleteGroupMember(context.Context, *GroupMember) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteGroupMember not implemented")
}
func (UnimplementedLDAPManagerServer) mustEmbedUnimplementedLDAPManagerServer() {}

// UnsafeLDAPManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LDAPManagerServer will
// result in compilation errors.
type UnsafeLDAPManagerServer interface {
	mustEmbedUnimplementedLDAPManagerServer()
}

func RegisterLDAPManagerServer(s grpc.ServiceRegistrar, srv LDAPManagerServer) {
	s.RegisterService(&LDAPManager_ServiceDesc, srv)
}

func _LDAPManager_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LDAPManagerServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ldapmanager.LDAPManager/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LDAPManagerServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LDAPManager_GetUserList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LDAPManagerServer).GetUserList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ldapmanager.LDAPManager/GetUserList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LDAPManagerServer).GetUserList(ctx, req.(*GetUserListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LDAPManager_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LDAPManagerServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ldapmanager.LDAPManager/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LDAPManagerServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LDAPManager_NewUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LDAPManagerServer).NewUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ldapmanager.LDAPManager/NewUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LDAPManagerServer).NewUser(ctx, req.(*NewUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LDAPManager_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LDAPManagerServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ldapmanager.LDAPManager/UpdateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LDAPManagerServer).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LDAPManager_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LDAPManagerServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ldapmanager.LDAPManager/DeleteUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LDAPManagerServer).DeleteUser(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LDAPManager_ChangePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangePasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LDAPManagerServer).ChangePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ldapmanager.LDAPManager/ChangePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LDAPManagerServer).ChangePassword(ctx, req.(*ChangePasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LDAPManager_NewGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LDAPManagerServer).NewGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ldapmanager.LDAPManager/NewGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LDAPManagerServer).NewGroup(ctx, req.(*NewGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LDAPManager_DeleteGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LDAPManagerServer).DeleteGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ldapmanager.LDAPManager/DeleteGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LDAPManagerServer).DeleteGroup(ctx, req.(*DeleteGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LDAPManager_UpdateGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LDAPManagerServer).UpdateGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ldapmanager.LDAPManager/UpdateGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LDAPManagerServer).UpdateGroup(ctx, req.(*UpdateGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LDAPManager_GetGroupList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LDAPManagerServer).GetGroupList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ldapmanager.LDAPManager/GetGroupList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LDAPManagerServer).GetGroupList(ctx, req.(*GetGroupListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LDAPManager_GetUserGroups_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserGroupsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LDAPManagerServer).GetUserGroups(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ldapmanager.LDAPManager/GetUserGroups",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LDAPManagerServer).GetUserGroups(ctx, req.(*GetUserGroupsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LDAPManager_IsGroupMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsGroupMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LDAPManagerServer).IsGroupMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ldapmanager.LDAPManager/IsGroupMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LDAPManagerServer).IsGroupMember(ctx, req.(*IsGroupMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LDAPManager_GetGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LDAPManagerServer).GetGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ldapmanager.LDAPManager/GetGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LDAPManagerServer).GetGroup(ctx, req.(*GetGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LDAPManager_AddGroupMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GroupMember)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LDAPManagerServer).AddGroupMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ldapmanager.LDAPManager/AddGroupMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LDAPManagerServer).AddGroupMember(ctx, req.(*GroupMember))
	}
	return interceptor(ctx, in, info, handler)
}

func _LDAPManager_DeleteGroupMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GroupMember)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LDAPManagerServer).DeleteGroupMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ldapmanager.LDAPManager/DeleteGroupMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LDAPManagerServer).DeleteGroupMember(ctx, req.(*GroupMember))
	}
	return interceptor(ctx, in, info, handler)
}

// LDAPManager_ServiceDesc is the grpc.ServiceDesc for LDAPManager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LDAPManager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ldapmanager.LDAPManager",
	HandlerType: (*LDAPManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _LDAPManager_Login_Handler,
		},
		{
			MethodName: "GetUserList",
			Handler:    _LDAPManager_GetUserList_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _LDAPManager_GetUser_Handler,
		},
		{
			MethodName: "NewUser",
			Handler:    _LDAPManager_NewUser_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _LDAPManager_UpdateUser_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _LDAPManager_DeleteUser_Handler,
		},
		{
			MethodName: "ChangePassword",
			Handler:    _LDAPManager_ChangePassword_Handler,
		},
		{
			MethodName: "NewGroup",
			Handler:    _LDAPManager_NewGroup_Handler,
		},
		{
			MethodName: "DeleteGroup",
			Handler:    _LDAPManager_DeleteGroup_Handler,
		},
		{
			MethodName: "UpdateGroup",
			Handler:    _LDAPManager_UpdateGroup_Handler,
		},
		{
			MethodName: "GetGroupList",
			Handler:    _LDAPManager_GetGroupList_Handler,
		},
		{
			MethodName: "GetUserGroups",
			Handler:    _LDAPManager_GetUserGroups_Handler,
		},
		{
			MethodName: "IsGroupMember",
			Handler:    _LDAPManager_IsGroupMember_Handler,
		},
		{
			MethodName: "GetGroup",
			Handler:    _LDAPManager_GetGroup_Handler,
		},
		{
			MethodName: "AddGroupMember",
			Handler:    _LDAPManager_AddGroupMember_Handler,
		},
		{
			MethodName: "DeleteGroupMember",
			Handler:    _LDAPManager_DeleteGroupMember_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ldap_manager.proto",
}
