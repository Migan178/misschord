import api from "../api/axios";
import { useSocket } from "../contexts/socket";
import type { MessageResponse } from "../models/message";
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

	async function handleSendMessage(data: FormData) {
		try {
			await api.post(`/users/me/channels/${channel}/messages`, {
				message: data.get("text"),
			});
		} catch (err) {
			console.error(err);
		}
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

			const res = await api.get<MessageResponse[]>(
				`/users/me/channels/${channel}/messages`,
			);

			for (const data of res.data) {
				setMessages(prev => [
					...prev,
					`${data.author.handle}: ${data.message}`,
				]);
			}
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
