import {jsoner} from "@/_helpers/jsoner";
import {withjwt} from "@/_helpers/withjwt";
import {Deserializer} from "ts-jsonapi";

export function get_countries(jwt: string): Promise<any> {
    return jsoner(withjwt(`http://api.localhost/country/`, jwt)).then((json) => {
        return new Deserializer({keyForAttribute: "snake_case"}).deserialize(json)
    })
}
