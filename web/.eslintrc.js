module.exports = {
    'env': {
        'browser': true,
        'es2021': true,
    },
    'extends': [
        'google',
    ],
    'parserOptions': {
        'ecmaVersion': 12,
        'sourceType': 'module',
    },
    'rules': {
        'indent': ['error', 4],
        'require-jsdoc': 'off',
        'max-len': ['error', 100],
    },
};
