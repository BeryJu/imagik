import resolve from "@rollup/plugin-node-resolve";
import commonjs from "@rollup/plugin-commonjs";
import cssimport from "rollup-plugin-cssimport";
import copy from "rollup-plugin-copy";
import minifyHTML from "rollup-plugin-minify-html-literals";
import { terser } from "rollup-plugin-terser";
import sourcemaps from "rollup-plugin-sourcemaps";
import image from "@rollup/plugin-image";

const resources = [
    { src: "src/index.html", dest: "../root" },
    { src: "src/style.css", dest: "../root" },
    { src: "assets/*", dest: "../root/assets" },
];

module.exports = [
    {
        input: "./src/main.js",
        output: [
            {
                format: "esm",
                file: "../root/bundle.js",
                sourcemap: true,
            },
        ],
        plugins: [
            resolve({ browser: true }),
            commonjs(),
            cssimport(),
            image({ dom: true }),
            sourcemaps(),
            process.env.NODE_ENV === "production" && minifyHTML(),
            process.env.NODE_ENV === "production" && terser(),
            copy({
                targets: [...resources],
                copyOnce: false,
            }),
        ],

        watch: {
            clearScreen: false,
        },
    },
];
