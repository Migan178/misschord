import api from "../api/axios";
import { useSocket } from "../contexts/socket";
import { OPCode } from "../models/websocket/data";
import {
	ChannelType,
	EventType,
	isMessageCreate,
} from "../models/websocket/events";
import { useEffect, useState } from "react";
import { Link } from "react-router";

export default function Root() {
	const { sendMessage, lastMessage, connectionStatus } = useSocket();

	const [messages, setMessages] = useState<string[]>([]);
	const [channel, setChannel] = useState(0);

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

	async function handleSetChannel(data: FormData) {
		if (!data.get("id")) return;

		if (channel > 0)
			sendMessage(
				JSON.stringify({
					op: OPCode.Dispatch,
					data: {
						id: channel,
						type: ChannelType.DM,
					},
					type: EventType.ChannelLeave,
				}),
			);

		setChannel(Number(data.get("id")?.toString()));

		try {
			await api.post("/users/me/channels", {
				recipient_id: channel,
			});

			sendMessage(
				JSON.stringify({
					op: OPCode.Dispatch,
					data: {
						id: channel,
						type: ChannelType.DM,
					},
					type: EventType.ChannelJoin,
				}),
			);
		} catch (err) {
			console.error(err);
		}
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
				<input type="number" name="id" required />
				<input type="submit" value="채널 설정" />
			</form>
			<form action={handleSendMessage}>
				<input type="text" name="text" required />
				<input type="submit" value="보내기" />
			</form>
		</div>
	);
}
