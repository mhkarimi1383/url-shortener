/* eslint-env node */
require('@rushstack/eslint-patch/modern-module-resolution');

module.exports = {
  root: true,
  plugins: ['waterfall'],
    parserOptions: {
    ecmaVersion: 'latest',
    extraFileExtensions: ['.vue'],
    parser: require.resolve('@typescript-eslint/parser'),
  },
    rules: {
    'waterfall/waterfall-objects': 'error',
    'waterfall/waterfall-imports': 'error',
    'waterfall/waterfall-requires': 'error',
    'no-unused-vars': ['error', { argsIgnorePattern: '^_$' }],
  },
    extends: [
    'plugin:vue/vue3-essential',
    'eslint:recommended',
    // '@vue/eslint-config-typescript',
    '@vue/eslint-config-prettier/skip-formatting',
    'plugin:@typescript-eslint/eslint-recommended',
  ],
};
