import { LitElement, html, css } from "lit";
import { until } from "lit-html/directives/until.js";
import { get } from "./services/api.js";

import "./ik-file.js";
import "./ik-directory.js";
import "./ik-detail.js";

class List extends LitElement {
    static get styles() {
        return css`
            :host {
                display: flex;
                flex-direction: row;
                flex-wrap: wrap;
                padding: 1rem;
                gap: 1rem;
            }
        `;
    }

    static get properties() {
        return {
            path: {
                attribute: true,
                type: String,
            },
        };
    }

    async listFiles(path) {
        return get(`./list?pathOffset=${encodeURIComponent(path)}`)
            .then(({ files }) => {
                if (files.length === 0) {
                    throw new Exception("empty files");
                }
                return files
                    .sort((a, b) => {
                        if (a.type !== b.type) {
                            return a.type.localeCompare(b.type);
                        } else {
                            return a.name.localeCompare(b.name);
                        }
                    })
                    .map(
                        (f) => html`
                            <ik-file path=${f.fullPath} mime=${f.mime} name=${f.name}> </ik-file>
                        `,
                    );
            })
            .catch((e) => {
                return html`<ik-detail path=${path}></ik-detail>`;
            });
    }

    render() {
        return html` ${until(this.listFiles(this.path), html`loading...`)} `;
    }
}
customElements.define("ik-list", List);
