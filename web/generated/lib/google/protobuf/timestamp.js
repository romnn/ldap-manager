"use strict";
/* eslint-disable */
Object.defineProperty(exports, "__esModule", { value: true });
exports.Timestamp = exports.protobufPackage = void 0;
exports.protobufPackage = "google.protobuf";
function createBaseTimestamp() {
    return { seconds: 0, nanos: 0 };
}
exports.Timestamp = {
    fromJSON: function (object) {
        return {
            seconds: isSet(object.seconds) ? Number(object.seconds) : 0,
            nanos: isSet(object.nanos) ? Number(object.nanos) : 0,
        };
    },
    toJSON: function (message) {
        var obj = {};
        message.seconds !== undefined && (obj.seconds = Math.round(message.seconds));
        message.nanos !== undefined && (obj.nanos = Math.round(message.nanos));
        return obj;
    },
};
function isSet(value) {
    return value !== null && value !== undefined;
}
//# sourceMappingURL=timestamp.js.map