import js from "@eslint/js";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { FlatCompat } from "@eslint/eslintrc";
import waterfall from "eslint-plugin-waterfall";
import { includeIgnoreFile } from "@eslint/compat"

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const compat = new FlatCompat({
    baseDirectory: __dirname,
    allConfig: js.configs.all,
    recommendedConfig: js.configs.recommended
});
const gitignorePath = path.resolve(__dirname, ".gitignore");

export default [...compat.extends(
    "plugin:vue/vue3-essential",
    "eslint:recommended",
    "@vue/eslint-config-prettier/skip-formatting",
    "plugin:@typescript-eslint/eslint-recommended",
), {

        plugins: {
        waterfall,
    },
        rules: {
        "waterfall/waterfall-objects": "error",
        "waterfall/waterfall-imports": "error",
        "waterfall/waterfall-requires": "error",

        "no-unused-vars": ["error", {
            argsIgnorePattern: "^_$",
        }],
    },
        languageOptions: {
        sourceType: "script",
        ecmaVersion: "latest",
                parserOptions: {
            extraFileExtensions: [".vue"],
            parser: "/home/karimi/go/src/github.com/mhkarimi1383/url-shortener/ui/node_modules/@typescript-eslint/parser/dist/index.js",
        },
    },
},
includeIgnoreFile(gitignorePath),
];
