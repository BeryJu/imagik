const baseUrl = new URL('/api/priv/', window.location);

export const getAuthorization = () => {
    const auth = window.localStorage.getItem('authorization');
    return auth || undefined;
};

export const setAuthorization = (user, pass) => {
    const auth = btoa(`${user}:${pass}`);
    window.localStorage.setItem('authorization', 'Basic '+auth);
};

export const logout = () => {
    window.localStorage.removeItem('authorization');
    window.location.reload();
};

export const get = async (url) => {
    return request(url);
};
export const post = async (url, body) => {
    return request(url, body);
};

export const request = async (url, body, options = {}) => {
    const headers = {
        'accepts': 'application/json',
        'authorization': getAuthorization(),
        ...options.headers,
    };

    if (!options.method && body) {
        options.method = 'POST';
    }

    if (options.method && options.method !== 'GET') {
        headers['content-type'] = 'application/json';

        if (body !== undefined && typeof body !== 'string') {
            options.body = JSON.stringify(body);
        }
    }

    return fetch(new URL(url, baseUrl), {
        ...options,
        headers,
        body,
    }).then(async (res) => {
        if (!res.ok) {
            if (res.status === 401) logout();

            return res.json().then(({error}) => {
                console.error(e);
                console.error('api error: ' + error);
                throw new Error(error);
            });
        }
        if (res.status == 201) {
            return res;
        }
        return res.json();
    }, (e) => {
        console.error(e);
        console.error('network unreachable: ' + e.message);
        throw e;
    });
};
