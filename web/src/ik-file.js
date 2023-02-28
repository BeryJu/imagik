import { LitElement, html, css } from "lit";

class File extends LitElement {
    static get styles() {
        return css`
            :host {
                display: flex;
                background-color: var(--color-primary-background-light);
                padding: 1rem;
                gap: 1rem;
                flex-direction: column;
            }
            :host([type=folder]) img {
                min-height: 75%;
            }
            .imageContainer {
                width: 10rem;
                height: 10rem;
                display: flex;
                justify-content: center;
                align-items: center;
            }
            img {
                max-width: 100%;
                max-height: 100%;
            }
            span {
                white-space: nowrap;
                overflow: hidden;
                text-overflow: ellipsis;
                max-width: 10rem;
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
            type: {
                attribute: true,
                type: String,
            },
        };
    }

    constructor() {
        super();
        this.addEventListener("click", () =>
            this.dispatchEvent(
                new CustomEvent("navigate", {
                    detail: this.path,
                    composed: true,
                    bubbles: true,
                }),
            ),
        );
    }

    isFolder() {
        return this.type === "folder";
    }

    getPreview() {
        if (this.mime && this.mime.startsWith("image")) {
            return this.path;
        } else if (this.isFolder()) {
            return "assets/icons/folder-line.svg";
        } else {
            return "assets/preview/unknown.svg";
        }
    }

    render() {
        return html`
            <div class="imageContainer">
                <img src=${this.getPreview()} loading="lazy" />
            </div>
            <span title=${this.name}>${this.name}</span>
        `;
    }
}
customElements.define("ik-file", File);
