import Any = jasmine.Any;
import {Deserializer} from "ts-jsonapi";
import {withjwt} from "@/_helpers/withjwt";
import {jsoner} from "@/_helpers/jsoner";

export function get_goods(jwt: string, description: String | null, page: number, limit: number): Promise<Any> {
    let descr = (description !== null || description == "") ? `&filter[description]=${description}` : "";
    return jsoner(withjwt(`http://api.localhost/inventory/good?page[number]=${page}&page[limit]=${limit}&page[order]=desc` + descr, jwt)).then((json) => {
        return new Deserializer({keyForAttribute: "snake_case"}).deserialize(json)
    })
}

export function get_goods_count(jwt: string, description: String): Promise<Number> {
    return jsoner(withjwt(`http://api.localhost/inventory/good/count` + ((description !== null || description == "") ? `?filter[description]=${description}` : ""), jwt)).then((json) => {
        return new Deserializer({keyForAttribute: "snake_case"}).deserialize(json)
    }).then((data) => {
        return data.value
    })
}