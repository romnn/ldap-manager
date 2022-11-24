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
    exports.Timestamp = exports.protobufPackage = void 0;
    exports.protobufPackage = "google.protobuf";
    function createBaseTimestamp() {
        return { seconds: 0, nanos: 0 };
    }
    exports.Timestamp = {
        fromJSON: function (object) {
            return {
                seconds: isSet(object.seconds) ? Number(object.seconds) : 0,
                nanos: isSet(object.nanos) ? Number(object.nanos) : 0
            };
        },
        toJSON: function (message) {
            var obj = {};
            message.seconds !== undefined && (obj.seconds = Math.round(message.seconds));
            message.nanos !== undefined && (obj.nanos = Math.round(message.nanos));
            return obj;
        }
    };
    function isSet(value) {
        return value !== null && value !== undefined;
    }
});
