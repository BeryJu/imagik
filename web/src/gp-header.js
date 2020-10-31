import {LitElement, html, css} from 'lit-element';

class GpApp extends LitElement {
    static get styles() {
        return css`
            :host {
                display: flex;
                flex-direction: row;
                align-items: center;
                background-color: var(--color-primary-background-dark);
                color: var(--color-primary-text);
                box-shadow: 0px 2px 3px 0px #0008;
                padding: 0 1rem;
                gap: 1rem;
            }
            #logo {
                height: 2rem;
            }
            .spacer {
                flex-grow: 1;
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

    logoClick() {
        const url = new URL(this.path + "/", window.location);
        const destUrl = new URL("..", url).pathname.slice(0, -1)
        this.dispatchEvent(new CustomEvent('navigate', {
            detail: destUrl, composed: true, bubbles: true,
        }));
    }

    headerPath() {
        if (this.path === "/") {
            return "gopyazo";
        } else {
            return decodeURIComponent(this.path)
        }
    }

    render() {
        return html`
            <img
                id="logo"
                src="assets/logo.svg"
                @click=${()=>this.logoClick()}
            ></img>
            <h1>${this.headerPath()}</h1>
            <div class="spacer"></div>
            <slot></slot>
        `;
    }
}
customElements.define('gp-header', GpApp);
