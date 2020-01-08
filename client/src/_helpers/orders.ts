import {withjwt} from "@/_helpers/withjwt";
import {jsoner} from "@/_helpers/jsoner";
import {Deserializer} from "ts-jsonapi";

export function get_orders(username: string, jwt: string): Promise<any> {
    return jsoner(withjwt(`http://api.localhost/user/${username}/invoice`, jwt)).then((json) => {
        return new Deserializer({keyForAttribute: "snake_case"}).deserialize(json)
    })
}

export function create_order(jwt: string, invoice: { country: string }, parts: [{ quantity: number, good_id: number }]) {
    console.log(invoice);
    console.log(parts);
    let headers = new Headers();
    headers.set('Authorization', 'Bearer ' + jwt);
    return fetch("http://api.localhost/invoice/",
        {
            method: 'POST',
            headers: headers,
            body: JSON.stringify({
                data: {
                    type: "invoices",
                    attributes: {
                        destination_country: invoice.country,

                    },
                    relationships: {
                        invoice_parts: {
                            data: parts.map((v) => {
                                return {type: "invoice_parts", id: v.good_id.toString()}
                            })
                        }
                    }
                },
                included: parts.map((v) => {
                    return {
                        type: "invoice_parts",
                        id: v.good_id.toString(),
                        attributes: v
                    }
                })
            }),
        })
}