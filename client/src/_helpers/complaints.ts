import {withjwt} from "@/_helpers/withjwt";
import {jsoner} from "@/_helpers/jsoner";
import {Deserializer} from "ts-jsonapi";

export function get_complaints(username: string, jwt: string): Promise<any> {
    return jsoner(withjwt(`http://api.localhost/user/${username}/complaints`, jwt)).then((json) => {
        return new Deserializer({keyForAttribute: "snake_case"}).deserialize(json)
    }).then((data) => {
        return data.map((v, i, _) => {
            v.status = v.reviewer !== null ? "reviewed" : "in progress";
            v.created_at = v.created_at ? new Date(v.created_at * 1000).toDateString() : "no timestamp";
            v.reviewed_at = v.reviewed_at ? new Date(v.reviewed_at * 1000).toDateString() : "no timestamp";
            return v
        });
    })
}

export function all_complaints(jwt: string): Promise<any> {
    return jsoner(withjwt(`http://api.localhost/complaint`, jwt)).then((json) => {
        return new Deserializer({keyForAttribute: "snake_case"}).deserialize(json)
    }).then((data) => {
        return data.map((v, i, _) => {
            v.status = v.reviewer !== null ? "reviewed" : "in progress";
            v.created_at = v.created_at ? new Date(v.created_at * 1000).toDateString() : "no timestamp";
            v.reviewed_at = v.reviewed_at ? new Date(v.reviewed_at * 1000).toDateString() : "no timestamp";
            return v
        });
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

export function review_complaint(id: string, jwt: string, text: string): Promise<any> {
    let headers = new Headers();
    headers.set('Authorization', 'Bearer ' + jwt);
    return jsoner(fetch(`http://api.localhost/complaint/${id}`,
        {
            method: 'PATCH',
            headers: headers,
            body: JSON.stringify({
                data: {
                    id: id,
                    type: "complaints",
                    attributes: {
                        answer: text
                    },
                }
            }),
        }))
        .then((json) => {
            return new Deserializer({keyForAttribute: "snake_case"}).deserialize(json)
        });
}