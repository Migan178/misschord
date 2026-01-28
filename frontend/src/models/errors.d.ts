export interface ErrorResponse {
	error: string;
}

export interface ErrorsResponse {
	errors: string[];
}

export interface WebsocketErrorData {
	error: string;
	code: number;
}
