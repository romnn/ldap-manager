export declare const protobufPackage = "ldapmanager";
export declare enum SortOrder {
    ASCENDING = 0,
    DESCENDING = 1,
    UNRECOGNIZED = -1
}
export declare function sortOrderFromJSON(object: any): SortOrder;
export declare function sortOrderToJSON(object: SortOrder): string;
export interface Empty {
}
export interface GetUserListRequest {
    start: number;
    end: number;
    sortOrder: SortOrder;
    sortKey: string;
    filter: string[];
}
export interface User {
    username: string;
    firstName: string;
    lastName: string;
    displayName: string;
    email: string;
    loginShell: string;
    homeDirectory: string;
    CN: string;
    DN: string;
    UID: number;
    GID: number;
}
export interface UserList {
    users: User[];
    total: number;
}
export interface AuthenticateUserRequest {
    username: string;
    password: string;
}
export interface GetUserRequest {
    username: string;
}
export interface NewUserRequest {
    firstName: string;
    lastName: string;
    UID: number;
    GID: number;
    loginShell: string;
    homeDirectory: string;
    username: string;
    email: string;
    password: string;
}
export interface UpdateUserRequest {
    username: string;
    update: NewUserRequest | undefined;
}
export interface DeleteUserRequest {
    username: string;
}
export interface NewGroupRequest {
    name: string;
    members: string[];
}
export interface DeleteGroupRequest {
    name: string;
}
export interface UpdateGroupRequest {
    name: string;
    newName: string;
    GID: number;
}
export interface GetGroupListRequest {
    start: number;
    end: number;
    sortOrder: SortOrder;
    sortKey: string;
    filter: string[];
}
export interface GroupList {
    groups: Group[];
    total: number;
}
export interface IsGroupMemberRequest {
    username: string;
    group: string;
}
export interface GroupMemberStatus {
    isMember: boolean;
}
export interface GetGroupRequest {
    start: number;
    end: number;
    sortOrder: SortOrder;
    sortKey: string;
    name: string;
}
export interface GetUserGroupsRequest {
    username: string;
}
export interface Group {
    name: string;
    members: GroupMember[];
    GID: number;
}
export interface GroupMember {
    group: string;
    username: string;
    dn: string;
}
export interface ChangePasswordRequest {
    username: string;
    password: string;
}
export interface LoginRequest {
    username: string;
    password: string;
}
export interface Token {
    token: string;
    username: string;
    UID: number;
    displayName: string;
    isAdmin: boolean;
    expires: Date | undefined;
}
export declare const Empty: {
    fromJSON(_: any): Empty;
    toJSON(_: Empty): unknown;
};
export declare const GetUserListRequest: {
    fromJSON(object: any): GetUserListRequest;
    toJSON(message: GetUserListRequest): unknown;
};
export declare const User: {
    fromJSON(object: any): User;
    toJSON(message: User): unknown;
};
export declare const UserList: {
    fromJSON(object: any): UserList;
    toJSON(message: UserList): unknown;
};
export declare const AuthenticateUserRequest: {
    fromJSON(object: any): AuthenticateUserRequest;
    toJSON(message: AuthenticateUserRequest): unknown;
};
export declare const GetUserRequest: {
    fromJSON(object: any): GetUserRequest;
    toJSON(message: GetUserRequest): unknown;
};
export declare const NewUserRequest: {
    fromJSON(object: any): NewUserRequest;
    toJSON(message: NewUserRequest): unknown;
};
export declare const UpdateUserRequest: {
    fromJSON(object: any): UpdateUserRequest;
    toJSON(message: UpdateUserRequest): unknown;
};
export declare const DeleteUserRequest: {
    fromJSON(object: any): DeleteUserRequest;
    toJSON(message: DeleteUserRequest): unknown;
};
export declare const NewGroupRequest: {
    fromJSON(object: any): NewGroupRequest;
    toJSON(message: NewGroupRequest): unknown;
};
export declare const DeleteGroupRequest: {
    fromJSON(object: any): DeleteGroupRequest;
    toJSON(message: DeleteGroupRequest): unknown;
};
export declare const UpdateGroupRequest: {
    fromJSON(object: any): UpdateGroupRequest;
    toJSON(message: UpdateGroupRequest): unknown;
};
export declare const GetGroupListRequest: {
    fromJSON(object: any): GetGroupListRequest;
    toJSON(message: GetGroupListRequest): unknown;
};
export declare const GroupList: {
    fromJSON(object: any): GroupList;
    toJSON(message: GroupList): unknown;
};
export declare const IsGroupMemberRequest: {
    fromJSON(object: any): IsGroupMemberRequest;
    toJSON(message: IsGroupMemberRequest): unknown;
};
export declare const GroupMemberStatus: {
    fromJSON(object: any): GroupMemberStatus;
    toJSON(message: GroupMemberStatus): unknown;
};
export declare const GetGroupRequest: {
    fromJSON(object: any): GetGroupRequest;
    toJSON(message: GetGroupRequest): unknown;
};
export declare const GetUserGroupsRequest: {
    fromJSON(object: any): GetUserGroupsRequest;
    toJSON(message: GetUserGroupsRequest): unknown;
};
export declare const Group: {
    fromJSON(object: any): Group;
    toJSON(message: Group): unknown;
};
export declare const GroupMember: {
    fromJSON(object: any): GroupMember;
    toJSON(message: GroupMember): unknown;
};
export declare const ChangePasswordRequest: {
    fromJSON(object: any): ChangePasswordRequest;
    toJSON(message: ChangePasswordRequest): unknown;
};
export declare const LoginRequest: {
    fromJSON(object: any): LoginRequest;
    toJSON(message: LoginRequest): unknown;
};
export declare const Token: {
    fromJSON(object: any): Token;
    toJSON(message: Token): unknown;
};
