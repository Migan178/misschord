import useUserDataStore from "../../stores/userData";
import { Navigate, Outlet } from "react-router";

export default function PublicLayout() {
	const state = useUserDataStore(state => state.state);

	if (state === "loading") return null;

	if (state === "authorized") return <Navigate to={"/"} />;

	return <Outlet />;
}
