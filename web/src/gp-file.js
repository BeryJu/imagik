import {LitElement, html, css} from 'lit-element';

class GpFile extends LitElement {
    static get styles() {
        return css`
            :host {
                display: flex;
                background-color: var(--color-primary-background-light);
                padding: 1rem;
                gap: 1rem;
                flex-direction: column;
                border-radius: 3px;
            }
            img {
                width: 5rem;
                height: 5rem;
                border-radius: 3px;
            }
            span {
                white-space: nowrap;
                overflow: hidden;
                text-overflow: ellipsis;
                max-width: 5rem;
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
            name: {
                attribute: true,
                type: String,
            },
        };
    }

    constructor() {
        super();
        this.addEventListener('click',
            ()=>this.dispatchEvent(new CustomEvent('navigate', {
                detail: this.path,
                composed: true,
                bubbles: true,
            })),
        );
    }

    isFolder() {
        return !this.mime;
    }

    getPreview() {
        if (this.mime && this.mime.startsWith('image')) {
            return this.path
        } else if (this.isFolder()) {
            return "assets/icons/folder-line.svg";
        } else {
            return 'assets/preview/unknown.svg';
        }
    };

    render() {
        return html`
            <img src=${this.getPreview()} lazy></img>
            <span title=${this.name}>${this.name}</span>
        `;
    }
}
customElements.define('gp-file', GpFile);
