import type { WebsocketData } from "../models/websocket/data";
import { createContext, useContext } from "react";
import type { SendMessage } from "react-use-websocket";

export type ConnectionStatus =
	| "연결 중"
	| "연결 됨"
	| "연결 해제 중"
	| "연결 해제 됨"
	| "인스턴스화 안 됨";

export interface SocketContextType {
	sendMessage: SendMessage;
	lastMessage: WebsocketData;
	connectionStatus: ConnectionStatus;
}

export const SocketContext = createContext<SocketContextType | null>(null);

export function useSocket() {
	const context = useContext(SocketContext);

	if (!context)
		throw new Error("useSocket은 SocketProvider내에서 써라 애송이");

	return context;
}
