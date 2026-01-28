import type { UserResponse } from "../user";
import { OPCode, type WebsocketData } from "./data";

export const EventType = {
	MessageCreate: "MESSAGE_CREATE",
	ChannelJoin: "CHANNEL_JOIN",
	ChannelLeave: "CHANNEL_LEAVE",
} as const;
export type EventType = (typeof EventType)[keyof typeof EventType];

export const ChannelType = {
	DM: "DM",
} as const;
export type ChannelType = (typeof ChannelType)[keyof typeof ChannelType];

export interface MessageCreateEvent {
	author: UserResponse;
	message: string;
	channelType: ChannelType;
	createdAt: string;
}

export interface WebsocketDispatchData extends WebsocketData {
	op: 0;
}

export interface WebsocketMessageCreateData extends WebsocketDispatchData {
	data: MessageCreateEvent;
}

export function isMessageCreate(
	msg: WebsocketData,
): msg is WebsocketMessageCreateData {
	return msg.op === OPCode.Dispatch && msg.type === EventType.MessageCreate;
}

export interface ChannelData {
	id: string;
	type: ChannelType;
}

export interface WebsocketChannelData extends WebsocketDispatchData {
	data: ChannelData;
}

export function isChannelJoin(msg: WebsocketData): msg is WebsocketChannelData {
	return msg.op === OPCode.Dispatch && msg.type == EventType.ChannelJoin;
}
