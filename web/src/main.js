import "construct-style-sheets-polyfill";
import "./ik-gate.js";
import * as Sentry from "@sentry/browser";
Sentry.init({
    dsn: "https://bc5df9f742f04f62bb1e4b37b227aa62@sentry.beryju.org/7",
    tunnel: "/api/pub/sentry",
});
