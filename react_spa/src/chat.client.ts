// @generated by protobuf-ts 2.9.0
// @generated from protobuf file "chat.proto" (syntax proto3)
// tslint:disable
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { Chat } from "./chat";
import type { ChatMessage } from "./chat";
import type { ServerStreamingCall } from "@protobuf-ts/runtime-rpc";
import type { Empty } from "./chat";
import type { SendMessageRequest } from "./chat";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { HistoryResponse } from "./chat";
import type { AuthRequest } from "./chat";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
/**
 * @generated from protobuf service Chat
 */
export interface IChatClient {
    /**
     * @generated from protobuf rpc: GetHistory(AuthRequest) returns (HistoryResponse);
     */
    getHistory(input: AuthRequest, options?: RpcOptions): UnaryCall<AuthRequest, HistoryResponse>;
    /**
     * @generated from protobuf rpc: SendMessage(SendMessageRequest) returns (Empty);
     */
    sendMessage(input: SendMessageRequest, options?: RpcOptions): UnaryCall<SendMessageRequest, Empty>;
    /**
     * @generated from protobuf rpc: ListenRequest(AuthRequest) returns (stream ChatMessage);
     */
    listenRequest(input: AuthRequest, options?: RpcOptions): ServerStreamingCall<AuthRequest, ChatMessage>;
}
/**
 * @generated from protobuf service Chat
 */
export class ChatClient implements IChatClient, ServiceInfo {
    typeName = Chat.typeName;
    methods = Chat.methods;
    options = Chat.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @generated from protobuf rpc: GetHistory(AuthRequest) returns (HistoryResponse);
     */
    getHistory(input: AuthRequest, options?: RpcOptions): UnaryCall<AuthRequest, HistoryResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<AuthRequest, HistoryResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: SendMessage(SendMessageRequest) returns (Empty);
     */
    sendMessage(input: SendMessageRequest, options?: RpcOptions): UnaryCall<SendMessageRequest, Empty> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<SendMessageRequest, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: ListenRequest(AuthRequest) returns (stream ChatMessage);
     */
    listenRequest(input: AuthRequest, options?: RpcOptions): ServerStreamingCall<AuthRequest, ChatMessage> {
        const method = this.methods[2], opt = this._transport.mergeOptions(options);
        return stackIntercept<AuthRequest, ChatMessage>("serverStreaming", this._transport, method, opt, input);
    }
}
