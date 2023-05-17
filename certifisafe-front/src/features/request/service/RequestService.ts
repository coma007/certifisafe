import { REQUESTS_BY_USER_SIGNING_URL, REQUESTS_BY_USER_URL } from "api";
import { Request } from "../types/Request";
import axios from "axios";

export const RequestService = {

    getAllByUserSigning: async (): Promise<Request[]> => {
        let url = REQUESTS_BY_USER_SIGNING_URL;
        const response = await axios.get(url);
        return response.data;
    },

    getByUser: async (): Promise<Request[]> => {
        let url = REQUESTS_BY_USER_URL;
        const response = await axios.get(url);
        return response.data;
    }
}