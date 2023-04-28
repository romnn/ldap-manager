/* eslint-disable */
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

function createBaseEmpty(): Empty {
  return {};
}

export const Empty = {
  fromJSON(_: any): Empty {
    return {};
  },

  toJSON(_: Empty): unknown {
    const obj: any = {};
    return obj;
  },
};

function createBaseGetUserListRequest(): GetUserListRequest {
  return { start: 0, end: 0, sortOrder: 0, sortKey: "", filter: [] };
}

export const GetUserListRequest = {
  fromJSON(object: any): GetUserListRequest {
    return {
      start: isSet(object.start) ? Number(object.start) : 0,
      end: isSet(object.end) ? Number(object.end) : 0,
      sortOrder: isSet(object.sortOrder) ? sortOrderFromJSON(object.sortOrder) : 0,
      sortKey: isSet(object.sortKey) ? String(object.sortKey) : "",
      filter: Array.isArray(object?.filter)
        ? object.filter.map((e: any) => String(e))
        : [],
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
};

function createBaseUser(): User {
  return {
    username: "",
    firstName: "",
    lastName: "",
    displayName: "",
    email: "",
    loginShell: "",
    homeDirectory: "",
    CN: "",
    DN: "",
    UID: 0,
    GID: 0,
  };
}

export const User = {
  fromJSON(object: any): User {
    return {
      username: isSet(object.username) ? String(object.username) : "",
      firstName: isSet(object.firstName) ? String(object.firstName) : "",
      lastName: isSet(object.lastName) ? String(object.lastName) : "",
      displayName: isSet(object.displayName) ? String(object.displayName) : "",
      email: isSet(object.email) ? String(object.email) : "",
      loginShell: isSet(object.loginShell) ? String(object.loginShell) : "",
      homeDirectory: isSet(object.homeDirectory) ? String(object.homeDirectory) : "",
      CN: isSet(object.CN) ? String(object.CN) : "",
      DN: isSet(object.DN) ? String(object.DN) : "",
      UID: isSet(object.UID) ? Number(object.UID) : 0,
      GID: isSet(object.GID) ? Number(object.GID) : 0,
    };
  },

  toJSON(message: User): unknown {
    const obj: any = {};
    message.username !== undefined && (obj.username = message.username);
    message.firstName !== undefined && (obj.firstName = message.firstName);
    message.lastName !== undefined && (obj.lastName = message.lastName);
    message.displayName !== undefined && (obj.displayName = message.displayName);
    message.email !== undefined && (obj.email = message.email);
    message.loginShell !== undefined && (obj.loginShell = message.loginShell);
    message.homeDirectory !== undefined && (obj.homeDirectory = message.homeDirectory);
    message.CN !== undefined && (obj.CN = message.CN);
    message.DN !== undefined && (obj.DN = message.DN);
    message.UID !== undefined && (obj.UID = Math.round(message.UID));
    message.GID !== undefined && (obj.GID = Math.round(message.GID));
    return obj;
  },
};

function createBaseUserList(): UserList {
  return { users: [], total: 0 };
}

export const UserList = {
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
};

function createBaseAuthenticateUserRequest(): AuthenticateUserRequest {
  return { username: "", password: "" };
}

export const AuthenticateUserRequest = {
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
};

function createBaseGetUserRequest(): GetUserRequest {
  return { username: "" };
}

export const GetUserRequest = {
  fromJSON(object: any): GetUserRequest {
    return { username: isSet(object.username) ? String(object.username) : "" };
  },

  toJSON(message: GetUserRequest): unknown {
    const obj: any = {};
    message.username !== undefined && (obj.username = message.username);
    return obj;
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
};

function createBaseUpdateUserRequest(): UpdateUserRequest {
  return { username: "", update: undefined };
}

export const UpdateUserRequest = {
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
};

function createBaseDeleteUserRequest(): DeleteUserRequest {
  return { username: "" };
}

export const DeleteUserRequest = {
  fromJSON(object: any): DeleteUserRequest {
    return { username: isSet(object.username) ? String(object.username) : "" };
  },

  toJSON(message: DeleteUserRequest): unknown {
    const obj: any = {};
    message.username !== undefined && (obj.username = message.username);
    return obj;
  },
};

function createBaseNewGroupRequest(): NewGroupRequest {
  return { name: "", members: [] };
}

export const NewGroupRequest = {
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
};

function createBaseDeleteGroupRequest(): DeleteGroupRequest {
  return { name: "" };
}

export const DeleteGroupRequest = {
  fromJSON(object: any): DeleteGroupRequest {
    return { name: isSet(object.name) ? String(object.name) : "" };
  },

  toJSON(message: DeleteGroupRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    return obj;
  },
};

function createBaseUpdateGroupRequest(): UpdateGroupRequest {
  return { name: "", newName: "", GID: 0 };
}

export const UpdateGroupRequest = {
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
};

function createBaseGetGroupListRequest(): GetGroupListRequest {
  return { start: 0, end: 0, sortOrder: 0, sortKey: "", filter: [] };
}

export const GetGroupListRequest = {
  fromJSON(object: any): GetGroupListRequest {
    return {
      start: isSet(object.start) ? Number(object.start) : 0,
      end: isSet(object.end) ? Number(object.end) : 0,
      sortOrder: isSet(object.sortOrder) ? sortOrderFromJSON(object.sortOrder) : 0,
      sortKey: isSet(object.sortKey) ? String(object.sortKey) : "",
      filter: Array.isArray(object?.filter)
        ? object.filter.map((e: any) => String(e))
        : [],
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
};

function createBaseGroupList(): GroupList {
  return { groups: [], total: 0 };
}

export const GroupList = {
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
};

function createBaseIsGroupMemberRequest(): IsGroupMemberRequest {
  return { username: "", group: "" };
}

export const IsGroupMemberRequest = {
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
};

function createBaseGroupMemberStatus(): GroupMemberStatus {
  return { isMember: false };
}

export const GroupMemberStatus = {
  fromJSON(object: any): GroupMemberStatus {
    return { isMember: isSet(object.isMember) ? Boolean(object.isMember) : false };
  },

  toJSON(message: GroupMemberStatus): unknown {
    const obj: any = {};
    message.isMember !== undefined && (obj.isMember = message.isMember);
    return obj;
  },
};

function createBaseGetGroupRequest(): GetGroupRequest {
  return { start: 0, end: 0, sortOrder: 0, sortKey: "", name: "" };
}

export const GetGroupRequest = {
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
};

function createBaseGetUserGroupsRequest(): GetUserGroupsRequest {
  return { username: "" };
}

export const GetUserGroupsRequest = {
  fromJSON(object: any): GetUserGroupsRequest {
    return { username: isSet(object.username) ? String(object.username) : "" };
  },

  toJSON(message: GetUserGroupsRequest): unknown {
    const obj: any = {};
    message.username !== undefined && (obj.username = message.username);
    return obj;
  },
};

function createBaseGroup(): Group {
  return { name: "", members: [], GID: 0 };
}

export const Group = {
  fromJSON(object: any): Group {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      members: Array.isArray(object?.members) ? object.members.map((e: any) => GroupMember.fromJSON(e)) : [],
      GID: isSet(object.GID) ? Number(object.GID) : 0,
    };
  },

  toJSON(message: Group): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    if (message.members) {
      obj.members = message.members.map((e) => e ? GroupMember.toJSON(e) : undefined);
    } else {
      obj.members = [];
    }
    message.GID !== undefined && (obj.GID = Math.round(message.GID));
    return obj;
  },
};

function createBaseGroupMember(): GroupMember {
  return { group: "", username: "", dn: "" };
}

export const GroupMember = {
  fromJSON(object: any): GroupMember {
    return {
      group: isSet(object.group) ? String(object.group) : "",
      username: isSet(object.username) ? String(object.username) : "",
      dn: isSet(object.dn) ? String(object.dn) : "",
    };
  },

  toJSON(message: GroupMember): unknown {
    const obj: any = {};
    message.group !== undefined && (obj.group = message.group);
    message.username !== undefined && (obj.username = message.username);
    message.dn !== undefined && (obj.dn = message.dn);
    return obj;
  },
};

function createBaseChangePasswordRequest(): ChangePasswordRequest {
  return { username: "", password: "" };
}

export const ChangePasswordRequest = {
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
};

function createBaseLoginRequest(): LoginRequest {
  return { username: "", password: "" };
}

export const LoginRequest = {
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
};

function createBaseToken(): Token {
  return { token: "", username: "", UID: 0, displayName: "", isAdmin: false, expires: undefined };
}

export const Token = {
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
};

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

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
