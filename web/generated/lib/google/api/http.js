/* eslint-disable */
(function (factory) {
    if (typeof module === "object" && typeof module.exports === "object") {
        var v = factory(require, exports);
        if (v !== undefined) module.exports = v;
    }
    else if (typeof define === "function" && define.amd) {
        define(["require", "exports"], factory);
    }
})(function (require, exports) {
    "use strict";
    exports.__esModule = true;
    exports.CustomHttpPattern = exports.HttpRule = exports.Http = exports.protobufPackage = void 0;
    exports.protobufPackage = "google.api";
    function createBaseHttp() {
        return { rules: [], fullyDecodeReservedExpansion: false };
    }
    exports.Http = {
        fromJSON: function (object) {
            return {
                rules: Array.isArray(object === null || object === void 0 ? void 0 : object.rules) ? object.rules.map(function (e) { return exports.HttpRule.fromJSON(e); }) : [],
                fullyDecodeReservedExpansion: isSet(object.fullyDecodeReservedExpansion)
                    ? Boolean(object.fullyDecodeReservedExpansion)
                    : false
            };
        },
        toJSON: function (message) {
            var obj = {};
            if (message.rules) {
                obj.rules = message.rules.map(function (e) { return e ? exports.HttpRule.toJSON(e) : undefined; });
            }
            else {
                obj.rules = [];
            }
            message.fullyDecodeReservedExpansion !== undefined &&
                (obj.fullyDecodeReservedExpansion = message.fullyDecodeReservedExpansion);
            return obj;
        }
    };
    function createBaseHttpRule() {
        return {
            selector: "",
            get: undefined,
            put: undefined,
            post: undefined,
            "delete": undefined,
            patch: undefined,
            custom: undefined,
            body: "",
            responseBody: "",
            additionalBindings: []
        };
    }
    exports.HttpRule = {
        fromJSON: function (object) {
            return {
                selector: isSet(object.selector) ? String(object.selector) : "",
                get: isSet(object.get) ? String(object.get) : undefined,
                put: isSet(object.put) ? String(object.put) : undefined,
                post: isSet(object.post) ? String(object.post) : undefined,
                "delete": isSet(object["delete"]) ? String(object["delete"]) : undefined,
                patch: isSet(object.patch) ? String(object.patch) : undefined,
                custom: isSet(object.custom) ? exports.CustomHttpPattern.fromJSON(object.custom) : undefined,
                body: isSet(object.body) ? String(object.body) : "",
                responseBody: isSet(object.responseBody) ? String(object.responseBody) : "",
                additionalBindings: Array.isArray(object === null || object === void 0 ? void 0 : object.additionalBindings)
                    ? object.additionalBindings.map(function (e) { return exports.HttpRule.fromJSON(e); })
                    : []
            };
        },
        toJSON: function (message) {
            var obj = {};
            message.selector !== undefined && (obj.selector = message.selector);
            message.get !== undefined && (obj.get = message.get);
            message.put !== undefined && (obj.put = message.put);
            message.post !== undefined && (obj.post = message.post);
            message["delete"] !== undefined && (obj["delete"] = message["delete"]);
            message.patch !== undefined && (obj.patch = message.patch);
            message.custom !== undefined &&
                (obj.custom = message.custom ? exports.CustomHttpPattern.toJSON(message.custom) : undefined);
            message.body !== undefined && (obj.body = message.body);
            message.responseBody !== undefined && (obj.responseBody = message.responseBody);
            if (message.additionalBindings) {
                obj.additionalBindings = message.additionalBindings.map(function (e) { return e ? exports.HttpRule.toJSON(e) : undefined; });
            }
            else {
                obj.additionalBindings = [];
            }
            return obj;
        }
    };
    function createBaseCustomHttpPattern() {
        return { kind: "", path: "" };
    }
    exports.CustomHttpPattern = {
        fromJSON: function (object) {
            return { kind: isSet(object.kind) ? String(object.kind) : "", path: isSet(object.path) ? String(object.path) : "" };
        },
        toJSON: function (message) {
            var obj = {};
            message.kind !== undefined && (obj.kind = message.kind);
            message.path !== undefined && (obj.path = message.path);
            return obj;
        }
    };
    function isSet(value) {
        return value !== null && value !== undefined;
    }
});
