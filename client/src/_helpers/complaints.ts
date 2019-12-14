import {withjwt} from "@/_helpers/withjwt";
import {jsoner} from "@/_helpers/jsoner";
import {Deserializer} from "ts-jsonapi";

export function get_complaints(username: string, jwt: string): Promise<any> {
    return jsoner(withjwt(`http://api.localhost/user/${username}/complaints`, jwt)).then((json) => {
        return new Deserializer({keyForAttribute: "snake_case"}).deserialize(json)
    })
}