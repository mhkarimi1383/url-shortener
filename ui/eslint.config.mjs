import eslint from '@eslint/js';
import eslintConfigPrettier from 'eslint-config-prettier';
import eslintPluginVue from 'eslint-plugin-vue';
import globals from 'globals';
import typescriptEslint from 'typescript-eslint';
import waterfall from 'eslint-plugin-waterfall';

export default typescriptEslint.config(
  { ignores: ['*.d.ts', '**/coverage', '**/dist'] },
  {
    extends: [
      eslint.configs.recommended,
      ...typescriptEslint.configs.recommended,
      ...eslintPluginVue.configs['flat/recommended'],
    ],
    files: ['**/*.{ts,vue}'],
    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
      globals: globals.browser,
      parserOptions: {
        parser: typescriptEslint.parser,
      },
    },
    plugins: {
      waterfall: waterfall,
    },
    rules: {
        "waterfall/waterfall-objects": "error",
        "waterfall/waterfall-imports": "error",
        "waterfall/waterfall-requires": "error",
        "@typescript-eslint/no-explicit-any": ["off"],
        "no-unused-vars": ["error", {
            argsIgnorePattern: "^_$",
        }],
    },
  },
  eslintConfigPrettier
);
