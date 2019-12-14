export function jsonresp(response: Response): Promise<any> {
    if (!response.ok) {
        console.log(response.status, response.statusText);
        return Promise.reject({
            status: response.status,
            details: response.statusText,
        });
    }
    return response.json()
}

export function jsoner(promise: Promise<Response>): Promise<any> {
    return promise.then(jsonresp)
}