/* eslint-disable */
import * as Long from "long";
import * as _m0 from "protobufjs/minimal";
import { Timestamp } from "./google/protobuf/timestamp";

export const protobufPackage = "ldapmanager";

export enum SortOrder {
  ASCENDING = 0,
  DESCENDING = 1,
  UNRECOGNIZED = -1,
}

export function sortOrderFromJSON(object: any): SortOrder {
  switch (object) {
    case 0:
    case "ASCENDING":
      return SortOrder.ASCENDING;
    case 1:
    case "DESCENDING":
      return SortOrder.DESCENDING;
    case -1:
    case "UNRECOGNIZED":
    default:
      return SortOrder.UNRECOGNIZED;
  }
}

export function sortOrderToJSON(object: SortOrder): string {
  switch (object) {
    case SortOrder.ASCENDING:
      return "ASCENDING";
    case SortOrder.DESCENDING:
      return "DESCENDING";
    case SortOrder.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

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
  CN: string;
  email: string;
  UID: number;
  GID: number;
  loginShell: string;
  homeDirectory: string;
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

function createBaseEmpty(): Empty {
  return {};
}

export const Empty = {
  encode(_: Empty, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Empty {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEmpty();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): Empty {
    return {};
  },

  toJSON(_: Empty): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Empty>, I>>(_: I): Empty {
    const message = createBaseEmpty();
    return message;
  },
};

function createBaseGetUserListRequest(): GetUserListRequest {
  return { start: 0, end: 0, sortOrder: 0, sortKey: "", filter: [] };
}

export const GetUserListRequest = {
  encode(message: GetUserListRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.start !== 0) {
      writer.uint32(8).int32(message.start);
    }
    if (message.end !== 0) {
      writer.uint32(16).int32(message.end);
    }
    if (message.sortOrder !== 0) {
      writer.uint32(24).int32(message.sortOrder);
    }
    if (message.sortKey !== "") {
      writer.uint32(34).string(message.sortKey);
    }
    for (const v of message.filter) {
      writer.uint32(82).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetUserListRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetUserListRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.start = reader.int32();
          break;
        case 2:
          message.end = reader.int32();
          break;
        case 3:
          message.sortOrder = reader.int32() as any;
          break;
        case 4:
          message.sortKey = reader.string();
          break;
        case 10:
          message.filter.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetUserListRequest {
    return {
      start: isSet(object.start) ? Number(object.start) : 0,
      end: isSet(object.end) ? Number(object.end) : 0,
      sortOrder: isSet(object.sortOrder) ? sortOrderFromJSON(object.sortOrder) : 0,
      sortKey: isSet(object.sortKey) ? String(object.sortKey) : "",
      filter: Array.isArray(object?.filter) ? object.filter.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: GetUserListRequest): unknown {
    const obj: any = {};
    message.start !== undefined && (obj.start = Math.round(message.start));
    message.end !== undefined && (obj.end = Math.round(message.end));
    message.sortOrder !== undefined && (obj.sortOrder = sortOrderToJSON(message.sortOrder));
    message.sortKey !== undefined && (obj.sortKey = message.sortKey);
    if (message.filter) {
      obj.filter = message.filter.map((e) => e);
    } else {
      obj.filter = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetUserListRequest>, I>>(object: I): GetUserListRequest {
    const message = createBaseGetUserListRequest();
    message.start = object.start ?? 0;
    message.end = object.end ?? 0;
    message.sortOrder = object.sortOrder ?? 0;
    message.sortKey = object.sortKey ?? "";
    message.filter = object.filter?.map((e) => e) || [];
    return message;
  },
};

function createBaseUser(): User {
  return {
    username: "",
    firstName: "",
    lastName: "",
    displayName: "",
    CN: "",
    email: "",
    UID: 0,
    GID: 0,
    loginShell: "",
    homeDirectory: "",
  };
}

export const User = {
  encode(message: User, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.username !== "") {
      writer.uint32(10).string(message.username);
    }
    if (message.firstName !== "") {
      writer.uint32(82).string(message.firstName);
    }
    if (message.lastName !== "") {
      writer.uint32(90).string(message.lastName);
    }
    if (message.displayName !== "") {
      writer.uint32(98).string(message.displayName);
    }
    if (message.CN !== "") {
      writer.uint32(106).string(message.CN);
    }
    if (message.email !== "") {
      writer.uint32(114).string(message.email);
    }
    if (message.UID !== 0) {
      writer.uint32(160).int32(message.UID);
    }
    if (message.GID !== 0) {
      writer.uint32(168).int64(message.GID);
    }
    if (message.loginShell !== "") {
      writer.uint32(242).string(message.loginShell);
    }
    if (message.homeDirectory !== "") {
      writer.uint32(250).string(message.homeDirectory);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): User {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUser();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.username = reader.string();
          break;
        case 10:
          message.firstName = reader.string();
          break;
        case 11:
          message.lastName = reader.string();
          break;
        case 12:
          message.displayName = reader.string();
          break;
        case 13:
          message.CN = reader.string();
          break;
        case 14:
          message.email = reader.string();
          break;
        case 20:
          message.UID = reader.int32();
          break;
        case 21:
          message.GID = longToNumber(reader.int64() as Long);
          break;
        case 30:
          message.loginShell = reader.string();
          break;
        case 31:
          message.homeDirectory = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): User {
    return {
      username: isSet(object.username) ? String(object.username) : "",
      firstName: isSet(object.firstName) ? String(object.firstName) : "",
      lastName: isSet(object.lastName) ? String(object.lastName) : "",
      displayName: isSet(object.displayName) ? String(object.displayName) : "",
      CN: isSet(object.CN) ? String(object.CN) : "",
      email: isSet(object.email) ? String(object.email) : "",
      UID: isSet(object.UID) ? Number(object.UID) : 0,
      GID: isSet(object.GID) ? Number(object.GID) : 0,
      loginShell: isSet(object.loginShell) ? String(object.loginShell) : "",
      homeDirectory: isSet(object.homeDirectory) ? String(object.homeDirectory) : "",
    };
  },

  toJSON(message: User): unknown {
    const obj: any = {};
    message.username !== undefined && (obj.username = message.username);
    message.firstName !== undefined && (obj.firstName = message.firstName);
    message.lastName !== undefined && (obj.lastName = message.lastName);
    message.displayName !== undefined && (obj.displayName = message.displayName);
    message.CN !== undefined && (obj.CN = message.CN);
    message.email !== undefined && (obj.email = message.email);
    message.UID !== undefined && (obj.UID = Math.round(message.UID));
    message.GID !== undefined && (obj.GID = Math.round(message.GID));
    message.loginShell !== undefined && (obj.loginShell = message.loginShell);
    message.homeDirectory !== undefined && (obj.homeDirectory = message.homeDirectory);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<User>, I>>(object: I): User {
    const message = createBaseUser();
    message.username = object.username ?? "";
    message.firstName = object.firstName ?? "";
    message.lastName = object.lastName ?? "";
    message.displayName = object.displayName ?? "";
    message.CN = object.CN ?? "";
    message.email = object.email ?? "";
    message.UID = object.UID ?? 0;
    message.GID = object.GID ?? 0;
    message.loginShell = object.loginShell ?? "";
    message.homeDirectory = object.homeDirectory ?? "";
    return message;
  },
};

function createBaseUserList(): UserList {
  return { users: [], total: 0 };
}

export const UserList = {
  encode(message: UserList, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.users) {
      User.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.total !== 0) {
      writer.uint32(80).int64(message.total);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserList {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserList();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.users.push(User.decode(reader, reader.uint32()));
          break;
        case 10:
          message.total = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserList {
    return {
      users: Array.isArray(object?.users) ? object.users.map((e: any) => User.fromJSON(e)) : [],
      total: isSet(object.total) ? Number(object.total) : 0,
    };
  },

  toJSON(message: UserList): unknown {
    const obj: any = {};
    if (message.users) {
      obj.users = message.users.map((e) => e ? User.toJSON(e) : undefined);
    } else {
      obj.users = [];
    }
    message.total !== undefined && (obj.total = Math.round(message.total));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserList>, I>>(object: I): UserList {
    const message = createBaseUserList();
    message.users = object.users?.map((e) => User.fromPartial(e)) || [];
    message.total = object.total ?? 0;
    return message;
  },
};

function createBaseAuthenticateUserRequest(): AuthenticateUserRequest {
  return { username: "", password: "" };
}

export const AuthenticateUserRequest = {
  encode(message: AuthenticateUserRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.username !== "") {
      writer.uint32(10).string(message.username);
    }
    if (message.password !== "") {
      writer.uint32(18).string(message.password);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AuthenticateUserRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAuthenticateUserRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.username = reader.string();
          break;
        case 2:
          message.password = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AuthenticateUserRequest {
    return {
      username: isSet(object.username) ? String(object.username) : "",
      password: isSet(object.password) ? String(object.password) : "",
    };
  },

  toJSON(message: AuthenticateUserRequest): unknown {
    const obj: any = {};
    message.username !== undefined && (obj.username = message.username);
    message.password !== undefined && (obj.password = message.password);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<AuthenticateUserRequest>, I>>(object: I): AuthenticateUserRequest {
    const message = createBaseAuthenticateUserRequest();
    message.username = object.username ?? "";
    message.password = object.password ?? "";
    return message;
  },
};

function createBaseGetUserRequest(): GetUserRequest {
  return { username: "" };
}

export const GetUserRequest = {
  encode(message: GetUserRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.username !== "") {
      writer.uint32(10).string(message.username);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetUserRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetUserRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.username = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetUserRequest {
    return { username: isSet(object.username) ? String(object.username) : "" };
  },

  toJSON(message: GetUserRequest): unknown {
    const obj: any = {};
    message.username !== undefined && (obj.username = message.username);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetUserRequest>, I>>(object: I): GetUserRequest {
    const message = createBaseGetUserRequest();
    message.username = object.username ?? "";
    return message;
  },
};

function createBaseNewUserRequest(): NewUserRequest {
  return {
    firstName: "",
    lastName: "",
    UID: 0,
    GID: 0,
    loginShell: "",
    homeDirectory: "",
    username: "",
    email: "",
    password: "",
  };
}

export const NewUserRequest = {
  encode(message: NewUserRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.firstName !== "") {
      writer.uint32(10).string(message.firstName);
    }
    if (message.lastName !== "") {
      writer.uint32(18).string(message.lastName);
    }
    if (message.UID !== 0) {
      writer.uint32(80).int64(message.UID);
    }
    if (message.GID !== 0) {
      writer.uint32(88).int64(message.GID);
    }
    if (message.loginShell !== "") {
      writer.uint32(98).string(message.loginShell);
    }
    if (message.homeDirectory !== "") {
      writer.uint32(106).string(message.homeDirectory);
    }
    if (message.username !== "") {
      writer.uint32(162).string(message.username);
    }
    if (message.email !== "") {
      writer.uint32(170).string(message.email);
    }
    if (message.password !== "") {
      writer.uint32(178).string(message.password);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): NewUserRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseNewUserRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.firstName = reader.string();
          break;
        case 2:
          message.lastName = reader.string();
          break;
        case 10:
          message.UID = longToNumber(reader.int64() as Long);
          break;
        case 11:
          message.GID = longToNumber(reader.int64() as Long);
          break;
        case 12:
          message.loginShell = reader.string();
          break;
        case 13:
          message.homeDirectory = reader.string();
          break;
        case 20:
          message.username = reader.string();
          break;
        case 21:
          message.email = reader.string();
          break;
        case 22:
          message.password = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): NewUserRequest {
    return {
      firstName: isSet(object.firstName) ? String(object.firstName) : "",
      lastName: isSet(object.lastName) ? String(object.lastName) : "",
      UID: isSet(object.UID) ? Number(object.UID) : 0,
      GID: isSet(object.GID) ? Number(object.GID) : 0,
      loginShell: isSet(object.loginShell) ? String(object.loginShell) : "",
      homeDirectory: isSet(object.homeDirectory) ? String(object.homeDirectory) : "",
      username: isSet(object.username) ? String(object.username) : "",
      email: isSet(object.email) ? String(object.email) : "",
      password: isSet(object.password) ? String(object.password) : "",
    };
  },

  toJSON(message: NewUserRequest): unknown {
    const obj: any = {};
    message.firstName !== undefined && (obj.firstName = message.firstName);
    message.lastName !== undefined && (obj.lastName = message.lastName);
    message.UID !== undefined && (obj.UID = Math.round(message.UID));
    message.GID !== undefined && (obj.GID = Math.round(message.GID));
    message.loginShell !== undefined && (obj.loginShell = message.loginShell);
    message.homeDirectory !== undefined && (obj.homeDirectory = message.homeDirectory);
    message.username !== undefined && (obj.username = message.username);
    message.email !== undefined && (obj.email = message.email);
    message.password !== undefined && (obj.password = message.password);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<NewUserRequest>, I>>(object: I): NewUserRequest {
    const message = createBaseNewUserRequest();
    message.firstName = object.firstName ?? "";
    message.lastName = object.lastName ?? "";
    message.UID = object.UID ?? 0;
    message.GID = object.GID ?? 0;
    message.loginShell = object.loginShell ?? "";
    message.homeDirectory = object.homeDirectory ?? "";
    message.username = object.username ?? "";
    message.email = object.email ?? "";
    message.password = object.password ?? "";
    return message;
  },
};

function createBaseUpdateUserRequest(): UpdateUserRequest {
  return { username: "", update: undefined };
}

export const UpdateUserRequest = {
  encode(message: UpdateUserRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.username !== "") {
      writer.uint32(10).string(message.username);
    }
    if (message.update !== undefined) {
      NewUserRequest.encode(message.update, writer.uint32(82).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateUserRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateUserRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.username = reader.string();
          break;
        case 10:
          message.update = NewUserRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateUserRequest {
    return {
      username: isSet(object.username) ? String(object.username) : "",
      update: isSet(object.update) ? NewUserRequest.fromJSON(object.update) : undefined,
    };
  },

  toJSON(message: UpdateUserRequest): unknown {
    const obj: any = {};
    message.username !== undefined && (obj.username = message.username);
    message.update !== undefined && (obj.update = message.update ? NewUserRequest.toJSON(message.update) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UpdateUserRequest>, I>>(object: I): UpdateUserRequest {
    const message = createBaseUpdateUserRequest();
    message.username = object.username ?? "";
    message.update = (object.update !== undefined && object.update !== null)
      ? NewUserRequest.fromPartial(object.update)
      : undefined;
    return message;
  },
};

function createBaseDeleteUserRequest(): DeleteUserRequest {
  return { username: "" };
}

export const DeleteUserRequest = {
  encode(message: DeleteUserRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.username !== "") {
      writer.uint32(10).string(message.username);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteUserRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteUserRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.username = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DeleteUserRequest {
    return { username: isSet(object.username) ? String(object.username) : "" };
  },

  toJSON(message: DeleteUserRequest): unknown {
    const obj: any = {};
    message.username !== undefined && (obj.username = message.username);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DeleteUserRequest>, I>>(object: I): DeleteUserRequest {
    const message = createBaseDeleteUserRequest();
    message.username = object.username ?? "";
    return message;
  },
};

function createBaseNewGroupRequest(): NewGroupRequest {
  return { name: "", members: [] };
}

export const NewGroupRequest = {
  encode(message: NewGroupRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    for (const v of message.members) {
      writer.uint32(18).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): NewGroupRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseNewGroupRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.members.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): NewGroupRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      members: Array.isArray(object?.members) ? object.members.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: NewGroupRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    if (message.members) {
      obj.members = message.members.map((e) => e);
    } else {
      obj.members = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<NewGroupRequest>, I>>(object: I): NewGroupRequest {
    const message = createBaseNewGroupRequest();
    message.name = object.name ?? "";
    message.members = object.members?.map((e) => e) || [];
    return message;
  },
};

function createBaseDeleteGroupRequest(): DeleteGroupRequest {
  return { name: "" };
}

export const DeleteGroupRequest = {
  encode(message: DeleteGroupRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteGroupRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteGroupRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DeleteGroupRequest {
    return { name: isSet(object.name) ? String(object.name) : "" };
  },

  toJSON(message: DeleteGroupRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DeleteGroupRequest>, I>>(object: I): DeleteGroupRequest {
    const message = createBaseDeleteGroupRequest();
    message.name = object.name ?? "";
    return message;
  },
};

function createBaseUpdateGroupRequest(): UpdateGroupRequest {
  return { name: "", newName: "", GID: 0 };
}

export const UpdateGroupRequest = {
  encode(message: UpdateGroupRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.newName !== "") {
      writer.uint32(18).string(message.newName);
    }
    if (message.GID !== 0) {
      writer.uint32(24).int64(message.GID);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateGroupRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateGroupRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.newName = reader.string();
          break;
        case 3:
          message.GID = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateGroupRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      newName: isSet(object.newName) ? String(object.newName) : "",
      GID: isSet(object.GID) ? Number(object.GID) : 0,
    };
  },

  toJSON(message: UpdateGroupRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.newName !== undefined && (obj.newName = message.newName);
    message.GID !== undefined && (obj.GID = Math.round(message.GID));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UpdateGroupRequest>, I>>(object: I): UpdateGroupRequest {
    const message = createBaseUpdateGroupRequest();
    message.name = object.name ?? "";
    message.newName = object.newName ?? "";
    message.GID = object.GID ?? 0;
    return message;
  },
};

function createBaseGetGroupListRequest(): GetGroupListRequest {
  return { start: 0, end: 0, sortOrder: 0, sortKey: "", filter: [] };
}

export const GetGroupListRequest = {
  encode(message: GetGroupListRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.start !== 0) {
      writer.uint32(8).int32(message.start);
    }
    if (message.end !== 0) {
      writer.uint32(16).int32(message.end);
    }
    if (message.sortOrder !== 0) {
      writer.uint32(24).int32(message.sortOrder);
    }
    if (message.sortKey !== "") {
      writer.uint32(34).string(message.sortKey);
    }
    for (const v of message.filter) {
      writer.uint32(82).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetGroupListRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetGroupListRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.start = reader.int32();
          break;
        case 2:
          message.end = reader.int32();
          break;
        case 3:
          message.sortOrder = reader.int32() as any;
          break;
        case 4:
          message.sortKey = reader.string();
          break;
        case 10:
          message.filter.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetGroupListRequest {
    return {
      start: isSet(object.start) ? Number(object.start) : 0,
      end: isSet(object.end) ? Number(object.end) : 0,
      sortOrder: isSet(object.sortOrder) ? sortOrderFromJSON(object.sortOrder) : 0,
      sortKey: isSet(object.sortKey) ? String(object.sortKey) : "",
      filter: Array.isArray(object?.filter) ? object.filter.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: GetGroupListRequest): unknown {
    const obj: any = {};
    message.start !== undefined && (obj.start = Math.round(message.start));
    message.end !== undefined && (obj.end = Math.round(message.end));
    message.sortOrder !== undefined && (obj.sortOrder = sortOrderToJSON(message.sortOrder));
    message.sortKey !== undefined && (obj.sortKey = message.sortKey);
    if (message.filter) {
      obj.filter = message.filter.map((e) => e);
    } else {
      obj.filter = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetGroupListRequest>, I>>(object: I): GetGroupListRequest {
    const message = createBaseGetGroupListRequest();
    message.start = object.start ?? 0;
    message.end = object.end ?? 0;
    message.sortOrder = object.sortOrder ?? 0;
    message.sortKey = object.sortKey ?? "";
    message.filter = object.filter?.map((e) => e) || [];
    return message;
  },
};

function createBaseGroupList(): GroupList {
  return { groups: [], total: 0 };
}

export const GroupList = {
  encode(message: GroupList, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.groups) {
      Group.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.total !== 0) {
      writer.uint32(80).int64(message.total);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GroupList {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGroupList();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.groups.push(Group.decode(reader, reader.uint32()));
          break;
        case 10:
          message.total = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GroupList {
    return {
      groups: Array.isArray(object?.groups) ? object.groups.map((e: any) => Group.fromJSON(e)) : [],
      total: isSet(object.total) ? Number(object.total) : 0,
    };
  },

  toJSON(message: GroupList): unknown {
    const obj: any = {};
    if (message.groups) {
      obj.groups = message.groups.map((e) => e ? Group.toJSON(e) : undefined);
    } else {
      obj.groups = [];
    }
    message.total !== undefined && (obj.total = Math.round(message.total));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GroupList>, I>>(object: I): GroupList {
    const message = createBaseGroupList();
    message.groups = object.groups?.map((e) => Group.fromPartial(e)) || [];
    message.total = object.total ?? 0;
    return message;
  },
};

function createBaseIsGroupMemberRequest(): IsGroupMemberRequest {
  return { username: "", group: "" };
}

export const IsGroupMemberRequest = {
  encode(message: IsGroupMemberRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.username !== "") {
      writer.uint32(10).string(message.username);
    }
    if (message.group !== "") {
      writer.uint32(18).string(message.group);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IsGroupMemberRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIsGroupMemberRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.username = reader.string();
          break;
        case 2:
          message.group = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): IsGroupMemberRequest {
    return {
      username: isSet(object.username) ? String(object.username) : "",
      group: isSet(object.group) ? String(object.group) : "",
    };
  },

  toJSON(message: IsGroupMemberRequest): unknown {
    const obj: any = {};
    message.username !== undefined && (obj.username = message.username);
    message.group !== undefined && (obj.group = message.group);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<IsGroupMemberRequest>, I>>(object: I): IsGroupMemberRequest {
    const message = createBaseIsGroupMemberRequest();
    message.username = object.username ?? "";
    message.group = object.group ?? "";
    return message;
  },
};

function createBaseGroupMemberStatus(): GroupMemberStatus {
  return { isMember: false };
}

export const GroupMemberStatus = {
  encode(message: GroupMemberStatus, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.isMember === true) {
      writer.uint32(8).bool(message.isMember);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GroupMemberStatus {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGroupMemberStatus();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.isMember = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GroupMemberStatus {
    return { isMember: isSet(object.isMember) ? Boolean(object.isMember) : false };
  },

  toJSON(message: GroupMemberStatus): unknown {
    const obj: any = {};
    message.isMember !== undefined && (obj.isMember = message.isMember);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GroupMemberStatus>, I>>(object: I): GroupMemberStatus {
    const message = createBaseGroupMemberStatus();
    message.isMember = object.isMember ?? false;
    return message;
  },
};

function createBaseGetGroupRequest(): GetGroupRequest {
  return { start: 0, end: 0, sortOrder: 0, sortKey: "", name: "" };
}

export const GetGroupRequest = {
  encode(message: GetGroupRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.start !== 0) {
      writer.uint32(8).int32(message.start);
    }
    if (message.end !== 0) {
      writer.uint32(16).int32(message.end);
    }
    if (message.sortOrder !== 0) {
      writer.uint32(24).int32(message.sortOrder);
    }
    if (message.sortKey !== "") {
      writer.uint32(34).string(message.sortKey);
    }
    if (message.name !== "") {
      writer.uint32(82).string(message.name);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetGroupRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetGroupRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.start = reader.int32();
          break;
        case 2:
          message.end = reader.int32();
          break;
        case 3:
          message.sortOrder = reader.int32() as any;
          break;
        case 4:
          message.sortKey = reader.string();
          break;
        case 10:
          message.name = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetGroupRequest {
    return {
      start: isSet(object.start) ? Number(object.start) : 0,
      end: isSet(object.end) ? Number(object.end) : 0,
      sortOrder: isSet(object.sortOrder) ? sortOrderFromJSON(object.sortOrder) : 0,
      sortKey: isSet(object.sortKey) ? String(object.sortKey) : "",
      name: isSet(object.name) ? String(object.name) : "",
    };
  },

  toJSON(message: GetGroupRequest): unknown {
    const obj: any = {};
    message.start !== undefined && (obj.start = Math.round(message.start));
    message.end !== undefined && (obj.end = Math.round(message.end));
    message.sortOrder !== undefined && (obj.sortOrder = sortOrderToJSON(message.sortOrder));
    message.sortKey !== undefined && (obj.sortKey = message.sortKey);
    message.name !== undefined && (obj.name = message.name);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetGroupRequest>, I>>(object: I): GetGroupRequest {
    const message = createBaseGetGroupRequest();
    message.start = object.start ?? 0;
    message.end = object.end ?? 0;
    message.sortOrder = object.sortOrder ?? 0;
    message.sortKey = object.sortKey ?? "";
    message.name = object.name ?? "";
    return message;
  },
};

function createBaseGetUserGroupsRequest(): GetUserGroupsRequest {
  return { username: "" };
}

export const GetUserGroupsRequest = {
  encode(message: GetUserGroupsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.username !== "") {
      writer.uint32(10).string(message.username);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetUserGroupsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetUserGroupsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.username = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetUserGroupsRequest {
    return { username: isSet(object.username) ? String(object.username) : "" };
  },

  toJSON(message: GetUserGroupsRequest): unknown {
    const obj: any = {};
    message.username !== undefined && (obj.username = message.username);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetUserGroupsRequest>, I>>(object: I): GetUserGroupsRequest {
    const message = createBaseGetUserGroupsRequest();
    message.username = object.username ?? "";
    return message;
  },
};

function createBaseGroup(): Group {
  return { name: "", members: [], GID: 0 };
}

export const Group = {
  encode(message: Group, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    for (const v of message.members) {
      writer.uint32(18).string(v!);
    }
    if (message.GID !== 0) {
      writer.uint32(24).int64(message.GID);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Group {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGroup();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.members.push(reader.string());
          break;
        case 3:
          message.GID = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Group {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      members: Array.isArray(object?.members) ? object.members.map((e: any) => String(e)) : [],
      GID: isSet(object.GID) ? Number(object.GID) : 0,
    };
  },

  toJSON(message: Group): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    if (message.members) {
      obj.members = message.members.map((e) => e);
    } else {
      obj.members = [];
    }
    message.GID !== undefined && (obj.GID = Math.round(message.GID));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Group>, I>>(object: I): Group {
    const message = createBaseGroup();
    message.name = object.name ?? "";
    message.members = object.members?.map((e) => e) || [];
    message.GID = object.GID ?? 0;
    return message;
  },
};

function createBaseGroupMember(): GroupMember {
  return { group: "", username: "" };
}

export const GroupMember = {
  encode(message: GroupMember, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.group !== "") {
      writer.uint32(10).string(message.group);
    }
    if (message.username !== "") {
      writer.uint32(18).string(message.username);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GroupMember {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGroupMember();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.group = reader.string();
          break;
        case 2:
          message.username = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GroupMember {
    return {
      group: isSet(object.group) ? String(object.group) : "",
      username: isSet(object.username) ? String(object.username) : "",
    };
  },

  toJSON(message: GroupMember): unknown {
    const obj: any = {};
    message.group !== undefined && (obj.group = message.group);
    message.username !== undefined && (obj.username = message.username);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GroupMember>, I>>(object: I): GroupMember {
    const message = createBaseGroupMember();
    message.group = object.group ?? "";
    message.username = object.username ?? "";
    return message;
  },
};

function createBaseChangePasswordRequest(): ChangePasswordRequest {
  return { username: "", password: "" };
}

export const ChangePasswordRequest = {
  encode(message: ChangePasswordRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.username !== "") {
      writer.uint32(10).string(message.username);
    }
    if (message.password !== "") {
      writer.uint32(18).string(message.password);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ChangePasswordRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChangePasswordRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.username = reader.string();
          break;
        case 2:
          message.password = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ChangePasswordRequest {
    return {
      username: isSet(object.username) ? String(object.username) : "",
      password: isSet(object.password) ? String(object.password) : "",
    };
  },

  toJSON(message: ChangePasswordRequest): unknown {
    const obj: any = {};
    message.username !== undefined && (obj.username = message.username);
    message.password !== undefined && (obj.password = message.password);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ChangePasswordRequest>, I>>(object: I): ChangePasswordRequest {
    const message = createBaseChangePasswordRequest();
    message.username = object.username ?? "";
    message.password = object.password ?? "";
    return message;
  },
};

function createBaseLoginRequest(): LoginRequest {
  return { username: "", password: "" };
}

export const LoginRequest = {
  encode(message: LoginRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.username !== "") {
      writer.uint32(10).string(message.username);
    }
    if (message.password !== "") {
      writer.uint32(18).string(message.password);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LoginRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLoginRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.username = reader.string();
          break;
        case 2:
          message.password = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): LoginRequest {
    return {
      username: isSet(object.username) ? String(object.username) : "",
      password: isSet(object.password) ? String(object.password) : "",
    };
  },

  toJSON(message: LoginRequest): unknown {
    const obj: any = {};
    message.username !== undefined && (obj.username = message.username);
    message.password !== undefined && (obj.password = message.password);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<LoginRequest>, I>>(object: I): LoginRequest {
    const message = createBaseLoginRequest();
    message.username = object.username ?? "";
    message.password = object.password ?? "";
    return message;
  },
};

function createBaseToken(): Token {
  return { token: "", username: "", UID: 0, displayName: "", isAdmin: false, expires: undefined };
}

export const Token = {
  encode(message: Token, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.token !== "") {
      writer.uint32(10).string(message.token);
    }
    if (message.username !== "") {
      writer.uint32(18).string(message.username);
    }
    if (message.UID !== 0) {
      writer.uint32(24).int32(message.UID);
    }
    if (message.displayName !== "") {
      writer.uint32(34).string(message.displayName);
    }
    if (message.isAdmin === true) {
      writer.uint32(40).bool(message.isAdmin);
    }
    if (message.expires !== undefined) {
      Timestamp.encode(toTimestamp(message.expires), writer.uint32(82).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Token {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseToken();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.token = reader.string();
          break;
        case 2:
          message.username = reader.string();
          break;
        case 3:
          message.UID = reader.int32();
          break;
        case 4:
          message.displayName = reader.string();
          break;
        case 5:
          message.isAdmin = reader.bool();
          break;
        case 10:
          message.expires = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Token {
    return {
      token: isSet(object.token) ? String(object.token) : "",
      username: isSet(object.username) ? String(object.username) : "",
      UID: isSet(object.UID) ? Number(object.UID) : 0,
      displayName: isSet(object.displayName) ? String(object.displayName) : "",
      isAdmin: isSet(object.isAdmin) ? Boolean(object.isAdmin) : false,
      expires: isSet(object.expires) ? fromJsonTimestamp(object.expires) : undefined,
    };
  },

  toJSON(message: Token): unknown {
    const obj: any = {};
    message.token !== undefined && (obj.token = message.token);
    message.username !== undefined && (obj.username = message.username);
    message.UID !== undefined && (obj.UID = Math.round(message.UID));
    message.displayName !== undefined && (obj.displayName = message.displayName);
    message.isAdmin !== undefined && (obj.isAdmin = message.isAdmin);
    message.expires !== undefined && (obj.expires = message.expires.toISOString());
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Token>, I>>(object: I): Token {
    const message = createBaseToken();
    message.token = object.token ?? "";
    message.username = object.username ?? "";
    message.UID = object.UID ?? 0;
    message.displayName = object.displayName ?? "";
    message.isAdmin = object.isAdmin ?? false;
    message.expires = object.expires ?? undefined;
    return message;
  },
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

export class LDAPManagerClientImpl implements LDAPManager {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || "ldapmanager.LDAPManager";
    this.rpc = rpc;
    this.Login = this.Login.bind(this);
    this.GetUserList = this.GetUserList.bind(this);
    this.GetUser = this.GetUser.bind(this);
    this.NewUser = this.NewUser.bind(this);
    this.UpdateUser = this.UpdateUser.bind(this);
    this.DeleteUser = this.DeleteUser.bind(this);
    this.ChangePassword = this.ChangePassword.bind(this);
    this.NewGroup = this.NewGroup.bind(this);
    this.DeleteGroup = this.DeleteGroup.bind(this);
    this.UpdateGroup = this.UpdateGroup.bind(this);
    this.GetGroupList = this.GetGroupList.bind(this);
    this.GetUserGroups = this.GetUserGroups.bind(this);
    this.IsGroupMember = this.IsGroupMember.bind(this);
    this.GetGroup = this.GetGroup.bind(this);
    this.AddGroupMember = this.AddGroupMember.bind(this);
    this.RemoveGroupMember = this.RemoveGroupMember.bind(this);
  }
  Login(request: LoginRequest): Promise<Token> {
    const data = LoginRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Login", data);
    return promise.then((data) => Token.decode(new _m0.Reader(data)));
  }

  GetUserList(request: GetUserListRequest): Promise<UserList> {
    const data = GetUserListRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetUserList", data);
    return promise.then((data) => UserList.decode(new _m0.Reader(data)));
  }

  GetUser(request: GetUserRequest): Promise<User> {
    const data = GetUserRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetUser", data);
    return promise.then((data) => User.decode(new _m0.Reader(data)));
  }

  NewUser(request: NewUserRequest): Promise<Empty> {
    const data = NewUserRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "NewUser", data);
    return promise.then((data) => Empty.decode(new _m0.Reader(data)));
  }

  UpdateUser(request: UpdateUserRequest): Promise<Token> {
    const data = UpdateUserRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "UpdateUser", data);
    return promise.then((data) => Token.decode(new _m0.Reader(data)));
  }

  DeleteUser(request: DeleteUserRequest): Promise<Empty> {
    const data = DeleteUserRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "DeleteUser", data);
    return promise.then((data) => Empty.decode(new _m0.Reader(data)));
  }

  ChangePassword(request: ChangePasswordRequest): Promise<Empty> {
    const data = ChangePasswordRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "ChangePassword", data);
    return promise.then((data) => Empty.decode(new _m0.Reader(data)));
  }

  NewGroup(request: NewGroupRequest): Promise<Empty> {
    const data = NewGroupRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "NewGroup", data);
    return promise.then((data) => Empty.decode(new _m0.Reader(data)));
  }

  DeleteGroup(request: DeleteGroupRequest): Promise<Empty> {
    const data = DeleteGroupRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "DeleteGroup", data);
    return promise.then((data) => Empty.decode(new _m0.Reader(data)));
  }

  UpdateGroup(request: UpdateGroupRequest): Promise<Empty> {
    const data = UpdateGroupRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "UpdateGroup", data);
    return promise.then((data) => Empty.decode(new _m0.Reader(data)));
  }

  GetGroupList(request: GetGroupListRequest): Promise<GroupList> {
    const data = GetGroupListRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetGroupList", data);
    return promise.then((data) => GroupList.decode(new _m0.Reader(data)));
  }

  GetUserGroups(request: GetUserGroupsRequest): Promise<GroupList> {
    const data = GetUserGroupsRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetUserGroups", data);
    return promise.then((data) => GroupList.decode(new _m0.Reader(data)));
  }

  IsGroupMember(request: IsGroupMemberRequest): Promise<GroupMemberStatus> {
    const data = IsGroupMemberRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "IsGroupMember", data);
    return promise.then((data) => GroupMemberStatus.decode(new _m0.Reader(data)));
  }

  GetGroup(request: GetGroupRequest): Promise<Group> {
    const data = GetGroupRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetGroup", data);
    return promise.then((data) => Group.decode(new _m0.Reader(data)));
  }

  AddGroupMember(request: GroupMember): Promise<Empty> {
    const data = GroupMember.encode(request).finish();
    const promise = this.rpc.request(this.service, "AddGroupMember", data);
    return promise.then((data) => Empty.decode(new _m0.Reader(data)));
  }

  RemoveGroupMember(request: GroupMember): Promise<Empty> {
    const data = GroupMember.encode(request).finish();
    const promise = this.rpc.request(this.service, "RemoveGroupMember", data);
    return promise.then((data) => Empty.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function toTimestamp(date: Date): Timestamp {
  const seconds = date.getTime() / 1_000;
  const nanos = (date.getTime() % 1_000) * 1_000_000;
  return { seconds, nanos };
}

function fromTimestamp(t: Timestamp): Date {
  let millis = t.seconds * 1_000;
  millis += t.nanos / 1_000_000;
  return new Date(millis);
}

function fromJsonTimestamp(o: any): Date {
  if (o instanceof Date) {
    return o;
  } else if (typeof o === "string") {
    return new Date(o);
  } else {
    return fromTimestamp(Timestamp.fromJSON(o));
  }
}

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

// If you get a compile-error about 'Constructor<Long> and ... have no overlap',
// add '--ts_proto_opt=esModuleInterop=true' as a flag when calling 'protoc'.
if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
