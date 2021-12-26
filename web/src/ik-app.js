import { LitElement, html, css } from "lit-element";
import "./ik-header.js";
import "./ik-drop.js";
import "./ik-list.js";
import { logout, request } from "./services/api.js";

class App extends LitElement {
    static get styles() {
        return css`
            :host {
                display: block;
            }
            ik-header a,
            ik-header a:visited {
                color: var(--color-primary);
            }
        `;
    }

    static get properties() {
        return {
            dragover: {
                attribute: true,
                type: Boolean,
            },
            path: {
                attribute: true,
                type: String,
                reflect: true,
            },
        };
    }

    constructor() {
        super();
        this.addEventListener(
            "dragover",
            (ev) => {
                ev.preventDefault();
                this.dragover = true;
            },
            false,
        );
        this.addEventListener(
            "dragleave",
            (ev) => {
                ev.preventDefault();
                this.dragover = false;
            },
            false,
        );
        this.addEventListener("drop", (ev) => {
            ev.preventDefault();
            this.dragover = false;
            this.handleDrop(ev);
        });
        this.addEventListener("update", (ev) => {
            ev.preventDefault();
            this.shadowRoot.querySelector("ik-list").requestUpdate();
        });

        this.navigate({ detail: window.location.hash.slice(1, Infinity) || "/" });
    }

    connectedCallback() {
        super.connectedCallback();
        window.addEventListener("hashchange", () => {
            this.path = window.location.hash.slice(1, Infinity);
        });
    }

    handleDrop(ev) {
        // Prevent default behavior (Prevent file from being opened)
        ev.preventDefault();

        if (ev.dataTransfer.items) {
            // Use DataTransferItemList interface to access the file(s)
            for (const item of ev.dataTransfer.items) {
                // If dropped items aren't files, reject them
                if (item.kind === "file") {
                    const file = item.getAsFile();
                    this.uploadFile(file);
                } else {
                    console.warn("... " + item.kind);
                }
            }
        }
    }

    uploadSelect(ev) {
        const input = document.createElement("input");
        input.setAttribute("type", "file");
        input.setAttribute("multiple", "true");
        input.addEventListener("change", (ev) => {
            const files = ev.target.files;
            for (const file of files) {
                this.uploadFile(file);
            }
        });
        input.click();
    }

    uploadFile(file) {
        request(`${this.path}/${file.name}`, file, {
            method: "PUT",
        })
            .catch((e) => {
                console.error(e);
            })
            .then((r) => {
                this.triggerUpdate();
            });
    }

    triggerUpdate() {
        this.dispatchEvent(
            new CustomEvent("update", {
                composed: true,
                bubbles: true,
            }),
        );
    }

    navigate({ detail }) {
        this.path = detail;
        if (detail == "/") {
            document.title = "imagik";
        } else {
            document.title = `imagik - ${detail}`;
        }
    }

    render() {
        if (window.location.hash !== "#" + this.path) window.location.hash = "#" + this.path;

        return html`
            <ik-header path=${this.path} @navigate=${(e) => this.navigate(e)}>
                <a @click=${() => this.uploadSelect()}>upload</a>
                |
                <a @click=${() => this.triggerUpdate()}>refresh</a>
                |
                <a @click=${logout}>logout</a>
            </ik-header>

            <ik-list path=${this.path} @navigate=${(e) => this.navigate(e)}></ik-list>

            <ik-drop ?show=${this.dragover}></ik-drop>
        `;
    }
}
customElements.define("ik-app", App);
