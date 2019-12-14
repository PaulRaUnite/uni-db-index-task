import {jsoner} from "@/_helpers/jsoner";

export function signup(username: string, password: string, name: string): Promise<null> {
    return jsoner(fetch(`http://api.localhost/user/signup`,
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
        }))
}
