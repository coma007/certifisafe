import { REQUESTS_BY_USER_SIGNING_URL, REQUESTS_BY_USER_URL, REQUESTS_URL } from "api";
import { CertificateRequest } from "../types/CertificateRequest";
import axios from "axios";

export const CertificateRequestService = {

    getAllByUserSigning: async (): Promise<CertificateRequest[]> => {
        let url = REQUESTS_BY_USER_SIGNING_URL;
        const response = await axios.get(url);
        return response.data;
    },

    getByUser: async (): Promise<CertificateRequest[]> => {
        let url = REQUESTS_BY_USER_URL;
        const response = await axios.get(url);
        return response.data;
    }
}