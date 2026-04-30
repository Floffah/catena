import createFetchClient from "openapi-fetch";
import createQueryClient from "openapi-react-query";
import {paths} from "../../types/api";

export const apiFetch = createFetchClient<paths>({
    baseUrl: "https://api.catena.build/",
});

export const $api = createQueryClient(apiFetch);