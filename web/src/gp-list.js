import {LitElement, html, css} from 'lit-element';
import {until} from 'lit-html/directives/until.js';
import {get} from './services/api.js';

import './gp-file.js';
import './gp-directory.js';

class GpList extends LitElement {
    static get styles() {
        return css`
            :host {
                display: flex;
                flex-direction: row;
                flex-wrap: wrap;
                justify-content: space-evenly;
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
            .then(({files})=>files.map((f)=>html`
                <gp-file path=${f.fullPath} mime=${f.mime}>
                    ${f.name}
                </gp-file>
            `));
    }

    render() {
        return html`
            ${until(this.listFiles(this.path), html`loading...`)}
        `;
    }
}
customElements.define('gp-list', GpList);
