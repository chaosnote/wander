/*eslint-disable block-scoped-var, id-length, no-control-regex, no-magic-numbers, no-prototype-builtins, no-redeclare, no-shadow, no-var, sort-vars*/
"use strict";

var $protobuf = require("protobufjs/minimal.js");

// Common aliases
var $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;

// Exported root namespace
var $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});

$root.GameMessage = (function() {

    /**
     * Properties of a GameMessage.
     * @exports IGameMessage
     * @interface IGameMessage
     * @property {string|null} [error] GameMessage error
     * @property {GameMessage.MessageType|null} [type] GameMessage type
     * @property {string|null} [action] GameMessage action
     * @property {Uint8Array|null} [payload] GameMessage payload
     * @property {number|Long|null} [timestamp] GameMessage timestamp
     */

    /**
     * Constructs a new GameMessage.
     * @exports GameMessage
     * @classdesc Represents a GameMessage.
     * @implements IGameMessage
     * @constructor
     * @param {IGameMessage=} [properties] Properties to set
     */
    function GameMessage(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * GameMessage error.
     * @member {string|null|undefined} error
     * @memberof GameMessage
     * @instance
     */
    GameMessage.prototype.error = null;

    /**
     * GameMessage type.
     * @member {GameMessage.MessageType} type
     * @memberof GameMessage
     * @instance
     */
    GameMessage.prototype.type = 0;

    /**
     * GameMessage action.
     * @member {string} action
     * @memberof GameMessage
     * @instance
     */
    GameMessage.prototype.action = "";

    /**
     * GameMessage payload.
     * @member {Uint8Array} payload
     * @memberof GameMessage
     * @instance
     */
    GameMessage.prototype.payload = $util.newBuffer([]);

    /**
     * GameMessage timestamp.
     * @member {number|Long} timestamp
     * @memberof GameMessage
     * @instance
     */
    GameMessage.prototype.timestamp = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

    // OneOf field names bound to virtual getters and setters
    var $oneOfFields;

    /**
     * GameMessage _error.
     * @member {"error"|undefined} _error
     * @memberof GameMessage
     * @instance
     */
    Object.defineProperty(GameMessage.prototype, "_error", {
        get: $util.oneOfGetter($oneOfFields = ["error"]),
        set: $util.oneOfSetter($oneOfFields)
    });

    /**
     * Creates a new GameMessage instance using the specified properties.
     * @function create
     * @memberof GameMessage
     * @static
     * @param {IGameMessage=} [properties] Properties to set
     * @returns {GameMessage} GameMessage instance
     */
    GameMessage.create = function create(properties) {
        return new GameMessage(properties);
    };

    /**
     * Encodes the specified GameMessage message. Does not implicitly {@link GameMessage.verify|verify} messages.
     * @function encode
     * @memberof GameMessage
     * @static
     * @param {IGameMessage} message GameMessage message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    GameMessage.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.error != null && Object.hasOwnProperty.call(message, "error"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.error);
        if (message.type != null && Object.hasOwnProperty.call(message, "type"))
            writer.uint32(/* id 2, wireType 0 =*/16).int32(message.type);
        if (message.action != null && Object.hasOwnProperty.call(message, "action"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.action);
        if (message.payload != null && Object.hasOwnProperty.call(message, "payload"))
            writer.uint32(/* id 4, wireType 2 =*/34).bytes(message.payload);
        if (message.timestamp != null && Object.hasOwnProperty.call(message, "timestamp"))
            writer.uint32(/* id 5, wireType 0 =*/40).int64(message.timestamp);
        return writer;
    };

    /**
     * Encodes the specified GameMessage message, length delimited. Does not implicitly {@link GameMessage.verify|verify} messages.
     * @function encodeDelimited
     * @memberof GameMessage
     * @static
     * @param {IGameMessage} message GameMessage message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    GameMessage.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a GameMessage message from the specified reader or buffer.
     * @function decode
     * @memberof GameMessage
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {GameMessage} GameMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    GameMessage.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.GameMessage();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1: {
                    message.error = reader.string();
                    break;
                }
            case 2: {
                    message.type = reader.int32();
                    break;
                }
            case 3: {
                    message.action = reader.string();
                    break;
                }
            case 4: {
                    message.payload = reader.bytes();
                    break;
                }
            case 5: {
                    message.timestamp = reader.int64();
                    break;
                }
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a GameMessage message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof GameMessage
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {GameMessage} GameMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    GameMessage.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a GameMessage message.
     * @function verify
     * @memberof GameMessage
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    GameMessage.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        var properties = {};
        if (message.error != null && message.hasOwnProperty("error")) {
            properties._error = 1;
            if (!$util.isString(message.error))
                return "error: string expected";
        }
        if (message.type != null && message.hasOwnProperty("type"))
            switch (message.type) {
            default:
                return "type: enum value expected";
            case 0:
            case 1:
            case 2:
            case 3:
                break;
            }
        if (message.action != null && message.hasOwnProperty("action"))
            if (!$util.isString(message.action))
                return "action: string expected";
        if (message.payload != null && message.hasOwnProperty("payload"))
            if (!(message.payload && typeof message.payload.length === "number" || $util.isString(message.payload)))
                return "payload: buffer expected";
        if (message.timestamp != null && message.hasOwnProperty("timestamp"))
            if (!$util.isInteger(message.timestamp) && !(message.timestamp && $util.isInteger(message.timestamp.low) && $util.isInteger(message.timestamp.high)))
                return "timestamp: integer|Long expected";
        return null;
    };

    /**
     * Creates a GameMessage message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof GameMessage
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {GameMessage} GameMessage
     */
    GameMessage.fromObject = function fromObject(object) {
        if (object instanceof $root.GameMessage)
            return object;
        var message = new $root.GameMessage();
        if (object.error != null)
            message.error = String(object.error);
        switch (object.type) {
        default:
            if (typeof object.type === "number") {
                message.type = object.type;
                break;
            }
            break;
        case "REQUEST":
        case 0:
            message.type = 0;
            break;
        case "RESPONSE":
        case 1:
            message.type = 1;
            break;
        case "NOTIFY":
        case 2:
            message.type = 2;
            break;
        case "ALERT":
        case 3:
            message.type = 3;
            break;
        }
        if (object.action != null)
            message.action = String(object.action);
        if (object.payload != null)
            if (typeof object.payload === "string")
                $util.base64.decode(object.payload, message.payload = $util.newBuffer($util.base64.length(object.payload)), 0);
            else if (object.payload.length >= 0)
                message.payload = object.payload;
        if (object.timestamp != null)
            if ($util.Long)
                (message.timestamp = $util.Long.fromValue(object.timestamp)).unsigned = false;
            else if (typeof object.timestamp === "string")
                message.timestamp = parseInt(object.timestamp, 10);
            else if (typeof object.timestamp === "number")
                message.timestamp = object.timestamp;
            else if (typeof object.timestamp === "object")
                message.timestamp = new $util.LongBits(object.timestamp.low >>> 0, object.timestamp.high >>> 0).toNumber();
        return message;
    };

    /**
     * Creates a plain object from a GameMessage message. Also converts values to other types if specified.
     * @function toObject
     * @memberof GameMessage
     * @static
     * @param {GameMessage} message GameMessage
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    GameMessage.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.type = options.enums === String ? "REQUEST" : 0;
            object.action = "";
            if (options.bytes === String)
                object.payload = "";
            else {
                object.payload = [];
                if (options.bytes !== Array)
                    object.payload = $util.newBuffer(object.payload);
            }
            if ($util.Long) {
                var long = new $util.Long(0, 0, false);
                object.timestamp = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
            } else
                object.timestamp = options.longs === String ? "0" : 0;
        }
        if (message.error != null && message.hasOwnProperty("error")) {
            object.error = message.error;
            if (options.oneofs)
                object._error = "error";
        }
        if (message.type != null && message.hasOwnProperty("type"))
            object.type = options.enums === String ? $root.GameMessage.MessageType[message.type] === undefined ? message.type : $root.GameMessage.MessageType[message.type] : message.type;
        if (message.action != null && message.hasOwnProperty("action"))
            object.action = message.action;
        if (message.payload != null && message.hasOwnProperty("payload"))
            object.payload = options.bytes === String ? $util.base64.encode(message.payload, 0, message.payload.length) : options.bytes === Array ? Array.prototype.slice.call(message.payload) : message.payload;
        if (message.timestamp != null && message.hasOwnProperty("timestamp"))
            if (typeof message.timestamp === "number")
                object.timestamp = options.longs === String ? String(message.timestamp) : message.timestamp;
            else
                object.timestamp = options.longs === String ? $util.Long.prototype.toString.call(message.timestamp) : options.longs === Number ? new $util.LongBits(message.timestamp.low >>> 0, message.timestamp.high >>> 0).toNumber() : message.timestamp;
        return object;
    };

    /**
     * Converts this GameMessage to JSON.
     * @function toJSON
     * @memberof GameMessage
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    GameMessage.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    /**
     * Gets the default type url for GameMessage
     * @function getTypeUrl
     * @memberof GameMessage
     * @static
     * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
     * @returns {string} The default type url
     */
    GameMessage.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
        if (typeUrlPrefix === undefined) {
            typeUrlPrefix = "type.googleapis.com";
        }
        return typeUrlPrefix + "/GameMessage";
    };

    /**
     * MessageType enum.
     * @name GameMessage.MessageType
     * @enum {number}
     * @property {number} REQUEST=0 REQUEST value
     * @property {number} RESPONSE=1 RESPONSE value
     * @property {number} NOTIFY=2 NOTIFY value
     * @property {number} ALERT=3 ALERT value
     */
    GameMessage.MessageType = (function() {
        var valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "REQUEST"] = 0;
        values[valuesById[1] = "RESPONSE"] = 1;
        values[valuesById[2] = "NOTIFY"] = 2;
        values[valuesById[3] = "ALERT"] = 3;
        return values;
    })();

    return GameMessage;
})();

module.exports = $root;
