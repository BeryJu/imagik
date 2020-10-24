/* eslint-disable import/no-extraneous-dependencies */
const { createDefaultConfig } = require('@open-wc/testing-karma');
const merge = require('deepmerge');

const rollupConfig = require('./rollup.config.js')[0];

rollupConfig.output = {
  format: 'esm',
  sourcemap: 'inline',
}
rollupConfig.external = ['sinon', 'chai', '@open-wc/testing'],

module.exports = config => {
  config.set(
    merge(createDefaultConfig(config), {
      browsers : ['FirefoxHeadless'], // chrome is in default config

      files: [
        // runs all files ending with .test in the test folder,
        // can be overwritten by passing a --grep flag. examples:
        //
        // npm run test -- --grep test/foo/bar.test.js
        // npm run test -- --grep test/bar/*
        { pattern: config.grep ? config.grep : 'test/**/*.test.js', type: 'module', watched: false },
      ],

      preprocessors: {
          'test/**/*.test.js': ['rollup'],
      },

      rollupPreprocessor: rollupConfig,

      esm: {
        // if you are using 'bare module imports' you will need this option
        nodeResolve: true,
        coverageExclude: ['assets/codemirror/**/*.js'],
      },

      reporters: ['junit'],

      junitReporter: {
        outputDir: 'junit',
      },

      coverageReporter: {
        reporters: [
          { type: 'cobertura', subdir: '.', file: 'cobertura.xml' },
        ]
      }
    }),
  );
  return config;
};
