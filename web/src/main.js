import "construct-style-sheets-polyfill";
import "./ik-gate.js";
import * as Sentry from "@sentry/browser";
import { Integrations as TracingIntegrations } from "@sentry/tracing";
Sentry.init({
    dsn: "https://759fc52c64334ea0a2978460a6469fd0@sentry.beryju.org/15",
    integrations: [new TracingIntegrations.BrowserTracing()],
    tracesSampleRate: 1,
    tunnel: "/api/pub/sentry",
});
