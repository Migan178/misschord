import api from "./api/axios";
import PrivateLayout from "./components/layouts/PrivateLayout";
import PublicLayout from "./components/layouts/PublicLayout";
import type { UserResponse } from "./models/user";
import Asdf from "./routes/Asdf";
import Login from "./routes/Login";
import Root from "./routes/Root";
import useUserDataStore from "./stores/userData";
import { useEffect } from "react";
import { BrowserRouter, Routes, Route } from "react-router";

export default function App() {
	const setHandle = useUserDataStore(state => state.setHandle);
	const setProfile = useUserDataStore(state => state.setProfile);
	const setState = useUserDataStore(state => state.setState);

	useEffect(() => {
		(async () => {
			try {
				const res = await api.get<UserResponse>("/users/me");

				setHandle(res.data.handle);
				setProfile(res.data.profile);
				setState("authorized");
			} catch (err) {
				console.error(err);
				setState("unauthorized");
			}
		})();
	}, [setHandle, setProfile, setState]);

	return (
		<BrowserRouter>
			<Routes>
				<Route element={<PrivateLayout />}>
					<Route path="/" element={<Root />} />
				</Route>
				<Route element={<PublicLayout />}>
					<Route path="/login" element={<Login />} />
				</Route>
				<Route path="/asdf" element={<Asdf />} />
			</Routes>
		</BrowserRouter>
	);
}
