import {jsoner} from "@/_helpers/jsoner";

export function login(username: string, password: string): Promise<string> {
    let headers = new Headers();
    headers.set('Authorization', 'Basic ' + btoa(username + ":" + password));
    return jsoner(fetch(`http://api.localhost/user/login`,
        {
            method: 'GET',
            headers: headers,
        })).then(data => {
            console.log(data);
        return data.data.id;
    })
}
