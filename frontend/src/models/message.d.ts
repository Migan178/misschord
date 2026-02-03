export interface MessageResponse {
	author: UserResponse;
	message: string;
	channelType: ChannelType;
	createdAt: string;
}
