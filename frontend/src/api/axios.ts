import type { ErrorResponse } from "../models/errors";
import axios, { type AxiosError } from "axios";

const api = axios.create({
	baseURL: import.meta.env.VITE_BACKEND_URL,
	timeout: 5_000,
	withCredentials: true,
	headers: {
		"Content-Type": "application/json",
	},
});

api.interceptors.response.use(
	response => {
		if (response.data.createdAt) {
			response.data.createdAt = new Date(response.data.createdAt);
		}

		if (response.data.expiresAt) {
			response.data.expiresAt = new Date(response.data.expireAt);
		}

		return response;
	},
	async (error: AxiosError) => {
		const originalReq = error.config;
		if (!originalReq) return Promise.reject(error);

		if (error.response) {
			if (
				error.response.status === 401 &&
				typeof error.response.data === "object"
			) {
				const data = error.response.data as ErrorResponse;

				if (data.error === "token is invalid") {
					window.location.href = "/login";
					return;
				}

				try {
					await api.post("/users/refresh");

					return api(originalReq);
				} catch (err) {
					console.error(err);

					window.location.href = "/login";
				}
			}
		}

		return Promise.reject(error);
	},
);

export default api;
