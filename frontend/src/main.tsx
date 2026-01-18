import App from "./App.tsx";
import Asdf from "./Asdf.tsx";
import "./index.css";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Route, Routes } from "react-router";

createRoot(document.getElementById("root")!).render(
	<StrictMode>
		<BrowserRouter>
			<Routes>
				<Route path="/" element={<App />} />
				<Route path="/asdf" element={<Asdf />} />
			</Routes>
		</BrowserRouter>
	</StrictMode>,
);
