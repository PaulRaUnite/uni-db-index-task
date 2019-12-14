export function login(username: string, password: string): Promise<string> {
    let headers = new Headers();
    headers.set('Authorization', 'Basic ' + btoa(username + ":" + password));
    return fetch(`http://api.localhost/user/login`,
        {
            method: 'GET',
            headers: headers,
        })
        .then((response: Response) => {
            if (!response.ok) {
                console.log(response.status, response.statusText);
                return Promise.reject({
                    status: response.status.toString(),
                    details: response.statusText,
                });
            }
            return response.json().then(data => {
                return data.data.id;
            })
        })
}
