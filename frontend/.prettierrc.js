import baseConfig from "@migan/prettier-config";

/** @type {import("prettier").Config} */
export default {
	...baseConfig,
	plugins: [
		"prettier-plugin-tailwindcss",
		"@trivago/prettier-plugin-sort-imports",
	],
};
