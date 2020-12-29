import {LitElement, html, css} from 'lit-element';

class Drop extends LitElement {
    static get styles() {
        return css`
            :host {
                pointer-events: none;
                z-index: 1;
            }
            :host([show]) div {
                opacity: 1;
            }
            div {
                display: flex;
                transition: opacity 0.26s ease-in-out;
                opacity: 0;
                position: fixed;
                top: 0;
                bottom: 0;
                left: 0;
                right: 0;
                backdrop-filter: blur(5px);
                background-color: #0003;
            }
            svg {
                color: white;
                position: absolute;
                top: 50%;
                left: 50%;
                transform: translate(-50%, -50%);
                width: 5rem;
                height: 5rem;
                margin: 0 auto;
                border: 2px solid var(--color-primary);
                border-radius: 5px;
                background-color: #fff5;
            }
        `;
    }

    render() {
        return html`
            <div>
                <svg width="1em" height="1em" viewbox="0 0 100 100">
                    <line
                        x1="50" y1="12" x2="50" y2="88"
                        stroke="currentColor" stroke-width="16px"
                    />
                    <line
                        x1="88" y1="50" x2="12" y2="50"
                        stroke="currentColor" stroke-width="16px"
                    />
                </svg>
            </div>
        `;
    }
}
customElements.define('ik-drop', Drop);
