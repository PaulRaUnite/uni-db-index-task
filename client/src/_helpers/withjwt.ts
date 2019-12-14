export function withjwt(query: string, jwt: string): Promise<Response> {
    let headers = new Headers();
    headers.set('Authorization', 'Bearer ' + jwt);
    return fetch(query,
        {
            method: 'GET',
            headers: headers,
        })
}