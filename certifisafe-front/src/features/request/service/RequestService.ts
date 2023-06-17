import { REQUESTS_BY_USER_SIGNING_URL, REQUESTS_BY_USER_URL, REQUESTS_URL, REQUEST_ACCEPT_URL, REQUEST_DECLINE_URL, REQUEST_DELETE_URL } from "api";
import { CreateRequestDTO, Request } from "../types/Request";
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
    },

    createRequest: async (values: CreateRequestDTO): Promise<Request[]> => {
        let url = REQUESTS_URL;
        const response = await axios.post(url, values);
        return response.data;
    },

    accept: async (id: number) => {
        let url = REQUEST_ACCEPT_URL(id);
        const response = await axios.patch(url);
    },

    decline: async (id: number, reason: string) => {
        let url = REQUEST_DECLINE_URL(id);
        const response = await axios.patch(url, { Reason: reason });
    },

    delete: async (id: number) => {
        let url = REQUEST_DELETE_URL(id);
        const response = await axios.patch(url);
    },

}