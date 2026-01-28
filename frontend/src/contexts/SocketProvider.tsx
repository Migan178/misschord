import { isHello, OPCode, type WebsocketData } from "../models/websocket/data";
import {
	type ConnectionStatus,
	type SocketContextType,
	SocketContext,
} from "./socket";
import { useEffect, useState, type ReactNode } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";

export default function SocketProvider({ children }: { children: ReactNode }) {
	const [heartbeatInterval, setHeartbeatInterval] = useState<number | null>(
		null,
	);

	const { readyState, lastJsonMessage, sendMessage } =
		useWebSocket<WebsocketData>("/api/v1/ws", {
			shouldReconnect: () => true,
			onOpen: () => console.log("connected"),
			onError: console.error,
			share: true,
		});

	const connectionStatus = {
		[ReadyState.CONNECTING]: "연결 중",
		[ReadyState.OPEN]: "연결 됨",
		[ReadyState.CLOSING]: "연결 해제 중",
		[ReadyState.CLOSED]: "연결 해제 됨",
		[ReadyState.UNINSTANTIATED]: "인스턴스화 안 됨",
	}[readyState] as ConnectionStatus;

	const value: SocketContextType = {
		sendMessage,
		lastMessage: lastJsonMessage,
		connectionStatus,
	};

	useEffect(() => {
		if (lastJsonMessage && isHello(lastJsonMessage)) {
			// eslint-disable-next-line react-hooks/set-state-in-effect
			setHeartbeatInterval(lastJsonMessage.data.heartbeatInterval);
		}
	}, [lastJsonMessage]);

	useEffect(() => {
		if (heartbeatInterval && readyState == ReadyState.OPEN) {
			const timer = setInterval(() => {
				sendMessage(JSON.stringify({ op: OPCode.HeartBeat }));
			}, heartbeatInterval);

			return () => clearInterval(timer);
		}
	}, [heartbeatInterval, readyState, sendMessage]);

	return (
		<SocketContext.Provider value={value}>
			{children}
		</SocketContext.Provider>
	);
}
