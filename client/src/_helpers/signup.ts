export function signup(username: string, password: string, name: string): Promise<null> {
    return fetch(`http://api.localhost/user/signup`,
        {
            method: 'POST',
            body: JSON.stringify({
                data: {
                    id: username,
                    type: "users",
                    attributes: {
                        name: name,
                        password: password,
                    }
                }
            }),
        })
        .then((response: Response) => {
            if (!response.ok) {
                console.log(response.status, response.statusText);
                return Promise.reject({
                    status: response.status.toString(),
                    details: response.statusText,
                });
            }
            return response.json()
        })
}
