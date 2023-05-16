import { LOGIN_URL, REGISTER_URL } from "api";
import axios from "axios";
import { Credentials, UserRegister } from "../types/User";

export const AuthService = {

  login: async (credentials: Credentials): Promise<string> => {
    let url = LOGIN_URL();
    let response = await axios.post(url, credentials);
    return response.data;
  },

  register: async (user: UserRegister): Promise<string> => {
    let url = REGISTER_URL();
    let response = await axios.post(url, user);
    return response.data;
  },

  logout: () => {
    localStorage.removeItem("token")
  }


}

axios.interceptors.request.use(
  config => {
    const token = localStorage.getItem("token")
    if (token) {
      config.headers['Authorization'] = token + ' Bearer'
    }
    return config
  },
  error => {
    Promise.reject(error)
  }
)