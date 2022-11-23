import * as _m0 from "protobufjs/minimal";
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
    members: string[];
    /** repeated User members = 2; */
    GID: number;
}
export interface GroupMember {
    group: string;
    username: string;
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
    encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Empty;
    fromJSON(_: any): Empty;
    toJSON(_: Empty): unknown;
    fromPartial<I extends {} & {} & { [K in Exclude<keyof I, never>]: never; }>(_: I): Empty;
};
export declare const GetUserListRequest: {
    encode(message: GetUserListRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetUserListRequest;
    fromJSON(object: any): GetUserListRequest;
    toJSON(message: GetUserListRequest): unknown;
    fromPartial<I extends {
        start?: number | undefined;
        end?: number | undefined;
        sortOrder?: SortOrder | undefined;
        sortKey?: string | undefined;
        filter?: string[] | undefined;
    } & {
        start?: number | undefined;
        end?: number | undefined;
        sortOrder?: SortOrder | undefined;
        sortKey?: string | undefined;
        filter?: (string[] & string[] & { [K in Exclude<keyof I["filter"], keyof string[]>]: never; }) | undefined;
    } & { [K_1 in Exclude<keyof I, keyof GetUserListRequest>]: never; }>(object: I): GetUserListRequest;
};
export declare const User: {
    encode(message: User, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): User;
    fromJSON(object: any): User;
    toJSON(message: User): unknown;
    fromPartial<I extends {
        username?: string | undefined;
        firstName?: string | undefined;
        lastName?: string | undefined;
        displayName?: string | undefined;
        email?: string | undefined;
        loginShell?: string | undefined;
        homeDirectory?: string | undefined;
        CN?: string | undefined;
        DN?: string | undefined;
        UID?: number | undefined;
        GID?: number | undefined;
    } & {
        username?: string | undefined;
        firstName?: string | undefined;
        lastName?: string | undefined;
        displayName?: string | undefined;
        email?: string | undefined;
        loginShell?: string | undefined;
        homeDirectory?: string | undefined;
        CN?: string | undefined;
        DN?: string | undefined;
        UID?: number | undefined;
        GID?: number | undefined;
    } & { [K in Exclude<keyof I, keyof User>]: never; }>(object: I): User;
};
export declare const UserList: {
    encode(message: UserList, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): UserList;
    fromJSON(object: any): UserList;
    toJSON(message: UserList): unknown;
    fromPartial<I extends {
        users?: {
            username?: string | undefined;
            firstName?: string | undefined;
            lastName?: string | undefined;
            displayName?: string | undefined;
            email?: string | undefined;
            loginShell?: string | undefined;
            homeDirectory?: string | undefined;
            CN?: string | undefined;
            DN?: string | undefined;
            UID?: number | undefined;
            GID?: number | undefined;
        }[] | undefined;
        total?: number | undefined;
    } & {
        users?: ({
            username?: string | undefined;
            firstName?: string | undefined;
            lastName?: string | undefined;
            displayName?: string | undefined;
            email?: string | undefined;
            loginShell?: string | undefined;
            homeDirectory?: string | undefined;
            CN?: string | undefined;
            DN?: string | undefined;
            UID?: number | undefined;
            GID?: number | undefined;
        }[] & ({
            username?: string | undefined;
            firstName?: string | undefined;
            lastName?: string | undefined;
            displayName?: string | undefined;
            email?: string | undefined;
            loginShell?: string | undefined;
            homeDirectory?: string | undefined;
            CN?: string | undefined;
            DN?: string | undefined;
            UID?: number | undefined;
            GID?: number | undefined;
        } & {
            username?: string | undefined;
            firstName?: string | undefined;
            lastName?: string | undefined;
            displayName?: string | undefined;
            email?: string | undefined;
            loginShell?: string | undefined;
            homeDirectory?: string | undefined;
            CN?: string | undefined;
            DN?: string | undefined;
            UID?: number | undefined;
            GID?: number | undefined;
        } & { [K in Exclude<keyof I["users"][number], keyof User>]: never; })[] & { [K_1 in Exclude<keyof I["users"], keyof {
            username?: string | undefined;
            firstName?: string | undefined;
            lastName?: string | undefined;
            displayName?: string | undefined;
            email?: string | undefined;
            loginShell?: string | undefined;
            homeDirectory?: string | undefined;
            CN?: string | undefined;
            DN?: string | undefined;
            UID?: number | undefined;
            GID?: number | undefined;
        }[]>]: never; }) | undefined;
        total?: number | undefined;
    } & { [K_2 in Exclude<keyof I, keyof UserList>]: never; }>(object: I): UserList;
};
export declare const AuthenticateUserRequest: {
    encode(message: AuthenticateUserRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): AuthenticateUserRequest;
    fromJSON(object: any): AuthenticateUserRequest;
    toJSON(message: AuthenticateUserRequest): unknown;
    fromPartial<I extends {
        username?: string | undefined;
        password?: string | undefined;
    } & {
        username?: string | undefined;
        password?: string | undefined;
    } & { [K in Exclude<keyof I, keyof AuthenticateUserRequest>]: never; }>(object: I): AuthenticateUserRequest;
};
export declare const GetUserRequest: {
    encode(message: GetUserRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetUserRequest;
    fromJSON(object: any): GetUserRequest;
    toJSON(message: GetUserRequest): unknown;
    fromPartial<I extends {
        username?: string | undefined;
    } & {
        username?: string | undefined;
    } & { [K in Exclude<keyof I, "username">]: never; }>(object: I): GetUserRequest;
};
export declare const NewUserRequest: {
    encode(message: NewUserRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): NewUserRequest;
    fromJSON(object: any): NewUserRequest;
    toJSON(message: NewUserRequest): unknown;
    fromPartial<I extends {
        firstName?: string | undefined;
        lastName?: string | undefined;
        UID?: number | undefined;
        GID?: number | undefined;
        loginShell?: string | undefined;
        homeDirectory?: string | undefined;
        username?: string | undefined;
        email?: string | undefined;
        password?: string | undefined;
    } & {
        firstName?: string | undefined;
        lastName?: string | undefined;
        UID?: number | undefined;
        GID?: number | undefined;
        loginShell?: string | undefined;
        homeDirectory?: string | undefined;
        username?: string | undefined;
        email?: string | undefined;
        password?: string | undefined;
    } & { [K in Exclude<keyof I, keyof NewUserRequest>]: never; }>(object: I): NewUserRequest;
};
export declare const UpdateUserRequest: {
    encode(message: UpdateUserRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): UpdateUserRequest;
    fromJSON(object: any): UpdateUserRequest;
    toJSON(message: UpdateUserRequest): unknown;
    fromPartial<I extends {
        username?: string | undefined;
        update?: {
            firstName?: string | undefined;
            lastName?: string | undefined;
            UID?: number | undefined;
            GID?: number | undefined;
            loginShell?: string | undefined;
            homeDirectory?: string | undefined;
            username?: string | undefined;
            email?: string | undefined;
            password?: string | undefined;
        } | undefined;
    } & {
        username?: string | undefined;
        update?: ({
            firstName?: string | undefined;
            lastName?: string | undefined;
            UID?: number | undefined;
            GID?: number | undefined;
            loginShell?: string | undefined;
            homeDirectory?: string | undefined;
            username?: string | undefined;
            email?: string | undefined;
            password?: string | undefined;
        } & {
            firstName?: string | undefined;
            lastName?: string | undefined;
            UID?: number | undefined;
            GID?: number | undefined;
            loginShell?: string | undefined;
            homeDirectory?: string | undefined;
            username?: string | undefined;
            email?: string | undefined;
            password?: string | undefined;
        } & { [K in Exclude<keyof I["update"], keyof NewUserRequest>]: never; }) | undefined;
    } & { [K_1 in Exclude<keyof I, keyof UpdateUserRequest>]: never; }>(object: I): UpdateUserRequest;
};
export declare const DeleteUserRequest: {
    encode(message: DeleteUserRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): DeleteUserRequest;
    fromJSON(object: any): DeleteUserRequest;
    toJSON(message: DeleteUserRequest): unknown;
    fromPartial<I extends {
        username?: string | undefined;
    } & {
        username?: string | undefined;
    } & { [K in Exclude<keyof I, "username">]: never; }>(object: I): DeleteUserRequest;
};
export declare const NewGroupRequest: {
    encode(message: NewGroupRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): NewGroupRequest;
    fromJSON(object: any): NewGroupRequest;
    toJSON(message: NewGroupRequest): unknown;
    fromPartial<I extends {
        name?: string | undefined;
        members?: string[] | undefined;
    } & {
        name?: string | undefined;
        members?: (string[] & string[] & { [K in Exclude<keyof I["members"], keyof string[]>]: never; }) | undefined;
    } & { [K_1 in Exclude<keyof I, keyof NewGroupRequest>]: never; }>(object: I): NewGroupRequest;
};
export declare const DeleteGroupRequest: {
    encode(message: DeleteGroupRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): DeleteGroupRequest;
    fromJSON(object: any): DeleteGroupRequest;
    toJSON(message: DeleteGroupRequest): unknown;
    fromPartial<I extends {
        name?: string | undefined;
    } & {
        name?: string | undefined;
    } & { [K in Exclude<keyof I, "name">]: never; }>(object: I): DeleteGroupRequest;
};
export declare const UpdateGroupRequest: {
    encode(message: UpdateGroupRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): UpdateGroupRequest;
    fromJSON(object: any): UpdateGroupRequest;
    toJSON(message: UpdateGroupRequest): unknown;
    fromPartial<I extends {
        name?: string | undefined;
        newName?: string | undefined;
        GID?: number | undefined;
    } & {
        name?: string | undefined;
        newName?: string | undefined;
        GID?: number | undefined;
    } & { [K in Exclude<keyof I, keyof UpdateGroupRequest>]: never; }>(object: I): UpdateGroupRequest;
};
export declare const GetGroupListRequest: {
    encode(message: GetGroupListRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetGroupListRequest;
    fromJSON(object: any): GetGroupListRequest;
    toJSON(message: GetGroupListRequest): unknown;
    fromPartial<I extends {
        start?: number | undefined;
        end?: number | undefined;
        sortOrder?: SortOrder | undefined;
        sortKey?: string | undefined;
        filter?: string[] | undefined;
    } & {
        start?: number | undefined;
        end?: number | undefined;
        sortOrder?: SortOrder | undefined;
        sortKey?: string | undefined;
        filter?: (string[] & string[] & { [K in Exclude<keyof I["filter"], keyof string[]>]: never; }) | undefined;
    } & { [K_1 in Exclude<keyof I, keyof GetGroupListRequest>]: never; }>(object: I): GetGroupListRequest;
};
export declare const GroupList: {
    encode(message: GroupList, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GroupList;
    fromJSON(object: any): GroupList;
    toJSON(message: GroupList): unknown;
    fromPartial<I extends {
        groups?: {
            name?: string | undefined;
            members?: string[] | undefined;
            GID?: number | undefined;
        }[] | undefined;
        total?: number | undefined;
    } & {
        groups?: ({
            name?: string | undefined;
            members?: string[] | undefined;
            GID?: number | undefined;
        }[] & ({
            name?: string | undefined;
            members?: string[] | undefined;
            GID?: number | undefined;
        } & {
            name?: string | undefined;
            members?: (string[] & string[] & { [K in Exclude<keyof I["groups"][number]["members"], keyof string[]>]: never; }) | undefined;
            GID?: number | undefined;
        } & { [K_1 in Exclude<keyof I["groups"][number], keyof Group>]: never; })[] & { [K_2 in Exclude<keyof I["groups"], keyof {
            name?: string | undefined;
            members?: string[] | undefined;
            GID?: number | undefined;
        }[]>]: never; }) | undefined;
        total?: number | undefined;
    } & { [K_3 in Exclude<keyof I, keyof GroupList>]: never; }>(object: I): GroupList;
};
export declare const IsGroupMemberRequest: {
    encode(message: IsGroupMemberRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): IsGroupMemberRequest;
    fromJSON(object: any): IsGroupMemberRequest;
    toJSON(message: IsGroupMemberRequest): unknown;
    fromPartial<I extends {
        username?: string | undefined;
        group?: string | undefined;
    } & {
        username?: string | undefined;
        group?: string | undefined;
    } & { [K in Exclude<keyof I, keyof IsGroupMemberRequest>]: never; }>(object: I): IsGroupMemberRequest;
};
export declare const GroupMemberStatus: {
    encode(message: GroupMemberStatus, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GroupMemberStatus;
    fromJSON(object: any): GroupMemberStatus;
    toJSON(message: GroupMemberStatus): unknown;
    fromPartial<I extends {
        isMember?: boolean | undefined;
    } & {
        isMember?: boolean | undefined;
    } & { [K in Exclude<keyof I, "isMember">]: never; }>(object: I): GroupMemberStatus;
};
export declare const GetGroupRequest: {
    encode(message: GetGroupRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetGroupRequest;
    fromJSON(object: any): GetGroupRequest;
    toJSON(message: GetGroupRequest): unknown;
    fromPartial<I extends {
        start?: number | undefined;
        end?: number | undefined;
        sortOrder?: SortOrder | undefined;
        sortKey?: string | undefined;
        name?: string | undefined;
    } & {
        start?: number | undefined;
        end?: number | undefined;
        sortOrder?: SortOrder | undefined;
        sortKey?: string | undefined;
        name?: string | undefined;
    } & { [K in Exclude<keyof I, keyof GetGroupRequest>]: never; }>(object: I): GetGroupRequest;
};
export declare const GetUserGroupsRequest: {
    encode(message: GetUserGroupsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetUserGroupsRequest;
    fromJSON(object: any): GetUserGroupsRequest;
    toJSON(message: GetUserGroupsRequest): unknown;
    fromPartial<I extends {
        username?: string | undefined;
    } & {
        username?: string | undefined;
    } & { [K in Exclude<keyof I, "username">]: never; }>(object: I): GetUserGroupsRequest;
};
export declare const Group: {
    encode(message: Group, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Group;
    fromJSON(object: any): Group;
    toJSON(message: Group): unknown;
    fromPartial<I extends {
        name?: string | undefined;
        members?: string[] | undefined;
        GID?: number | undefined;
    } & {
        name?: string | undefined;
        members?: (string[] & string[] & { [K in Exclude<keyof I["members"], keyof string[]>]: never; }) | undefined;
        GID?: number | undefined;
    } & { [K_1 in Exclude<keyof I, keyof Group>]: never; }>(object: I): Group;
};
export declare const GroupMember: {
    encode(message: GroupMember, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GroupMember;
    fromJSON(object: any): GroupMember;
    toJSON(message: GroupMember): unknown;
    fromPartial<I extends {
        group?: string | undefined;
        username?: string | undefined;
    } & {
        group?: string | undefined;
        username?: string | undefined;
    } & { [K in Exclude<keyof I, keyof GroupMember>]: never; }>(object: I): GroupMember;
};
export declare const ChangePasswordRequest: {
    encode(message: ChangePasswordRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ChangePasswordRequest;
    fromJSON(object: any): ChangePasswordRequest;
    toJSON(message: ChangePasswordRequest): unknown;
    fromPartial<I extends {
        username?: string | undefined;
        password?: string | undefined;
    } & {
        username?: string | undefined;
        password?: string | undefined;
    } & { [K in Exclude<keyof I, keyof ChangePasswordRequest>]: never; }>(object: I): ChangePasswordRequest;
};
export declare const LoginRequest: {
    encode(message: LoginRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): LoginRequest;
    fromJSON(object: any): LoginRequest;
    toJSON(message: LoginRequest): unknown;
    fromPartial<I extends {
        username?: string | undefined;
        password?: string | undefined;
    } & {
        username?: string | undefined;
        password?: string | undefined;
    } & { [K in Exclude<keyof I, keyof LoginRequest>]: never; }>(object: I): LoginRequest;
};
export declare const Token: {
    encode(message: Token, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Token;
    fromJSON(object: any): Token;
    toJSON(message: Token): unknown;
    fromPartial<I extends {
        token?: string | undefined;
        username?: string | undefined;
        UID?: number | undefined;
        displayName?: string | undefined;
        isAdmin?: boolean | undefined;
        expires?: Date | undefined;
    } & {
        token?: string | undefined;
        username?: string | undefined;
        UID?: number | undefined;
        displayName?: string | undefined;
        isAdmin?: boolean | undefined;
        expires?: Date | undefined;
    } & { [K in Exclude<keyof I, keyof Token>]: never; }>(object: I): Token;
};
export interface LDAPManager {
    /** Authentication */
    Login(request: LoginRequest): Promise<Token>;
    /** Users */
    GetUserList(request: GetUserListRequest): Promise<UserList>;
    GetUser(request: GetUserRequest): Promise<User>;
    NewUser(request: NewUserRequest): Promise<Empty>;
    UpdateUser(request: UpdateUserRequest): Promise<Token>;
    DeleteUser(request: DeleteUserRequest): Promise<Empty>;
    ChangePassword(request: ChangePasswordRequest): Promise<Empty>;
    /** Groups */
    NewGroup(request: NewGroupRequest): Promise<Empty>;
    DeleteGroup(request: DeleteGroupRequest): Promise<Empty>;
    UpdateGroup(request: UpdateGroupRequest): Promise<Empty>;
    GetGroupList(request: GetGroupListRequest): Promise<GroupList>;
    GetUserGroups(request: GetUserGroupsRequest): Promise<GroupList>;
    /** Group members */
    IsGroupMember(request: IsGroupMemberRequest): Promise<GroupMemberStatus>;
    GetGroup(request: GetGroupRequest): Promise<Group>;
    AddGroupMember(request: GroupMember): Promise<Empty>;
    RemoveGroupMember(request: GroupMember): Promise<Empty>;
}
export declare class LDAPManagerClientImpl implements LDAPManager {
    private readonly rpc;
    private readonly service;
    constructor(rpc: Rpc, opts?: {
        service?: string;
    });
    Login(request: LoginRequest): Promise<Token>;
    GetUserList(request: GetUserListRequest): Promise<UserList>;
    GetUser(request: GetUserRequest): Promise<User>;
    NewUser(request: NewUserRequest): Promise<Empty>;
    UpdateUser(request: UpdateUserRequest): Promise<Token>;
    DeleteUser(request: DeleteUserRequest): Promise<Empty>;
    ChangePassword(request: ChangePasswordRequest): Promise<Empty>;
    NewGroup(request: NewGroupRequest): Promise<Empty>;
    DeleteGroup(request: DeleteGroupRequest): Promise<Empty>;
    UpdateGroup(request: UpdateGroupRequest): Promise<Empty>;
    GetGroupList(request: GetGroupListRequest): Promise<GroupList>;
    GetUserGroups(request: GetUserGroupsRequest): Promise<GroupList>;
    IsGroupMember(request: IsGroupMemberRequest): Promise<GroupMemberStatus>;
    GetGroup(request: GetGroupRequest): Promise<Group>;
    AddGroupMember(request: GroupMember): Promise<Empty>;
    RemoveGroupMember(request: GroupMember): Promise<Empty>;
}
interface Rpc {
    request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}
type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;
export type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P : P & {
    [K in keyof P]: Exact<P[K], I[K]>;
} & {
    [K in Exclude<keyof I, KeysOfUnion<P>>]: never;
};
export {};
