"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Token = exports.LoginRequest = exports.ChangePasswordRequest = exports.GroupMember = exports.Group = exports.GetUserGroupsRequest = exports.GetGroupRequest = exports.GroupMemberStatus = exports.IsGroupMemberRequest = exports.GroupList = exports.GetGroupListRequest = exports.UpdateGroupRequest = exports.DeleteGroupRequest = exports.NewGroupRequest = exports.DeleteUserRequest = exports.UpdateUserRequest = exports.NewUserRequest = exports.GetUserRequest = exports.AuthenticateUserRequest = exports.UserList = exports.User = exports.GetUserListRequest = exports.Empty = exports.sortOrderToJSON = exports.sortOrderFromJSON = exports.SortOrder = exports.protobufPackage = void 0;
/* eslint-disable */
var timestamp_1 = require("./google/protobuf/timestamp");
exports.protobufPackage = "ldapmanager";
var SortOrder;
(function (SortOrder) {
    SortOrder[SortOrder["ASCENDING"] = 0] = "ASCENDING";
    SortOrder[SortOrder["DESCENDING"] = 1] = "DESCENDING";
    SortOrder[SortOrder["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(SortOrder = exports.SortOrder || (exports.SortOrder = {}));
function sortOrderFromJSON(object) {
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
exports.sortOrderFromJSON = sortOrderFromJSON;
function sortOrderToJSON(object) {
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
exports.sortOrderToJSON = sortOrderToJSON;
function createBaseEmpty() {
    return {};
}
exports.Empty = {
    fromJSON: function (_) {
        return {};
    },
    toJSON: function (_) {
        var obj = {};
        return obj;
    },
};
function createBaseGetUserListRequest() {
    return { start: 0, end: 0, sortOrder: 0, sortKey: "", filter: [] };
}
exports.GetUserListRequest = {
    fromJSON: function (object) {
        return {
            start: isSet(object.start) ? Number(object.start) : 0,
            end: isSet(object.end) ? Number(object.end) : 0,
            sortOrder: isSet(object.sortOrder) ? sortOrderFromJSON(object.sortOrder) : 0,
            sortKey: isSet(object.sortKey) ? String(object.sortKey) : "",
            filter: Array.isArray(object === null || object === void 0 ? void 0 : object.filter)
                ? object.filter.map(function (e) { return String(e); })
                : [],
        };
    },
    toJSON: function (message) {
        var obj = {};
        message.start !== undefined && (obj.start = Math.round(message.start));
        message.end !== undefined && (obj.end = Math.round(message.end));
        message.sortOrder !== undefined && (obj.sortOrder = sortOrderToJSON(message.sortOrder));
        message.sortKey !== undefined && (obj.sortKey = message.sortKey);
        if (message.filter) {
            obj.filter = message.filter.map(function (e) { return e; });
        }
        else {
            obj.filter = [];
        }
        return obj;
    },
};
function createBaseUser() {
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
exports.User = {
    fromJSON: function (object) {
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
    toJSON: function (message) {
        var obj = {};
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
function createBaseUserList() {
    return { users: [], total: 0 };
}
exports.UserList = {
    fromJSON: function (object) {
        return {
            users: Array.isArray(object === null || object === void 0 ? void 0 : object.users) ? object.users.map(function (e) { return exports.User.fromJSON(e); }) : [],
            total: isSet(object.total) ? Number(object.total) : 0,
        };
    },
    toJSON: function (message) {
        var obj = {};
        if (message.users) {
            obj.users = message.users.map(function (e) { return e ? exports.User.toJSON(e) : undefined; });
        }
        else {
            obj.users = [];
        }
        message.total !== undefined && (obj.total = Math.round(message.total));
        return obj;
    },
};
function createBaseAuthenticateUserRequest() {
    return { username: "", password: "" };
}
exports.AuthenticateUserRequest = {
    fromJSON: function (object) {
        return {
            username: isSet(object.username) ? String(object.username) : "",
            password: isSet(object.password) ? String(object.password) : "",
        };
    },
    toJSON: function (message) {
        var obj = {};
        message.username !== undefined && (obj.username = message.username);
        message.password !== undefined && (obj.password = message.password);
        return obj;
    },
};
function createBaseGetUserRequest() {
    return { username: "" };
}
exports.GetUserRequest = {
    fromJSON: function (object) {
        return { username: isSet(object.username) ? String(object.username) : "" };
    },
    toJSON: function (message) {
        var obj = {};
        message.username !== undefined && (obj.username = message.username);
        return obj;
    },
};
function createBaseNewUserRequest() {
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
exports.NewUserRequest = {
    fromJSON: function (object) {
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
    toJSON: function (message) {
        var obj = {};
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
function createBaseUpdateUserRequest() {
    return { username: "", update: undefined };
}
exports.UpdateUserRequest = {
    fromJSON: function (object) {
        return {
            username: isSet(object.username) ? String(object.username) : "",
            update: isSet(object.update) ? exports.NewUserRequest.fromJSON(object.update) : undefined,
        };
    },
    toJSON: function (message) {
        var obj = {};
        message.username !== undefined && (obj.username = message.username);
        message.update !== undefined && (obj.update = message.update ? exports.NewUserRequest.toJSON(message.update) : undefined);
        return obj;
    },
};
function createBaseDeleteUserRequest() {
    return { username: "" };
}
exports.DeleteUserRequest = {
    fromJSON: function (object) {
        return { username: isSet(object.username) ? String(object.username) : "" };
    },
    toJSON: function (message) {
        var obj = {};
        message.username !== undefined && (obj.username = message.username);
        return obj;
    },
};
function createBaseNewGroupRequest() {
    return { name: "", members: [] };
}
exports.NewGroupRequest = {
    fromJSON: function (object) {
        return {
            name: isSet(object.name) ? String(object.name) : "",
            members: Array.isArray(object === null || object === void 0 ? void 0 : object.members) ? object.members.map(function (e) { return String(e); }) : [],
        };
    },
    toJSON: function (message) {
        var obj = {};
        message.name !== undefined && (obj.name = message.name);
        if (message.members) {
            obj.members = message.members.map(function (e) { return e; });
        }
        else {
            obj.members = [];
        }
        return obj;
    },
};
function createBaseDeleteGroupRequest() {
    return { name: "" };
}
exports.DeleteGroupRequest = {
    fromJSON: function (object) {
        return { name: isSet(object.name) ? String(object.name) : "" };
    },
    toJSON: function (message) {
        var obj = {};
        message.name !== undefined && (obj.name = message.name);
        return obj;
    },
};
function createBaseUpdateGroupRequest() {
    return { name: "", newName: "", GID: 0 };
}
exports.UpdateGroupRequest = {
    fromJSON: function (object) {
        return {
            name: isSet(object.name) ? String(object.name) : "",
            newName: isSet(object.newName) ? String(object.newName) : "",
            GID: isSet(object.GID) ? Number(object.GID) : 0,
        };
    },
    toJSON: function (message) {
        var obj = {};
        message.name !== undefined && (obj.name = message.name);
        message.newName !== undefined && (obj.newName = message.newName);
        message.GID !== undefined && (obj.GID = Math.round(message.GID));
        return obj;
    },
};
function createBaseGetGroupListRequest() {
    return { start: 0, end: 0, sortOrder: 0, sortKey: "", filter: [] };
}
exports.GetGroupListRequest = {
    fromJSON: function (object) {
        return {
            start: isSet(object.start) ? Number(object.start) : 0,
            end: isSet(object.end) ? Number(object.end) : 0,
            sortOrder: isSet(object.sortOrder) ? sortOrderFromJSON(object.sortOrder) : 0,
            sortKey: isSet(object.sortKey) ? String(object.sortKey) : "",
            filter: Array.isArray(object === null || object === void 0 ? void 0 : object.filter)
                ? object.filter.map(function (e) { return String(e); })
                : [],
        };
    },
    toJSON: function (message) {
        var obj = {};
        message.start !== undefined && (obj.start = Math.round(message.start));
        message.end !== undefined && (obj.end = Math.round(message.end));
        message.sortOrder !== undefined && (obj.sortOrder = sortOrderToJSON(message.sortOrder));
        message.sortKey !== undefined && (obj.sortKey = message.sortKey);
        if (message.filter) {
            obj.filter = message.filter.map(function (e) { return e; });
        }
        else {
            obj.filter = [];
        }
        return obj;
    },
};
function createBaseGroupList() {
    return { groups: [], total: 0 };
}
exports.GroupList = {
    fromJSON: function (object) {
        return {
            groups: Array.isArray(object === null || object === void 0 ? void 0 : object.groups) ? object.groups.map(function (e) { return exports.Group.fromJSON(e); }) : [],
            total: isSet(object.total) ? Number(object.total) : 0,
        };
    },
    toJSON: function (message) {
        var obj = {};
        if (message.groups) {
            obj.groups = message.groups.map(function (e) { return e ? exports.Group.toJSON(e) : undefined; });
        }
        else {
            obj.groups = [];
        }
        message.total !== undefined && (obj.total = Math.round(message.total));
        return obj;
    },
};
function createBaseIsGroupMemberRequest() {
    return { username: "", group: "" };
}
exports.IsGroupMemberRequest = {
    fromJSON: function (object) {
        return {
            username: isSet(object.username) ? String(object.username) : "",
            group: isSet(object.group) ? String(object.group) : "",
        };
    },
    toJSON: function (message) {
        var obj = {};
        message.username !== undefined && (obj.username = message.username);
        message.group !== undefined && (obj.group = message.group);
        return obj;
    },
};
function createBaseGroupMemberStatus() {
    return { isMember: false };
}
exports.GroupMemberStatus = {
    fromJSON: function (object) {
        return { isMember: isSet(object.isMember) ? Boolean(object.isMember) : false };
    },
    toJSON: function (message) {
        var obj = {};
        message.isMember !== undefined && (obj.isMember = message.isMember);
        return obj;
    },
};
function createBaseGetGroupRequest() {
    return { start: 0, end: 0, sortOrder: 0, sortKey: "", name: "" };
}
exports.GetGroupRequest = {
    fromJSON: function (object) {
        return {
            start: isSet(object.start) ? Number(object.start) : 0,
            end: isSet(object.end) ? Number(object.end) : 0,
            sortOrder: isSet(object.sortOrder) ? sortOrderFromJSON(object.sortOrder) : 0,
            sortKey: isSet(object.sortKey) ? String(object.sortKey) : "",
            name: isSet(object.name) ? String(object.name) : "",
        };
    },
    toJSON: function (message) {
        var obj = {};
        message.start !== undefined && (obj.start = Math.round(message.start));
        message.end !== undefined && (obj.end = Math.round(message.end));
        message.sortOrder !== undefined && (obj.sortOrder = sortOrderToJSON(message.sortOrder));
        message.sortKey !== undefined && (obj.sortKey = message.sortKey);
        message.name !== undefined && (obj.name = message.name);
        return obj;
    },
};
function createBaseGetUserGroupsRequest() {
    return { username: "" };
}
exports.GetUserGroupsRequest = {
    fromJSON: function (object) {
        return { username: isSet(object.username) ? String(object.username) : "" };
    },
    toJSON: function (message) {
        var obj = {};
        message.username !== undefined && (obj.username = message.username);
        return obj;
    },
};
function createBaseGroup() {
    return { name: "", members: [], GID: 0 };
}
exports.Group = {
    fromJSON: function (object) {
        return {
            name: isSet(object.name) ? String(object.name) : "",
            members: Array.isArray(object === null || object === void 0 ? void 0 : object.members) ? object.members.map(function (e) { return String(e); }) : [],
            GID: isSet(object.GID) ? Number(object.GID) : 0,
        };
    },
    toJSON: function (message) {
        var obj = {};
        message.name !== undefined && (obj.name = message.name);
        if (message.members) {
            obj.members = message.members.map(function (e) { return e; });
        }
        else {
            obj.members = [];
        }
        message.GID !== undefined && (obj.GID = Math.round(message.GID));
        return obj;
    },
};
function createBaseGroupMember() {
    return { group: "", username: "" };
}
exports.GroupMember = {
    fromJSON: function (object) {
        return {
            group: isSet(object.group) ? String(object.group) : "",
            username: isSet(object.username) ? String(object.username) : "",
        };
    },
    toJSON: function (message) {
        var obj = {};
        message.group !== undefined && (obj.group = message.group);
        message.username !== undefined && (obj.username = message.username);
        return obj;
    },
};
function createBaseChangePasswordRequest() {
    return { username: "", password: "" };
}
exports.ChangePasswordRequest = {
    fromJSON: function (object) {
        return {
            username: isSet(object.username) ? String(object.username) : "",
            password: isSet(object.password) ? String(object.password) : "",
        };
    },
    toJSON: function (message) {
        var obj = {};
        message.username !== undefined && (obj.username = message.username);
        message.password !== undefined && (obj.password = message.password);
        return obj;
    },
};
function createBaseLoginRequest() {
    return { username: "", password: "" };
}
exports.LoginRequest = {
    fromJSON: function (object) {
        return {
            username: isSet(object.username) ? String(object.username) : "",
            password: isSet(object.password) ? String(object.password) : "",
        };
    },
    toJSON: function (message) {
        var obj = {};
        message.username !== undefined && (obj.username = message.username);
        message.password !== undefined && (obj.password = message.password);
        return obj;
    },
};
function createBaseToken() {
    return { token: "", username: "", UID: 0, displayName: "", isAdmin: false, expires: undefined };
}
exports.Token = {
    fromJSON: function (object) {
        return {
            token: isSet(object.token) ? String(object.token) : "",
            username: isSet(object.username) ? String(object.username) : "",
            UID: isSet(object.UID) ? Number(object.UID) : 0,
            displayName: isSet(object.displayName) ? String(object.displayName) : "",
            isAdmin: isSet(object.isAdmin) ? Boolean(object.isAdmin) : false,
            expires: isSet(object.expires) ? fromJsonTimestamp(object.expires) : undefined,
        };
    },
    toJSON: function (message) {
        var obj = {};
        message.token !== undefined && (obj.token = message.token);
        message.username !== undefined && (obj.username = message.username);
        message.UID !== undefined && (obj.UID = Math.round(message.UID));
        message.displayName !== undefined && (obj.displayName = message.displayName);
        message.isAdmin !== undefined && (obj.isAdmin = message.isAdmin);
        message.expires !== undefined && (obj.expires = message.expires.toISOString());
        return obj;
    },
};
function fromTimestamp(t) {
    var millis = t.seconds * 1000;
    millis += t.nanos / 1000000;
    return new Date(millis);
}
function fromJsonTimestamp(o) {
    if (o instanceof Date) {
        return o;
    }
    else if (typeof o === "string") {
        return new Date(o);
    }
    else {
        return fromTimestamp(timestamp_1.Timestamp.fromJSON(o));
    }
}
function isSet(value) {
    return value !== null && value !== undefined;
}
//# sourceMappingURL=ldap_manager.js.map