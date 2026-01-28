import type { ErrorResponse } from "../models/errors";
import axios from "axios";
import createAuthRefreshInterceptor from "axios-auth-refresh";

const api = axios.create({
	baseURL: import.meta.env.VITE_BACKEND_URL,
	timeout: 5_000,
	withCredentials: true,
	headers: {
		"Content-Type": "application/json",
	},
});

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function refreshToken(_failedRequest: unknown) {
	return api
		.post("/users/refresh")
		.then(Promise.resolve)
		.catch(Promise.reject);
}

createAuthRefreshInterceptor(api, refreshToken, {
	shouldRefresh: err =>
		err.response?.status == 401 &&
		(err.response?.data as ErrorResponse).error === "token is expired",
});

export default api;
