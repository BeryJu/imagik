import {LitElement, html, css} from 'lit-element';
import { until } from 'lit-html/directives/until.js';
import { get } from './services/api.js';

class Detail extends LitElement {
    static get styles() {
        return css`
            :host {
                display: flex;
                height: 100%;
            }
            div.sidebar {
                width: 15rem;
                background-color: var(--color-primary-background-dark);
                color: var(--color-primary-text);
                box-shadow: 0px 2px 3px 0px #0008;
                padding: 1rem;
                display: flex;
                flex-direction: column;
                margin-right: 1rem;
            }
            div.sidebar > div.group {
                display: flex;
                flex-direction: column;
            }
            input, button {
                margin: 1rem 0;
                padding: 0 1rem;
                line-height: 3rem;
                background-color: var(--color-primary-background-light);
                outline: none;
                border: 0;
                color: var(--color-primary-text);
                transition: 0.25s border-bottom;
                transition: 0.25s background-color;
                border-bottom: 1px solid rgba(0,0,0,0);
            }
            input:focus, textarea:focus, select:focus{
                outline: none;
                border-bottom: 1px solid var(--color-primary);
            }
            button:hover, input[type=submit]:hover {
                background-color: var(--color-primary);
            }
        `;
    }

    static get properties() {
        return {
            path: {
                attribute: true,
                type: String,
            },
            meta: {
                type: Object
            }
        };
    }

    getPreview(meta) {
        if (meta.mime && meta.mime.startsWith('image')) {
            return this.path
        } else {
            return 'assets/preview/unknown.svg';
        }
    }

    hashField([hashName, hashValue]) {
        const baseLoc = window.location.protocol + "//" + window.location.host;
        const fullLink = `${baseLoc}/${hashValue}`
        return html`<div class="group">
                <label>${hashName}</label>
                <input type="text" @click=${() => navigator.clipboard.writeText(fullLink)} value="${hashValue}" readonly>
            </div>`;
    }

    async loadData() {
        return get(`${this.path}?meta`).then(r => {
            return html`
                <div class="sidebar">
                    <div class="group">
                        <label>Name</label>
                        <input type="text" value=${r.name} readonly>
                    </div>
                    <div class="group">
                        <label>Creation Date</label>
                        <input type="text" value=${r.creationDate} readonly>
                    </div>
                    <div class="group">
                        <label>Size</label>
                        <input type="text" value=${r.size} readonly>
                    </div>
                    <hr>
                    ${Object.entries(r.hashes).map(this.hashField)}
                </div>
                <div class="preview">
                    <img src=${this.getPreview(r)}></img>
                </div>
            `;
        })
    }

    render() {
        return html`
            ${until(this.loadData(), html`loading...`)}
        `;
    }
}
customElements.define('ik-detail', Detail);
