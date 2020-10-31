import {LitElement, html, css} from 'lit-element';
import {until} from 'lit-html/directives/until.js';
import {get} from './services/api.js';

import './gp-file.js';
import './gp-directory.js';
import './gp-detail.js';

class GpList extends LitElement {
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
            .then(({ files }) => files.sort((a, b) => {
                if (a.type !== b.type) {
                    return a.type.localeCompare(b.type);
                } else {
                    return a.name.localeCompare(b.name);
                }
            })
            .map((f) => html`
                <gp-file path=${f.fullPath} mime=${f.mime} name=${f.name}>
                </gp-file>
            `))
            .catch(e => {
                return html`<gp-detail path=${path}></gp-detail>`;
            });
    }

    render() {
        return html`
            ${until(this.listFiles(this.path), html`loading...`)}
        `;
    }
}
customElements.define('gp-list', GpList);
