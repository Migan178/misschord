import type { UserResponse } from "../user";
import type { EventType } from "./events";

export const OPCode = {
	Dispatch: 0,
	HeartBeat: 1,
	HeartBeatACK: 2,
	Hello: 3,
	Identify: 4,
	Ready: 5,
	Error: 6,
} as const;
export type OPCode = (typeof OPCode)[keyof typeof OPCode];

export interface IdentifyData {
	token: string;
}

export interface HelloData {
	heartbeatInterval: number;
	message?: string;
}

export interface ReadyData {
	user: UserResponse;
}

export interface WebsocketData {
	op: OPCode;
	data?: unknown;
	type?: EventType;
}

export interface WebsocketHelloData extends WebsocketData {
	op: 3;
	data: HelloData;
}

export function isHello(msg: WebsocketData): msg is WebsocketHelloData {
	return msg.op === OPCode.Hello;
}

export interface WebsocketReadyData extends WebsocketData {
	op: 5;
	data: ReadyData;
}

export function isReady(msg: WebsocketData): msg is WebsocketReadyData {
	return msg.op === OPCode.Ready;
}
