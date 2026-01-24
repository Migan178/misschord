import api from "../api/axios";
import { useState, type ChangeEvent, type FormEvent } from "react";
import { useNavigate } from "react-router";

export default function Login() {
	const navigate = useNavigate();

	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");

	function handleEmailChange(e: ChangeEvent<HTMLInputElement>) {
		setEmail(e.target.value);
	}

	function handlePasswordChange(e: ChangeEvent<HTMLInputElement>) {
		setPassword(e.target.value);
	}

	async function handleSubmit(e: FormEvent<HTMLFormElement>) {
		e.preventDefault();

		try {
			await api.post("/users/login", {
				email,
				password,
			});

			navigate("/");
		} catch (err) {
			console.error(err);
			alert("에러");
		}
	}

	return (
		<form onSubmit={handleSubmit}>
			<div>
				<label>이메일</label>
				<input
					type="email"
					name="email"
					value={email}
					onChange={handleEmailChange}
				/>
			</div>
			<div>
				<label>비밀번호</label>
				<input
					type="password"
					name="password"
					value={password}
					onChange={handlePasswordChange}
				/>
			</div>
			<button type="submit">로그인</button>
		</form>
	);
}
