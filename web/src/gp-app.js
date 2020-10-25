import {LitElement, html, css} from 'lit-element';
import './gp-header.js';
import './gp-drop.js';
import './gp-list.js';
import {logout} from './services/api.js';

class GpApp extends LitElement {
    static get styles() {
        return css`
            :host {
                display: block;
                min-height: 100vh;
            }
            gp-header a, gp-header a:visited {
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
        this.addEventListener('dragover', (ev)=>{
            ev.preventDefault();
            this.dragover = true;
        }, false);
        this.addEventListener('dragleave', (ev)=>{
            ev.preventDefault();
            this.dragover = false;
        }, false);
        this.addEventListener('drop', (ev)=>{
            ev.preventDefault();
            this.dragover = false;
            this.handleDrop(ev);
        });

        this.path = window.location.hash.slice(1, Infinity) || '/';
    }

    connectedCallback() {
        super.connectedCallback();
        window.addEventListener('hashchange',
            ()=>this.path=window.location.hash.slice(1, Infinity),
        );
    }

    handleDrop(ev) {
        // Prevent default behavior (Prevent file from being opened)
        ev.preventDefault();

        if (ev.dataTransfer.items) {
            // Use DataTransferItemList interface to access the file(s)
            for (const item of ev.dataTransfer.items) {
                // If dropped items aren't files, reject them
                if (item.kind === 'file') {
                    const file = item.getAsFile();
                    this.uploadFile(file);
                } else {
                    console.warn('... ' + item.kind);
                }
            }
        }
    }

    uploadSelect(ev) {
        const input = document.createElement('input');
        input.setAttribute('type', 'file');
        input.setAttribute('accept', 'image/*');
        input.setAttribute('multiple', 'true');
        input.addEventListener('change', (ev)=>{
            const files = ev.target.files;
            for (const file of files) {
                this.uploadFile(file);
            }
        });
        input.click();
    }

    uploadFile(file) {
        console.log(file.name, URL.createObjectURL(file));
    }

    navigate({detail}) {
        this.path = detail;
    }

    render() {
        console.log(this.path);
        if (window.location.hash !== '#'+this.path) window.location.hash='#'+this.path;

        return html`
            <gp-header path=${this.path} @navigate=${(e)=>this.navigate(e)}>
                <a @click=${()=>this.uploadSelect()}>upload</a>
                |
                <a @click=${()=>window.location.reload()}>refresh</a>
                |
                <a @click=${logout}>logout</a>
            </gp-header>

            <gp-list path=${this.path} @navigate=${(e)=>this.navigate(e)}></gp-list>

            <gp-drop ?show=${this.dragover}></gp-drop>
        `;
    }
}
customElements.define('gp-app', GpApp);
