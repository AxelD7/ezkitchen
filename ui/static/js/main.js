function csrfFetch(url, options = {}) {
    const token = document
        .querySelector('meta[name="csrf-token"]')
        ?.getAttribute("content")

    return fetch(url, {
        credentials: "same-origin",
        headers: {
            "Content-Type": "application/json",
            "X-CSRF-Token": token,
            ...options.headers,
        },
        ...options,
    })
}
