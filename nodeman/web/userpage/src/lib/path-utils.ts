import { User, UserAPIData, UserID } from "./types";
import { config } from "../config/config";

function makeURL(components: string[]): string {
    var current = window.location.origin;
    for( var i = 0; i < components.length; i++ ) {
        if( current.endsWith("/") ) {
            current = new URL(components[i], current).toString();
        } else {
            current = new URL(components[i], current + "/").toString();
        }
    }
    return current;
}

export const PathTools = {
    userpagePath(path: string): string {
        return makeURL([config.USERPAGE_URLPATH, path]);
    },
    apiPath(path: string): string {
        return makeURL([config.API_URLPATH, path]);
    }
}
