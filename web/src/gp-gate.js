import {LitElement, html, css} from 'lit-element';
import './gp-header.js';
import './gp-app.js';
import {getAuthorization, setAuthorization} from './services/api.js';

class GpGate extends LitElement {
    static get styles() {
        return css`
            :host {
                display: block;
                min-height: 100vh;
            }
        `;
    }

    submitLogin(ev) {
        ev.preventDefault();
        const elements = ev.submitter.form.elements;
        setAuthorization(
            elements.namedItem('username').value,
            elements.namedItem('password').value,
        );
        this.requestUpdate();
    }

    render() {
        if (getAuthorization()) {
            return html`
                <gp-app></gp-app>
            `;
        } else {
            return html`
                <gp-header></gp-header>
                <form @submit=${this.submitLogin}>
                    <input type="text" placeholder="username" name="username" />
                    <input type="password" placeholder="password" name="password" />
                    <input type="submit" value="submit" />
                </form>
            `;
        }
    }
}
customElements.define('gp-gate', GpGate);
