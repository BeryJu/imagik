import { LitElement, html, css } from "lit";
import "./ik-header.js";
import "./ik-app.js";
import { isLoggedIn, request, get } from "./services/api.js";
import { until } from "lit-html/directives/until.js";

class Gate extends LitElement {
    static get styles() {
        return css`
            :host {
                display: block;
            }
            div {
                display: flex;
                justify-content: center;
                align-items: center;
                height: 60vh;
            }
            form {
                display: flex;
                flex-direction: column;
                justify-content: space-between;
                width: 30rem;
                margin: auto;
                padding: 3rem;
                background-color: var(--color-primary-background-dark);
                box-shadow: 0px 2px 3px 0px #0008;
            }
            form > input,
            form > button {
                margin: 1rem 0;
                padding: 0 1rem;
                line-height: 3rem;
                background-color: var(--color-primary-background-light);
                outline: none;
                border: 0;
                color: var(--color-primary-text);
                transition: 0.25s border-bottom;
                transition: 0.25s background-color;
                border-bottom: 1px solid rgba(0, 0, 0, 0);
            }
            input:focus,
            textarea:focus,
            select:focus {
                outline: none;
                border-bottom: 1px solid var(--color-primary);
            }
            button:hover,
            input[type="submit"]:hover {
                background-color: var(--color-primary);
            }
        `;
    }

    submitLogin(ev) {
        ev.preventDefault();
        const elements = ev.submitter.form.elements;
        request("/api/pub/auth/login", null, {
            method: "POST",
            headers: {
                authorization:
                    "Basic " +
                    btoa(
                        `${elements.namedItem("username").value}:${
                            elements.namedItem("password").value
                        }`,
                    ),
            },
        });
        this.requestUpdate();
    }

    render() {
        if (isLoggedIn()) {
            return html`<ik-app></ik-app> `;
        } else {
            return html`
                <ik-header></ik-header>
                ${until(
                    get("/api/pub/auth/driver").then((res) => {
                        if (res.type === "static") {
                            return html`<div>
                                <form @submit=${this.submitLogin}>
                                    <input
                                        type="text"
                                        placeholder="username"
                                        name="username"
                                        required
                                    />
                                    <input
                                        type="password"
                                        placeholder="password"
                                        name="password"
                                        required
                                    />
                                    <input type="submit" value="login" />
                                </form>
                            </div>`;
                        }
                    }),
                )}
            `;
        }
    }
}
customElements.define("ik-gate", Gate);
