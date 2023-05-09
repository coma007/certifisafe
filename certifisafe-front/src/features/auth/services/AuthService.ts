import { LOGIN_URL } from "api";
import { Credentials } from "../models/AuthModels";
import axios from "axios";

export const AuthService = {

    login: async (credentials: Credentials): Promise<string> => {
        let url = LOGIN_URL();

        let response = await axios.post(url, credentials);

        // const response = await axios.post(url, credentials);
        return response.data;
    },
}