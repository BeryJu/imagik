class PathHandler extends EventTarget {
    constructor() {
        super();

        this.url = new URL('/', window.location);
    }

    get path() {
        return this.url.pathname;
    }
    set path(v) {
        this.url.pathname = v;
        this.emit();
    }

    on(name, handler) {
        this.addEventListener(name, handler);
    }

    emit() {
        this.dispatchEvent(new CustomEvent('change', {detail: this.path}));
    }

    get() {
        return this.path;
    }
    set(v) {
        this.path = v;
    }

    getAbsolute(v) {
        return new URL(v.startsWith('./') ? v : './' + v, this.url).pathname;
    }

    navigate(v) {
        this.url = new URL(v, this.url);
    }
    up() {
        this.navigate('..');
    }
    into(v) {
        this.navigate(v.startsWith('./') ? v : './' + v);
    }
}

export const path = new PathHandler('/');
export default path;

export const relate = (a, b) => {
    const url = new URL(a, window.location);
    return new URL('./' + b, url).pathname;
};
