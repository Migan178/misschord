import useUserDataStore from "../../stores/userData";
import { Navigate, Outlet } from "react-router";

export default function PrivateLayout() {
	const state = useUserDataStore(state => state.state);

	if (state === "loading") return null;

	if (state === "unauthorized") return <Navigate to="/login" replace />;

	return <Outlet />;
}
