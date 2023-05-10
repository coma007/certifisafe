import { LOGIN_URL, REGISTER_URL } from "api";
import { Credentials, UserRegister } from "../models/AuthModels";
import axios from "axios";

export const AuthService = {

    login: async (credentials: Credentials): Promise<string> => {
        let url = LOGIN_URL();

        let response = await axios.post(url, credentials);

        // const response = await axios.post(url, credentials);
        return response.data;
    },

    register: async (user: UserRegister): Promise<string> => {
        let url = REGISTER_URL();

        let response = await axios.post(url, user);

        // const response = await axios.post(url, credentials);
        return response.data;
    },

    
}

axios.interceptors.request.use(
    config => {
      const token = localStorage.getItem("token")
      if (token) {
        config.headers['Authorization'] = token + ' Bearer'
      }
      // config.headers['Content-Type'] = 'application/json';
      return config
    },
    error => {
      Promise.reject(error)
    }
  )