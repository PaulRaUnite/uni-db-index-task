export function login(user_id: number): Promise<string> {
    return fetch(`http://api.localhost/user/${user_id}/login`)
        .then(handleResponse)
}

function handleResponse(response: Response): Promise<string> {
    return response.json().then(data => {
        if (!response.ok) {
            const error = (data && data.message) || response.statusText;
            return Promise.reject(error);
        }

        return data.data.id;
    });
}