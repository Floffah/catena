import createFetchClient from "openapi-fetch";
import createQueryClient from "openapi-react-query";

import { paths } from "../../types/api";

export const apiFetch = createFetchClient<paths>({
    baseUrl: "http://localhost:8080/",
});

export const $api = createQueryClient(apiFetch);
