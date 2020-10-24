import {LitElement, html, css} from 'lit-element';

class GpFile extends LitElement {
    static get styles() {
        return css`
            :host {
                display: flex;
                background-color: var(--color-primary-background-light);
                padding: 1rem;
                border-radius: 5px;
                gap: 1rem;
            }
            img {
                width: 5rem;
                height: 5rem;
                background-color: #111;
                border-radius: 1rem;
            }
        `;
    }

    static get properties() {
        return {
            path: {
                attribute: true,
                type: String,
            },
            mime: {
                attribute: true,
                type: String,
            },
        };
    }

    constructor() {
        super();
        this.addEventListener('click',
            ()=>this.isFolder() ? this.dispatchEvent(new CustomEvent('navigate', {
                detail: this.path,
                composed: true,
                bubbles: true,
            })):'',
        );
    }

    isFolder() {
        return !this.mime;
    }

    getPreview() {
        if (this.mime && this.mime.startsWith('image')) return this.path;
        return 'assets/preview/unknown.svg';
    };

    render() {
        return html`
            <img src=${this.getPreview()}></img>
            <slot></slot>
        `;
    }
}
customElements.define('gp-file', GpFile);
