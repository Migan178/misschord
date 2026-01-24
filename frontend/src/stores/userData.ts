import { create } from "zustand";

type State = "authorized" | "unauthorized" | "loading";

interface UserData {
	name: string;
	profile: string;
	handle: string;
	state: State;

	setName: (name: string) => void;
	setProfile: (name: string) => void;
	setHandle: (name: string) => void;
	setState: (state: State) => void;
}

const useUserDataStore = create<UserData>(set => ({
	name: "User",
	profile: "/defaults/default_profile.png",
	handle: "default_user",
	state: "loading",

	setName: name => set(() => ({ name })),
	setProfile: profile => set(() => ({ profile })),
	setHandle: handle => set(() => ({ handle })),
	setState: state => set(() => ({ state })),
}));

export default useUserDataStore;
