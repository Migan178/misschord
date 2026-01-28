import { useSocket } from "../contexts/socket";
import { OPCode } from "../models/websocket/data";
import { EventType, isMessageCreate } from "../models/websocket/events";
import { useEffect, useState } from "react";
import { Link } from "react-router";

export default function Root() {
	const { sendMessage, lastMessage, connectionStatus } = useSocket();

	const [messages, setMessages] = useState<string[]>([]);
	const [channel, setChannel] = useState("");

	function handleSendMessage(data: FormData) {
		sendMessage(
			JSON.stringify({
				op: OPCode.Dispatch,
				data: {
					message: data.get("text"),
					channel: {
						id: channel,
						type: "DM",
					},
				},
				type: EventType.MessageCreate,
			}),
		);
	}

	function handleSetChannel(data: FormData) {
		setChannel(data.get("id")!.toString());

		sendMessage(
			JSON.stringify({
				op: 0,
				data: {
					id: channel,
					type: "DM",
				},
				type: EventType.ChannelJoin,
			}),
		);
	}

	useEffect(() => {
		console.log(lastMessage);
		if (!lastMessage || !isMessageCreate(lastMessage)) return;
		// eslint-disable-next-line react-hooks/set-state-in-effect
		setMessages(prev => [
			...prev,
			`${lastMessage.data.author.handle}: ${lastMessage.data.message}`,
		]);
	}, [lastMessage]);

	return (
		<div>
			<h1 className="font-bold text-2xl">Hello, World!</h1>
			<Link to="/asdf">test react router</Link>
			<hr />
			<h1>WEBSOCKET TEST</h1>
			<p>{connectionStatus}</p>
			<ul>
				{messages.map(msg => (
					<li>
						<p>{msg}</p>
					</li>
				))}
			</ul>
			<form action={handleSetChannel}>
				<input type="text" name="id" required />
				<input type="submit" value="채널 설정" />
			</form>
			<form action={handleSendMessage}>
				<input type="text" name="text" required />
				<input type="submit" value="보내기" />
			</form>
		</div>
	);
}
