/* eslint-disable */
import { grpc } from "@improbable-eng/grpc-web";
import { BrowserHeaders } from "browser-headers";
import * as _m0 from "protobufjs/minimal";
import { Observable } from "rxjs";
import { share } from "rxjs/operators";
import Long = require("long");

export const protobufPackage = "";

export interface Empty {
}

export interface AuthRequest {
  UUID: string;
  ChatUUID: string;
}

export interface SendMessageRequest {
  Msg: string;
  AuthData: AuthRequest | undefined;
}

export interface ChatMessage {
  UUID: string;
  UserUUID: string;
  UserName: string;
  Time: number;
  Msg: string;
}

export interface HistoryResponse {
  Messages: ChatMessage[];
}

function createBaseEmpty(): Empty {
  return {};
}

export const Empty = {
  encode(_: Empty, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Empty {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEmpty();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
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

  create<I extends Exact<DeepPartial<Empty>, I>>(base?: I): Empty {
    return Empty.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Empty>, I>>(_: I): Empty {
    const message = createBaseEmpty();
    return message;
  },
};

function createBaseAuthRequest(): AuthRequest {
  return { UUID: "", ChatUUID: "" };
}

export const AuthRequest = {
  encode(message: AuthRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.UUID !== "") {
      writer.uint32(10).string(message.UUID);
    }
    if (message.ChatUUID !== "") {
      writer.uint32(18).string(message.ChatUUID);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AuthRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAuthRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.UUID = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.ChatUUID = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AuthRequest {
    return {
      UUID: isSet(object.UUID) ? String(object.UUID) : "",
      ChatUUID: isSet(object.ChatUUID) ? String(object.ChatUUID) : "",
    };
  },

  toJSON(message: AuthRequest): unknown {
    const obj: any = {};
    if (message.UUID !== "") {
      obj.UUID = message.UUID;
    }
    if (message.ChatUUID !== "") {
      obj.ChatUUID = message.ChatUUID;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AuthRequest>, I>>(base?: I): AuthRequest {
    return AuthRequest.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<AuthRequest>, I>>(object: I): AuthRequest {
    const message = createBaseAuthRequest();
    message.UUID = object.UUID ?? "";
    message.ChatUUID = object.ChatUUID ?? "";
    return message;
  },
};

function createBaseSendMessageRequest(): SendMessageRequest {
  return { Msg: "", AuthData: undefined };
}

export const SendMessageRequest = {
  encode(message: SendMessageRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Msg !== "") {
      writer.uint32(10).string(message.Msg);
    }
    if (message.AuthData !== undefined) {
      AuthRequest.encode(message.AuthData, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SendMessageRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSendMessageRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Msg = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.AuthData = AuthRequest.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SendMessageRequest {
    return {
      Msg: isSet(object.Msg) ? String(object.Msg) : "",
      AuthData: isSet(object.AuthData) ? AuthRequest.fromJSON(object.AuthData) : undefined,
    };
  },

  toJSON(message: SendMessageRequest): unknown {
    const obj: any = {};
    if (message.Msg !== "") {
      obj.Msg = message.Msg;
    }
    if (message.AuthData !== undefined) {
      obj.AuthData = AuthRequest.toJSON(message.AuthData);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SendMessageRequest>, I>>(base?: I): SendMessageRequest {
    return SendMessageRequest.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<SendMessageRequest>, I>>(object: I): SendMessageRequest {
    const message = createBaseSendMessageRequest();
    message.Msg = object.Msg ?? "";
    message.AuthData = (object.AuthData !== undefined && object.AuthData !== null)
      ? AuthRequest.fromPartial(object.AuthData)
      : undefined;
    return message;
  },
};

function createBaseChatMessage(): ChatMessage {
  return { UUID: "", UserUUID: "", UserName: "", Time: 0, Msg: "" };
}

export const ChatMessage = {
  encode(message: ChatMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.UUID !== "") {
      writer.uint32(10).string(message.UUID);
    }
    if (message.UserUUID !== "") {
      writer.uint32(18).string(message.UserUUID);
    }
    if (message.UserName !== "") {
      writer.uint32(26).string(message.UserName);
    }
    if (message.Time !== 0) {
      writer.uint32(32).int64(message.Time);
    }
    if (message.Msg !== "") {
      writer.uint32(42).string(message.Msg);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ChatMessage {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChatMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.UUID = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.UserUUID = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.UserName = reader.string();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.Time = longToNumber(reader.int64() as Long);
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.Msg = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ChatMessage {
    return {
      UUID: isSet(object.UUID) ? String(object.UUID) : "",
      UserUUID: isSet(object.UserUUID) ? String(object.UserUUID) : "",
      UserName: isSet(object.UserName) ? String(object.UserName) : "",
      Time: isSet(object.Time) ? Number(object.Time) : 0,
      Msg: isSet(object.Msg) ? String(object.Msg) : "",
    };
  },

  toJSON(message: ChatMessage): unknown {
    const obj: any = {};
    if (message.UUID !== "") {
      obj.UUID = message.UUID;
    }
    if (message.UserUUID !== "") {
      obj.UserUUID = message.UserUUID;
    }
    if (message.UserName !== "") {
      obj.UserName = message.UserName;
    }
    if (message.Time !== 0) {
      obj.Time = Math.round(message.Time);
    }
    if (message.Msg !== "") {
      obj.Msg = message.Msg;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ChatMessage>, I>>(base?: I): ChatMessage {
    return ChatMessage.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ChatMessage>, I>>(object: I): ChatMessage {
    const message = createBaseChatMessage();
    message.UUID = object.UUID ?? "";
    message.UserUUID = object.UserUUID ?? "";
    message.UserName = object.UserName ?? "";
    message.Time = object.Time ?? 0;
    message.Msg = object.Msg ?? "";
    return message;
  },
};

function createBaseHistoryResponse(): HistoryResponse {
  return { Messages: [] };
}

export const HistoryResponse = {
  encode(message: HistoryResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.Messages) {
      ChatMessage.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HistoryResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHistoryResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Messages.push(ChatMessage.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): HistoryResponse {
    return {
      Messages: Array.isArray(object?.Messages) ? object.Messages.map((e: any) => ChatMessage.fromJSON(e)) : [],
    };
  },

  toJSON(message: HistoryResponse): unknown {
    const obj: any = {};
    if (message.Messages?.length) {
      obj.Messages = message.Messages.map((e) => ChatMessage.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<HistoryResponse>, I>>(base?: I): HistoryResponse {
    return HistoryResponse.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<HistoryResponse>, I>>(object: I): HistoryResponse {
    const message = createBaseHistoryResponse();
    message.Messages = object.Messages?.map((e) => ChatMessage.fromPartial(e)) || [];
    return message;
  },
};

export interface Chat {
  GetHistory(request: DeepPartial<AuthRequest>, metadata?: grpc.Metadata): Promise<HistoryResponse>;
  SendMessage(request: DeepPartial<SendMessageRequest>, metadata?: grpc.Metadata): Promise<Empty>;
  ListenRequest(request: DeepPartial<AuthRequest>, metadata?: grpc.Metadata): Observable<ChatMessage>;
}

export class ChatClientImpl implements Chat {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.GetHistory = this.GetHistory.bind(this);
    this.SendMessage = this.SendMessage.bind(this);
    this.ListenRequest = this.ListenRequest.bind(this);
  }

  GetHistory(request: DeepPartial<AuthRequest>, metadata?: grpc.Metadata): Promise<HistoryResponse> {
    return this.rpc.unary(ChatGetHistoryDesc, AuthRequest.fromPartial(request), metadata);
  }

  SendMessage(request: DeepPartial<SendMessageRequest>, metadata?: grpc.Metadata): Promise<Empty> {
    return this.rpc.unary(ChatSendMessageDesc, SendMessageRequest.fromPartial(request), metadata);
  }

  ListenRequest(request: DeepPartial<AuthRequest>, metadata?: grpc.Metadata): Observable<ChatMessage> {
    return this.rpc.invoke(ChatListenRequestDesc, AuthRequest.fromPartial(request), metadata);
  }
}

export const ChatDesc = { serviceName: "Chat" };

export const ChatGetHistoryDesc: UnaryMethodDefinitionish = {
  methodName: "GetHistory",
  service: ChatDesc,
  requestStream: false,
  responseStream: false,
  requestType: {
    serializeBinary() {
      return AuthRequest.encode(this).finish();
    },
  } as any,
  responseType: {
    deserializeBinary(data: Uint8Array) {
      const value = HistoryResponse.decode(data);
      return {
        ...value,
        toObject() {
          return value;
        },
      };
    },
  } as any,
};

export const ChatSendMessageDesc: UnaryMethodDefinitionish = {
  methodName: "SendMessage",
  service: ChatDesc,
  requestStream: false,
  responseStream: false,
  requestType: {
    serializeBinary() {
      return SendMessageRequest.encode(this).finish();
    },
  } as any,
  responseType: {
    deserializeBinary(data: Uint8Array) {
      const value = Empty.decode(data);
      return {
        ...value,
        toObject() {
          return value;
        },
      };
    },
  } as any,
};

export const ChatListenRequestDesc: UnaryMethodDefinitionish = {
  methodName: "ListenRequest",
  service: ChatDesc,
  requestStream: false,
  responseStream: true,
  requestType: {
    serializeBinary() {
      return AuthRequest.encode(this).finish();
    },
  } as any,
  responseType: {
    deserializeBinary(data: Uint8Array) {
      const value = ChatMessage.decode(data);
      return {
        ...value,
        toObject() {
          return value;
        },
      };
    },
  } as any,
};

interface UnaryMethodDefinitionishR extends grpc.UnaryMethodDefinition<any, any> {
  requestStream: any;
  responseStream: any;
}

type UnaryMethodDefinitionish = UnaryMethodDefinitionishR;

interface Rpc {
  unary<T extends UnaryMethodDefinitionish>(
    methodDesc: T,
    request: any,
    metadata: grpc.Metadata | undefined,
  ): Promise<any>;
  invoke<T extends UnaryMethodDefinitionish>(
    methodDesc: T,
    request: any,
    metadata: grpc.Metadata | undefined,
  ): Observable<any>;
}

export class GrpcWebImpl {
  private host: string;
  private options: {
    transport?: grpc.TransportFactory;
    streamingTransport?: grpc.TransportFactory;
    debug?: boolean;
    metadata?: grpc.Metadata;
    upStreamRetryCodes?: number[];
  };

  constructor(
    host: string,
    options: {
      transport?: grpc.TransportFactory;
      streamingTransport?: grpc.TransportFactory;
      debug?: boolean;
      metadata?: grpc.Metadata;
      upStreamRetryCodes?: number[];
    },
  ) {
    this.host = host;
    this.options = options;
  }

  unary<T extends UnaryMethodDefinitionish>(
    methodDesc: T,
    _request: any,
    metadata: grpc.Metadata | undefined,
  ): Promise<any> {
    const request = { ..._request, ...methodDesc.requestType };
    const maybeCombinedMetadata = metadata && this.options.metadata
      ? new BrowserHeaders({ ...this.options?.metadata.headersMap, ...metadata?.headersMap })
      : metadata ?? this.options.metadata;
    return new Promise((resolve, reject) => {
      grpc.unary(methodDesc, {
        request,
        host: this.host,
        metadata: maybeCombinedMetadata ?? {},
        ...(this.options.transport !== undefined ? { transport: this.options.transport } : {}),
        debug: this.options.debug ?? false,
        onEnd: function (response) {
          if (response.status === grpc.Code.OK) {
            resolve(response.message!.toObject());
          } else {
            const err = new GrpcWebError(response.statusMessage, response.status, response.trailers);
            reject(err);
          }
        },
      });
    });
  }

  invoke<T extends UnaryMethodDefinitionish>(
    methodDesc: T,
    _request: any,
    metadata: grpc.Metadata | undefined,
  ): Observable<any> {
    const upStreamCodes = this.options.upStreamRetryCodes ?? [];
    const DEFAULT_TIMEOUT_TIME: number = 3_000;
    const request = { ..._request, ...methodDesc.requestType };
    const transport = this.options.streamingTransport ?? this.options.transport;
    const maybeCombinedMetadata = metadata && this.options.metadata
      ? new BrowserHeaders({ ...this.options?.metadata.headersMap, ...metadata?.headersMap })
      : metadata ?? this.options.metadata;
    return new Observable((observer) => {
      const upStream = (() => {
        const client = grpc.invoke(methodDesc, {
          host: this.host,
          request,
          ...(transport !== undefined ? { transport } : {}),
          metadata: maybeCombinedMetadata ?? {},
          debug: this.options.debug ?? false,
          onMessage: (next) => observer.next(next),
          onEnd: (code: grpc.Code, message: string, trailers: grpc.Metadata) => {
            if (code === 0) {
              observer.complete();
            } else if (upStreamCodes.includes(code)) {
              setTimeout(upStream, DEFAULT_TIMEOUT_TIME);
            } else {
              const err = new Error(message) as any;
              err.code = code;
              err.metadata = trailers;
              observer.error(err);
            }
          },
        });
        observer.add(() => {
          return client.close();
        });
      });
      upStream();
    }).pipe(share());
  }
}

declare const self: any | undefined;
declare const window: any | undefined;
declare const global: any | undefined;
const tsProtoGlobalThis: any = (() => {
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

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new tsProtoGlobalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

export class GrpcWebError extends tsProtoGlobalThis.Error {
  constructor(message: string, public code: grpc.Code, public metadata: grpc.Metadata) {
    super(message);
  }
}
