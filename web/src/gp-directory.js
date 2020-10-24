import {LitElement, html, css} from 'lit-element';

class GpDirectory extends LitElement {
    static get styles() {
        return css`
            :host {
                display: block;
                background-color: var(--color-primary-background-light);
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

    render() {
        return html`
            <slot></slot>
        `;
    }
}
customElements.define('gp-directory', GpDirectory);
