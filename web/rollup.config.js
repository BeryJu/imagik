import resolve from "@rollup/plugin-node-resolve";
import commonjs from "@rollup/plugin-commonjs";
import cssimport from "rollup-plugin-cssimport";
import copy from "rollup-plugin-copy";
import minifyHTML from "rollup-plugin-minify-html-literals-v3";
import terser from "@rollup/plugin-terser";
import image from "@rollup/plugin-image";

const resources = [
    { src: "src/style.css", dest: "./dist" },
    { src: "assets/*", dest: "./dist/assets" },
];

export default [
    {
        input: "./src/main.js",
        output: [
            {
                format: "esm",
                file: "./dist/bundle.js",
                sourcemap: true,
            },
        ],
        plugins: [
            resolve({ browser: true }),
            commonjs(),
            cssimport(),
            image({ dom: true }),
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
