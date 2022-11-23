"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LDAPManagerClientImpl = exports.Token = exports.LoginRequest = exports.ChangePasswordRequest = exports.GroupMember = exports.Group = exports.GetUserGroupsRequest = exports.GetGroupRequest = exports.GroupMemberStatus = exports.IsGroupMemberRequest = exports.GroupList = exports.GetGroupListRequest = exports.UpdateGroupRequest = exports.DeleteGroupRequest = exports.NewGroupRequest = exports.DeleteUserRequest = exports.UpdateUserRequest = exports.NewUserRequest = exports.GetUserRequest = exports.AuthenticateUserRequest = exports.UserList = exports.User = exports.GetUserListRequest = exports.Empty = exports.sortOrderToJSON = exports.sortOrderFromJSON = exports.SortOrder = exports.protobufPackage = void 0;
/* eslint-disable */
var Long = require("long");
var _m0 = require("protobufjs/minimal");
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
    encode: function (_, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseEmpty();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON: function (_) {
        return {};
    },
    toJSON: function (_) {
        var obj = {};
        return obj;
    },
    fromPartial: function (_) {
        var message = createBaseEmpty();
        return message;
    },
};
function createBaseGetUserListRequest() {
    return { start: 0, end: 0, sortOrder: 0, sortKey: "", filter: [] };
}
exports.GetUserListRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
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
        for (var _i = 0, _a = message.filter; _i < _a.length; _i++) {
            var v = _a[_i];
            writer.uint32(82).string(v);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseGetUserListRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.start = reader.int32();
                    break;
                case 2:
                    message.end = reader.int32();
                    break;
                case 3:
                    message.sortOrder = reader.int32();
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
    fromJSON: function (object) {
        return {
            start: isSet(object.start) ? Number(object.start) : 0,
            end: isSet(object.end) ? Number(object.end) : 0,
            sortOrder: isSet(object.sortOrder) ? sortOrderFromJSON(object.sortOrder) : 0,
            sortKey: isSet(object.sortKey) ? String(object.sortKey) : "",
            filter: Array.isArray(object === null || object === void 0 ? void 0 : object.filter) ? object.filter.map(function (e) { return String(e); }) : [],
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
    fromPartial: function (object) {
        var _a, _b, _c, _d, _e;
        var message = createBaseGetUserListRequest();
        message.start = (_a = object.start) !== null && _a !== void 0 ? _a : 0;
        message.end = (_b = object.end) !== null && _b !== void 0 ? _b : 0;
        message.sortOrder = (_c = object.sortOrder) !== null && _c !== void 0 ? _c : 0;
        message.sortKey = (_d = object.sortKey) !== null && _d !== void 0 ? _d : "";
        message.filter = ((_e = object.filter) === null || _e === void 0 ? void 0 : _e.map(function (e) { return e; })) || [];
        return message;
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
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
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
        if (message.email !== "") {
            writer.uint32(106).string(message.email);
        }
        if (message.loginShell !== "") {
            writer.uint32(162).string(message.loginShell);
        }
        if (message.homeDirectory !== "") {
            writer.uint32(170).string(message.homeDirectory);
        }
        if (message.CN !== "") {
            writer.uint32(242).string(message.CN);
        }
        if (message.DN !== "") {
            writer.uint32(250).string(message.DN);
        }
        if (message.UID !== 0) {
            writer.uint32(256).int32(message.UID);
        }
        if (message.GID !== 0) {
            writer.uint32(264).int64(message.GID);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseUser();
        while (reader.pos < end) {
            var tag = reader.uint32();
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
                    message.email = reader.string();
                    break;
                case 20:
                    message.loginShell = reader.string();
                    break;
                case 21:
                    message.homeDirectory = reader.string();
                    break;
                case 30:
                    message.CN = reader.string();
                    break;
                case 31:
                    message.DN = reader.string();
                    break;
                case 32:
                    message.UID = reader.int32();
                    break;
                case 33:
                    message.GID = longToNumber(reader.int64());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
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
    fromPartial: function (object) {
        var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l;
        var message = createBaseUser();
        message.username = (_a = object.username) !== null && _a !== void 0 ? _a : "";
        message.firstName = (_b = object.firstName) !== null && _b !== void 0 ? _b : "";
        message.lastName = (_c = object.lastName) !== null && _c !== void 0 ? _c : "";
        message.displayName = (_d = object.displayName) !== null && _d !== void 0 ? _d : "";
        message.email = (_e = object.email) !== null && _e !== void 0 ? _e : "";
        message.loginShell = (_f = object.loginShell) !== null && _f !== void 0 ? _f : "";
        message.homeDirectory = (_g = object.homeDirectory) !== null && _g !== void 0 ? _g : "";
        message.CN = (_h = object.CN) !== null && _h !== void 0 ? _h : "";
        message.DN = (_j = object.DN) !== null && _j !== void 0 ? _j : "";
        message.UID = (_k = object.UID) !== null && _k !== void 0 ? _k : 0;
        message.GID = (_l = object.GID) !== null && _l !== void 0 ? _l : 0;
        return message;
    },
};
function createBaseUserList() {
    return { users: [], total: 0 };
}
exports.UserList = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
        for (var _i = 0, _a = message.users; _i < _a.length; _i++) {
            var v = _a[_i];
            exports.User.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.total !== 0) {
            writer.uint32(80).int64(message.total);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseUserList();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.users.push(exports.User.decode(reader, reader.uint32()));
                    break;
                case 10:
                    message.total = longToNumber(reader.int64());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
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
    fromPartial: function (object) {
        var _a, _b;
        var message = createBaseUserList();
        message.users = ((_a = object.users) === null || _a === void 0 ? void 0 : _a.map(function (e) { return exports.User.fromPartial(e); })) || [];
        message.total = (_b = object.total) !== null && _b !== void 0 ? _b : 0;
        return message;
    },
};
function createBaseAuthenticateUserRequest() {
    return { username: "", password: "" };
}
exports.AuthenticateUserRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
        if (message.username !== "") {
            writer.uint32(10).string(message.username);
        }
        if (message.password !== "") {
            writer.uint32(18).string(message.password);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseAuthenticateUserRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
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
    fromPartial: function (object) {
        var _a, _b;
        var message = createBaseAuthenticateUserRequest();
        message.username = (_a = object.username) !== null && _a !== void 0 ? _a : "";
        message.password = (_b = object.password) !== null && _b !== void 0 ? _b : "";
        return message;
    },
};
function createBaseGetUserRequest() {
    return { username: "" };
}
exports.GetUserRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
        if (message.username !== "") {
            writer.uint32(10).string(message.username);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseGetUserRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
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
    fromJSON: function (object) {
        return { username: isSet(object.username) ? String(object.username) : "" };
    },
    toJSON: function (message) {
        var obj = {};
        message.username !== undefined && (obj.username = message.username);
        return obj;
    },
    fromPartial: function (object) {
        var _a;
        var message = createBaseGetUserRequest();
        message.username = (_a = object.username) !== null && _a !== void 0 ? _a : "";
        return message;
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
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
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
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseNewUserRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.firstName = reader.string();
                    break;
                case 2:
                    message.lastName = reader.string();
                    break;
                case 10:
                    message.UID = longToNumber(reader.int64());
                    break;
                case 11:
                    message.GID = longToNumber(reader.int64());
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
    fromPartial: function (object) {
        var _a, _b, _c, _d, _e, _f, _g, _h, _j;
        var message = createBaseNewUserRequest();
        message.firstName = (_a = object.firstName) !== null && _a !== void 0 ? _a : "";
        message.lastName = (_b = object.lastName) !== null && _b !== void 0 ? _b : "";
        message.UID = (_c = object.UID) !== null && _c !== void 0 ? _c : 0;
        message.GID = (_d = object.GID) !== null && _d !== void 0 ? _d : 0;
        message.loginShell = (_e = object.loginShell) !== null && _e !== void 0 ? _e : "";
        message.homeDirectory = (_f = object.homeDirectory) !== null && _f !== void 0 ? _f : "";
        message.username = (_g = object.username) !== null && _g !== void 0 ? _g : "";
        message.email = (_h = object.email) !== null && _h !== void 0 ? _h : "";
        message.password = (_j = object.password) !== null && _j !== void 0 ? _j : "";
        return message;
    },
};
function createBaseUpdateUserRequest() {
    return { username: "", update: undefined };
}
exports.UpdateUserRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
        if (message.username !== "") {
            writer.uint32(10).string(message.username);
        }
        if (message.update !== undefined) {
            exports.NewUserRequest.encode(message.update, writer.uint32(82).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseUpdateUserRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.username = reader.string();
                    break;
                case 10:
                    message.update = exports.NewUserRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
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
    fromPartial: function (object) {
        var _a;
        var message = createBaseUpdateUserRequest();
        message.username = (_a = object.username) !== null && _a !== void 0 ? _a : "";
        message.update = (object.update !== undefined && object.update !== null)
            ? exports.NewUserRequest.fromPartial(object.update)
            : undefined;
        return message;
    },
};
function createBaseDeleteUserRequest() {
    return { username: "" };
}
exports.DeleteUserRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
        if (message.username !== "") {
            writer.uint32(10).string(message.username);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseDeleteUserRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
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
    fromJSON: function (object) {
        return { username: isSet(object.username) ? String(object.username) : "" };
    },
    toJSON: function (message) {
        var obj = {};
        message.username !== undefined && (obj.username = message.username);
        return obj;
    },
    fromPartial: function (object) {
        var _a;
        var message = createBaseDeleteUserRequest();
        message.username = (_a = object.username) !== null && _a !== void 0 ? _a : "";
        return message;
    },
};
function createBaseNewGroupRequest() {
    return { name: "", members: [] };
}
exports.NewGroupRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
        if (message.name !== "") {
            writer.uint32(10).string(message.name);
        }
        for (var _i = 0, _a = message.members; _i < _a.length; _i++) {
            var v = _a[_i];
            writer.uint32(18).string(v);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseNewGroupRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
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
    fromPartial: function (object) {
        var _a, _b;
        var message = createBaseNewGroupRequest();
        message.name = (_a = object.name) !== null && _a !== void 0 ? _a : "";
        message.members = ((_b = object.members) === null || _b === void 0 ? void 0 : _b.map(function (e) { return e; })) || [];
        return message;
    },
};
function createBaseDeleteGroupRequest() {
    return { name: "" };
}
exports.DeleteGroupRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
        if (message.name !== "") {
            writer.uint32(10).string(message.name);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseDeleteGroupRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
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
    fromJSON: function (object) {
        return { name: isSet(object.name) ? String(object.name) : "" };
    },
    toJSON: function (message) {
        var obj = {};
        message.name !== undefined && (obj.name = message.name);
        return obj;
    },
    fromPartial: function (object) {
        var _a;
        var message = createBaseDeleteGroupRequest();
        message.name = (_a = object.name) !== null && _a !== void 0 ? _a : "";
        return message;
    },
};
function createBaseUpdateGroupRequest() {
    return { name: "", newName: "", GID: 0 };
}
exports.UpdateGroupRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
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
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseUpdateGroupRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.name = reader.string();
                    break;
                case 2:
                    message.newName = reader.string();
                    break;
                case 3:
                    message.GID = longToNumber(reader.int64());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
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
    fromPartial: function (object) {
        var _a, _b, _c;
        var message = createBaseUpdateGroupRequest();
        message.name = (_a = object.name) !== null && _a !== void 0 ? _a : "";
        message.newName = (_b = object.newName) !== null && _b !== void 0 ? _b : "";
        message.GID = (_c = object.GID) !== null && _c !== void 0 ? _c : 0;
        return message;
    },
};
function createBaseGetGroupListRequest() {
    return { start: 0, end: 0, sortOrder: 0, sortKey: "", filter: [] };
}
exports.GetGroupListRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
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
        for (var _i = 0, _a = message.filter; _i < _a.length; _i++) {
            var v = _a[_i];
            writer.uint32(82).string(v);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseGetGroupListRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.start = reader.int32();
                    break;
                case 2:
                    message.end = reader.int32();
                    break;
                case 3:
                    message.sortOrder = reader.int32();
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
    fromJSON: function (object) {
        return {
            start: isSet(object.start) ? Number(object.start) : 0,
            end: isSet(object.end) ? Number(object.end) : 0,
            sortOrder: isSet(object.sortOrder) ? sortOrderFromJSON(object.sortOrder) : 0,
            sortKey: isSet(object.sortKey) ? String(object.sortKey) : "",
            filter: Array.isArray(object === null || object === void 0 ? void 0 : object.filter) ? object.filter.map(function (e) { return String(e); }) : [],
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
    fromPartial: function (object) {
        var _a, _b, _c, _d, _e;
        var message = createBaseGetGroupListRequest();
        message.start = (_a = object.start) !== null && _a !== void 0 ? _a : 0;
        message.end = (_b = object.end) !== null && _b !== void 0 ? _b : 0;
        message.sortOrder = (_c = object.sortOrder) !== null && _c !== void 0 ? _c : 0;
        message.sortKey = (_d = object.sortKey) !== null && _d !== void 0 ? _d : "";
        message.filter = ((_e = object.filter) === null || _e === void 0 ? void 0 : _e.map(function (e) { return e; })) || [];
        return message;
    },
};
function createBaseGroupList() {
    return { groups: [], total: 0 };
}
exports.GroupList = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
        for (var _i = 0, _a = message.groups; _i < _a.length; _i++) {
            var v = _a[_i];
            exports.Group.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.total !== 0) {
            writer.uint32(80).int64(message.total);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseGroupList();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.groups.push(exports.Group.decode(reader, reader.uint32()));
                    break;
                case 10:
                    message.total = longToNumber(reader.int64());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
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
    fromPartial: function (object) {
        var _a, _b;
        var message = createBaseGroupList();
        message.groups = ((_a = object.groups) === null || _a === void 0 ? void 0 : _a.map(function (e) { return exports.Group.fromPartial(e); })) || [];
        message.total = (_b = object.total) !== null && _b !== void 0 ? _b : 0;
        return message;
    },
};
function createBaseIsGroupMemberRequest() {
    return { username: "", group: "" };
}
exports.IsGroupMemberRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
        if (message.username !== "") {
            writer.uint32(10).string(message.username);
        }
        if (message.group !== "") {
            writer.uint32(18).string(message.group);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseIsGroupMemberRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
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
    fromPartial: function (object) {
        var _a, _b;
        var message = createBaseIsGroupMemberRequest();
        message.username = (_a = object.username) !== null && _a !== void 0 ? _a : "";
        message.group = (_b = object.group) !== null && _b !== void 0 ? _b : "";
        return message;
    },
};
function createBaseGroupMemberStatus() {
    return { isMember: false };
}
exports.GroupMemberStatus = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
        if (message.isMember === true) {
            writer.uint32(8).bool(message.isMember);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseGroupMemberStatus();
        while (reader.pos < end) {
            var tag = reader.uint32();
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
    fromJSON: function (object) {
        return { isMember: isSet(object.isMember) ? Boolean(object.isMember) : false };
    },
    toJSON: function (message) {
        var obj = {};
        message.isMember !== undefined && (obj.isMember = message.isMember);
        return obj;
    },
    fromPartial: function (object) {
        var _a;
        var message = createBaseGroupMemberStatus();
        message.isMember = (_a = object.isMember) !== null && _a !== void 0 ? _a : false;
        return message;
    },
};
function createBaseGetGroupRequest() {
    return { start: 0, end: 0, sortOrder: 0, sortKey: "", name: "" };
}
exports.GetGroupRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
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
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseGetGroupRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.start = reader.int32();
                    break;
                case 2:
                    message.end = reader.int32();
                    break;
                case 3:
                    message.sortOrder = reader.int32();
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
    fromPartial: function (object) {
        var _a, _b, _c, _d, _e;
        var message = createBaseGetGroupRequest();
        message.start = (_a = object.start) !== null && _a !== void 0 ? _a : 0;
        message.end = (_b = object.end) !== null && _b !== void 0 ? _b : 0;
        message.sortOrder = (_c = object.sortOrder) !== null && _c !== void 0 ? _c : 0;
        message.sortKey = (_d = object.sortKey) !== null && _d !== void 0 ? _d : "";
        message.name = (_e = object.name) !== null && _e !== void 0 ? _e : "";
        return message;
    },
};
function createBaseGetUserGroupsRequest() {
    return { username: "" };
}
exports.GetUserGroupsRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
        if (message.username !== "") {
            writer.uint32(10).string(message.username);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseGetUserGroupsRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
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
    fromJSON: function (object) {
        return { username: isSet(object.username) ? String(object.username) : "" };
    },
    toJSON: function (message) {
        var obj = {};
        message.username !== undefined && (obj.username = message.username);
        return obj;
    },
    fromPartial: function (object) {
        var _a;
        var message = createBaseGetUserGroupsRequest();
        message.username = (_a = object.username) !== null && _a !== void 0 ? _a : "";
        return message;
    },
};
function createBaseGroup() {
    return { name: "", members: [], GID: 0 };
}
exports.Group = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
        if (message.name !== "") {
            writer.uint32(10).string(message.name);
        }
        for (var _i = 0, _a = message.members; _i < _a.length; _i++) {
            var v = _a[_i];
            writer.uint32(18).string(v);
        }
        if (message.GID !== 0) {
            writer.uint32(24).int64(message.GID);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseGroup();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.name = reader.string();
                    break;
                case 2:
                    message.members.push(reader.string());
                    break;
                case 3:
                    message.GID = longToNumber(reader.int64());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
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
    fromPartial: function (object) {
        var _a, _b, _c;
        var message = createBaseGroup();
        message.name = (_a = object.name) !== null && _a !== void 0 ? _a : "";
        message.members = ((_b = object.members) === null || _b === void 0 ? void 0 : _b.map(function (e) { return e; })) || [];
        message.GID = (_c = object.GID) !== null && _c !== void 0 ? _c : 0;
        return message;
    },
};
function createBaseGroupMember() {
    return { group: "", username: "" };
}
exports.GroupMember = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
        if (message.group !== "") {
            writer.uint32(10).string(message.group);
        }
        if (message.username !== "") {
            writer.uint32(18).string(message.username);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseGroupMember();
        while (reader.pos < end) {
            var tag = reader.uint32();
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
    fromPartial: function (object) {
        var _a, _b;
        var message = createBaseGroupMember();
        message.group = (_a = object.group) !== null && _a !== void 0 ? _a : "";
        message.username = (_b = object.username) !== null && _b !== void 0 ? _b : "";
        return message;
    },
};
function createBaseChangePasswordRequest() {
    return { username: "", password: "" };
}
exports.ChangePasswordRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
        if (message.username !== "") {
            writer.uint32(10).string(message.username);
        }
        if (message.password !== "") {
            writer.uint32(18).string(message.password);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseChangePasswordRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
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
    fromPartial: function (object) {
        var _a, _b;
        var message = createBaseChangePasswordRequest();
        message.username = (_a = object.username) !== null && _a !== void 0 ? _a : "";
        message.password = (_b = object.password) !== null && _b !== void 0 ? _b : "";
        return message;
    },
};
function createBaseLoginRequest() {
    return { username: "", password: "" };
}
exports.LoginRequest = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
        if (message.username !== "") {
            writer.uint32(10).string(message.username);
        }
        if (message.password !== "") {
            writer.uint32(18).string(message.password);
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseLoginRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
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
    fromPartial: function (object) {
        var _a, _b;
        var message = createBaseLoginRequest();
        message.username = (_a = object.username) !== null && _a !== void 0 ? _a : "";
        message.password = (_b = object.password) !== null && _b !== void 0 ? _b : "";
        return message;
    },
};
function createBaseToken() {
    return { token: "", username: "", UID: 0, displayName: "", isAdmin: false, expires: undefined };
}
exports.Token = {
    encode: function (message, writer) {
        if (writer === void 0) { writer = _m0.Writer.create(); }
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
            timestamp_1.Timestamp.encode(toTimestamp(message.expires), writer.uint32(82).fork()).ldelim();
        }
        return writer;
    },
    decode: function (input, length) {
        var reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        var end = length === undefined ? reader.len : reader.pos + length;
        var message = createBaseToken();
        while (reader.pos < end) {
            var tag = reader.uint32();
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
                    message.expires = fromTimestamp(timestamp_1.Timestamp.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
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
    fromPartial: function (object) {
        var _a, _b, _c, _d, _e, _f;
        var message = createBaseToken();
        message.token = (_a = object.token) !== null && _a !== void 0 ? _a : "";
        message.username = (_b = object.username) !== null && _b !== void 0 ? _b : "";
        message.UID = (_c = object.UID) !== null && _c !== void 0 ? _c : 0;
        message.displayName = (_d = object.displayName) !== null && _d !== void 0 ? _d : "";
        message.isAdmin = (_e = object.isAdmin) !== null && _e !== void 0 ? _e : false;
        message.expires = (_f = object.expires) !== null && _f !== void 0 ? _f : undefined;
        return message;
    },
};
var LDAPManagerClientImpl = /** @class */ (function () {
    function LDAPManagerClientImpl(rpc, opts) {
        this.service = (opts === null || opts === void 0 ? void 0 : opts.service) || "ldapmanager.LDAPManager";
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
    LDAPManagerClientImpl.prototype.Login = function (request) {
        var data = exports.LoginRequest.encode(request).finish();
        var promise = this.rpc.request(this.service, "Login", data);
        return promise.then(function (data) { return exports.Token.decode(new _m0.Reader(data)); });
    };
    LDAPManagerClientImpl.prototype.GetUserList = function (request) {
        var data = exports.GetUserListRequest.encode(request).finish();
        var promise = this.rpc.request(this.service, "GetUserList", data);
        return promise.then(function (data) { return exports.UserList.decode(new _m0.Reader(data)); });
    };
    LDAPManagerClientImpl.prototype.GetUser = function (request) {
        var data = exports.GetUserRequest.encode(request).finish();
        var promise = this.rpc.request(this.service, "GetUser", data);
        return promise.then(function (data) { return exports.User.decode(new _m0.Reader(data)); });
    };
    LDAPManagerClientImpl.prototype.NewUser = function (request) {
        var data = exports.NewUserRequest.encode(request).finish();
        var promise = this.rpc.request(this.service, "NewUser", data);
        return promise.then(function (data) { return exports.Empty.decode(new _m0.Reader(data)); });
    };
    LDAPManagerClientImpl.prototype.UpdateUser = function (request) {
        var data = exports.UpdateUserRequest.encode(request).finish();
        var promise = this.rpc.request(this.service, "UpdateUser", data);
        return promise.then(function (data) { return exports.Token.decode(new _m0.Reader(data)); });
    };
    LDAPManagerClientImpl.prototype.DeleteUser = function (request) {
        var data = exports.DeleteUserRequest.encode(request).finish();
        var promise = this.rpc.request(this.service, "DeleteUser", data);
        return promise.then(function (data) { return exports.Empty.decode(new _m0.Reader(data)); });
    };
    LDAPManagerClientImpl.prototype.ChangePassword = function (request) {
        var data = exports.ChangePasswordRequest.encode(request).finish();
        var promise = this.rpc.request(this.service, "ChangePassword", data);
        return promise.then(function (data) { return exports.Empty.decode(new _m0.Reader(data)); });
    };
    LDAPManagerClientImpl.prototype.NewGroup = function (request) {
        var data = exports.NewGroupRequest.encode(request).finish();
        var promise = this.rpc.request(this.service, "NewGroup", data);
        return promise.then(function (data) { return exports.Empty.decode(new _m0.Reader(data)); });
    };
    LDAPManagerClientImpl.prototype.DeleteGroup = function (request) {
        var data = exports.DeleteGroupRequest.encode(request).finish();
        var promise = this.rpc.request(this.service, "DeleteGroup", data);
        return promise.then(function (data) { return exports.Empty.decode(new _m0.Reader(data)); });
    };
    LDAPManagerClientImpl.prototype.UpdateGroup = function (request) {
        var data = exports.UpdateGroupRequest.encode(request).finish();
        var promise = this.rpc.request(this.service, "UpdateGroup", data);
        return promise.then(function (data) { return exports.Empty.decode(new _m0.Reader(data)); });
    };
    LDAPManagerClientImpl.prototype.GetGroupList = function (request) {
        var data = exports.GetGroupListRequest.encode(request).finish();
        var promise = this.rpc.request(this.service, "GetGroupList", data);
        return promise.then(function (data) { return exports.GroupList.decode(new _m0.Reader(data)); });
    };
    LDAPManagerClientImpl.prototype.GetUserGroups = function (request) {
        var data = exports.GetUserGroupsRequest.encode(request).finish();
        var promise = this.rpc.request(this.service, "GetUserGroups", data);
        return promise.then(function (data) { return exports.GroupList.decode(new _m0.Reader(data)); });
    };
    LDAPManagerClientImpl.prototype.IsGroupMember = function (request) {
        var data = exports.IsGroupMemberRequest.encode(request).finish();
        var promise = this.rpc.request(this.service, "IsGroupMember", data);
        return promise.then(function (data) { return exports.GroupMemberStatus.decode(new _m0.Reader(data)); });
    };
    LDAPManagerClientImpl.prototype.GetGroup = function (request) {
        var data = exports.GetGroupRequest.encode(request).finish();
        var promise = this.rpc.request(this.service, "GetGroup", data);
        return promise.then(function (data) { return exports.Group.decode(new _m0.Reader(data)); });
    };
    LDAPManagerClientImpl.prototype.AddGroupMember = function (request) {
        var data = exports.GroupMember.encode(request).finish();
        var promise = this.rpc.request(this.service, "AddGroupMember", data);
        return promise.then(function (data) { return exports.Empty.decode(new _m0.Reader(data)); });
    };
    LDAPManagerClientImpl.prototype.RemoveGroupMember = function (request) {
        var data = exports.GroupMember.encode(request).finish();
        var promise = this.rpc.request(this.service, "RemoveGroupMember", data);
        return promise.then(function (data) { return exports.Empty.decode(new _m0.Reader(data)); });
    };
    return LDAPManagerClientImpl;
}());
exports.LDAPManagerClientImpl = LDAPManagerClientImpl;
var globalThis = (function () {
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
function toTimestamp(date) {
    var seconds = date.getTime() / 1000;
    var nanos = (date.getTime() % 1000) * 1000000;
    return { seconds: seconds, nanos: nanos };
}
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
function longToNumber(long) {
    if (long.gt(Number.MAX_SAFE_INTEGER)) {
        throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
    }
    return long.toNumber();
}
// If you get a compile-error about 'Constructor<Long> and ... have no overlap',
// add '--ts_proto_opt=esModuleInterop=true' as a flag when calling 'protoc'.
if (_m0.util.Long !== Long) {
    _m0.util.Long = Long;
    _m0.configure();
}
function isSet(value) {
    return value !== null && value !== undefined;
}
