import {withjwt} from "@/_helpers/withjwt";
import {jsoner} from "@/_helpers/jsoner";
import {Deserializer} from "ts-jsonapi";

export function get_complaints(username: string, jwt: string): Promise<any> {
    return jsoner(withjwt(`http://api.localhost/user/${username}/complaints`, jwt)).then((json) => {
        return new Deserializer({keyForAttribute: "snake_case"}).deserialize(json)
    })
}

export function new_complaint(username: string, jwt: string, text: string): Promise<any> {
    console.log(text);
    let headers = new Headers();
    headers.set('Authorization', 'Bearer ' + jwt);
    return jsoner(fetch("http://api.localhost/complaint",
        {
            method: 'POST',
            headers: headers,
            body: JSON.stringify({
                data: {
                    type: "complaints",
                    attributes: {
                        body: text
                    },
                }
            }),
        }))
        .then((json) => {
            return new Deserializer({keyForAttribute: "snake_case"}).deserialize(json)
        });
}