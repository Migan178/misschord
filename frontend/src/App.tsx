import { useEffect, useRef, useState, type FormEvent } from "react";
import { Link } from "react-router";
import useWebSocket, { ReadyState } from "react-use-websocket";

export default function App() {
	const websocketServerUrl = "ws://localhost:8080/ws";

	const { sendMessage, lastMessage, readyState } = useWebSocket(
		websocketServerUrl,
		{
			shouldReconnect: () => true,
			onOpen: () => console.log("connected"),
			onError: console.error,
		},
	);

	const [messages, setMessages] = useState<string[]>([]);
	const inputRef = useRef<HTMLInputElement>(null);

	const connectionStatus = {
		[ReadyState.CONNECTING]: "연결 중",
		[ReadyState.OPEN]: "연결 됨",
		[ReadyState.CLOSING]: "연결 해제 중",
		[ReadyState.CLOSED]: "연결 해제 됨",
		[ReadyState.UNINSTANTIATED]: "인스턴스화 안 됨",
	}[readyState];

	useEffect(() => {
		if (lastMessage)
			// eslint-disable-next-line react-hooks/set-state-in-effect
			setMessages(prev => [...prev, lastMessage.data]);
	}, [lastMessage]);

	function handleSubmit(e: FormEvent<HTMLFormElement>) {
		e.preventDefault();

		sendMessage(inputRef.current!.value);
	}

	return (
		<div>
			<h1 className="font-bold text-2xl">Hello, World!</h1>
			<Link to="/asdf">test react router</Link>
			<hr />
			<h1>WEBSOCKET TEST</h1>
			<p>{connectionStatus}</p>
			<ul>
				{messages.map(msg => (
					<li key={msg}>
						<p>{msg}</p>
					</li>
				))}
			</ul>
			<form onSubmit={handleSubmit}>
				<input type="text" ref={inputRef} required />
				<input type="submit" value="보내기" />
			</form>
		</div>
	);
}
