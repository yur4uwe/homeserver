const tokenKey = "authTokenHS";

document.body.addEventListener("htmx:configRequest", function (event) {
    cookieKVPs = document.cookie
        .replace(/^\s+/, '')
        .split(';')
        .map((kvp) => {
            const [key, value] = kvp.split('=', 2);
            return { key: key.trim(), value: value ? value.trim() : '' };
        });

    const token = cookieKVPs.find(kvp => kvp.key === tokenKey);
    if (token) {
        event.detail.headers["Authorization"] = "Bearer " + token.value;
    }
});