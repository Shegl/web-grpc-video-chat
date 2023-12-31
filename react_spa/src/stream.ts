// @generated by protobuf-ts 2.9.0
// @generated from protobuf file "stream.proto" (syntax proto3)
// tslint:disable
import { ServiceType } from "@protobuf-ts/runtime-rpc";
import { WireType } from "@protobuf-ts/runtime";
import type { BinaryWriteOptions } from "@protobuf-ts/runtime";
import type { IBinaryWriter } from "@protobuf-ts/runtime";
import { UnknownFieldHandler } from "@protobuf-ts/runtime";
import type { BinaryReadOptions } from "@protobuf-ts/runtime";
import type { IBinaryReader } from "@protobuf-ts/runtime";
import type { PartialMessage } from "@protobuf-ts/runtime";
import { reflectionMergePartial } from "@protobuf-ts/runtime";
import { MESSAGE_TYPE } from "@protobuf-ts/runtime";
import { MessageType } from "@protobuf-ts/runtime";
/**
 * @generated from protobuf message Ack
 */
export interface Ack {
}
/**
 * @generated from protobuf message StateMessage
 */
export interface StateMessage {
    /**
     * @generated from protobuf field: int64 Time = 1 [json_name = "Time"];
     */
    time: bigint;
    /**
     * @generated from protobuf field: string UUID = 2 [json_name = "UUID"];
     */
    uUID: string;
    /**
     * @generated from protobuf field: User Author = 3 [json_name = "Author"];
     */
    author?: User;
    /**
     * @generated from protobuf field: User Guest = 4 [json_name = "Guest"];
     */
    guest?: User;
}
/**
 * @generated from protobuf message User
 */
export interface User {
    /**
     * @generated from protobuf field: bool IsCamEnabled = 1 [json_name = "IsCamEnabled"];
     */
    isCamEnabled: boolean;
    /**
     * @generated from protobuf field: bool IsMuted = 2 [json_name = "IsMuted"];
     */
    isMuted: boolean;
    /**
     * @generated from protobuf field: string UserUUID = 3 [json_name = "UserUUID"];
     */
    userUUID: string;
    /**
     * @generated from protobuf field: string UserName = 4 [json_name = "UserName"];
     */
    userName: string;
    /**
     * @generated from protobuf field: string UserRoom = 5 [json_name = "UserRoom"];
     */
    userRoom: string;
}
/**
 * @generated from protobuf message AVFrameData
 */
export interface AVFrameData {
    /**
     * @generated from protobuf field: string UserUUID = 1 [json_name = "UserUUID"];
     */
    userUUID: string;
    /**
     * @generated from protobuf field: bytes FrameData = 2 [json_name = "FrameData"];
     */
    frameData: Uint8Array;
}
// @generated message type with reflection information, may provide speed optimized methods
class Ack$Type extends MessageType<Ack> {
    constructor() {
        super("Ack", []);
    }
    create(value?: PartialMessage<Ack>): Ack {
        const message = {};
        globalThis.Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial<Ack>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Ack): Ack {
        return target ?? this.create();
    }
    internalBinaryWrite(message: Ack, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message Ack
 */
export const Ack = new Ack$Type();
// @generated message type with reflection information, may provide speed optimized methods
class StateMessage$Type extends MessageType<StateMessage> {
    constructor() {
        super("StateMessage", [
            { no: 1, name: "Time", kind: "scalar", jsonName: "Time", T: 3 /*ScalarType.INT64*/, L: 0 /*LongType.BIGINT*/ },
            { no: 2, name: "UUID", kind: "scalar", jsonName: "UUID", T: 9 /*ScalarType.STRING*/ },
            { no: 3, name: "Author", kind: "message", jsonName: "Author", T: () => User },
            { no: 4, name: "Guest", kind: "message", jsonName: "Guest", T: () => User }
        ]);
    }
    create(value?: PartialMessage<StateMessage>): StateMessage {
        const message = { time: 0n, uUID: "" };
        globalThis.Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial<StateMessage>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: StateMessage): StateMessage {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* int64 Time = 1 [json_name = "Time"];*/ 1:
                    message.time = reader.int64().toBigInt();
                    break;
                case /* string UUID = 2 [json_name = "UUID"];*/ 2:
                    message.uUID = reader.string();
                    break;
                case /* User Author = 3 [json_name = "Author"];*/ 3:
                    message.author = User.internalBinaryRead(reader, reader.uint32(), options, message.author);
                    break;
                case /* User Guest = 4 [json_name = "Guest"];*/ 4:
                    message.guest = User.internalBinaryRead(reader, reader.uint32(), options, message.guest);
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: StateMessage, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* int64 Time = 1 [json_name = "Time"]; */
        if (message.time !== 0n)
            writer.tag(1, WireType.Varint).int64(message.time);
        /* string UUID = 2 [json_name = "UUID"]; */
        if (message.uUID !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.uUID);
        /* User Author = 3 [json_name = "Author"]; */
        if (message.author)
            User.internalBinaryWrite(message.author, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* User Guest = 4 [json_name = "Guest"]; */
        if (message.guest)
            User.internalBinaryWrite(message.guest, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message StateMessage
 */
export const StateMessage = new StateMessage$Type();
// @generated message type with reflection information, may provide speed optimized methods
class User$Type extends MessageType<User> {
    constructor() {
        super("User", [
            { no: 1, name: "IsCamEnabled", kind: "scalar", jsonName: "IsCamEnabled", T: 8 /*ScalarType.BOOL*/ },
            { no: 2, name: "IsMuted", kind: "scalar", jsonName: "IsMuted", T: 8 /*ScalarType.BOOL*/ },
            { no: 3, name: "UserUUID", kind: "scalar", jsonName: "UserUUID", T: 9 /*ScalarType.STRING*/ },
            { no: 4, name: "UserName", kind: "scalar", jsonName: "UserName", T: 9 /*ScalarType.STRING*/ },
            { no: 5, name: "UserRoom", kind: "scalar", jsonName: "UserRoom", T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<User>): User {
        const message = { isCamEnabled: false, isMuted: false, userUUID: "", userName: "", userRoom: "" };
        globalThis.Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial<User>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: User): User {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* bool IsCamEnabled = 1 [json_name = "IsCamEnabled"];*/ 1:
                    message.isCamEnabled = reader.bool();
                    break;
                case /* bool IsMuted = 2 [json_name = "IsMuted"];*/ 2:
                    message.isMuted = reader.bool();
                    break;
                case /* string UserUUID = 3 [json_name = "UserUUID"];*/ 3:
                    message.userUUID = reader.string();
                    break;
                case /* string UserName = 4 [json_name = "UserName"];*/ 4:
                    message.userName = reader.string();
                    break;
                case /* string UserRoom = 5 [json_name = "UserRoom"];*/ 5:
                    message.userRoom = reader.string();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: User, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* bool IsCamEnabled = 1 [json_name = "IsCamEnabled"]; */
        if (message.isCamEnabled !== false)
            writer.tag(1, WireType.Varint).bool(message.isCamEnabled);
        /* bool IsMuted = 2 [json_name = "IsMuted"]; */
        if (message.isMuted !== false)
            writer.tag(2, WireType.Varint).bool(message.isMuted);
        /* string UserUUID = 3 [json_name = "UserUUID"]; */
        if (message.userUUID !== "")
            writer.tag(3, WireType.LengthDelimited).string(message.userUUID);
        /* string UserName = 4 [json_name = "UserName"]; */
        if (message.userName !== "")
            writer.tag(4, WireType.LengthDelimited).string(message.userName);
        /* string UserRoom = 5 [json_name = "UserRoom"]; */
        if (message.userRoom !== "")
            writer.tag(5, WireType.LengthDelimited).string(message.userRoom);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message User
 */
export const User = new User$Type();
// @generated message type with reflection information, may provide speed optimized methods
class AVFrameData$Type extends MessageType<AVFrameData> {
    constructor() {
        super("AVFrameData", [
            { no: 1, name: "UserUUID", kind: "scalar", jsonName: "UserUUID", T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "FrameData", kind: "scalar", jsonName: "FrameData", T: 12 /*ScalarType.BYTES*/ }
        ]);
    }
    create(value?: PartialMessage<AVFrameData>): AVFrameData {
        const message = { userUUID: "", frameData: new Uint8Array(0) };
        globalThis.Object.defineProperty(message, MESSAGE_TYPE, { enumerable: false, value: this });
        if (value !== undefined)
            reflectionMergePartial<AVFrameData>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: AVFrameData): AVFrameData {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string UserUUID = 1 [json_name = "UserUUID"];*/ 1:
                    message.userUUID = reader.string();
                    break;
                case /* bytes FrameData = 2 [json_name = "FrameData"];*/ 2:
                    message.frameData = reader.bytes();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: AVFrameData, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string UserUUID = 1 [json_name = "UserUUID"]; */
        if (message.userUUID !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.userUUID);
        /* bytes FrameData = 2 [json_name = "FrameData"]; */
        if (message.frameData.length)
            writer.tag(2, WireType.LengthDelimited).bytes(message.frameData);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message AVFrameData
 */
export const AVFrameData = new AVFrameData$Type();
/**
 * @generated ServiceType for protobuf service Stream
 */
export const Stream = new ServiceType("Stream", [
    { name: "StreamState", serverStreaming: true, options: {}, I: User, O: StateMessage },
    { name: "ChangeState", options: {}, I: User, O: Ack },
    { name: "AVStream", serverStreaming: true, options: {}, I: User, O: AVFrameData }
]);
