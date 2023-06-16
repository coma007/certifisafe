import { GET_USER_URL, LOGIN_URL, PASSWORD_RESET_REQUEST_URL, PASSWORD_RESET_URL, REGISTER_URL, TWO_FACTOR_AUTH_URL } from "api";
import axios from "axios";
import { Credentials, User, UserRegister } from "../types/User";
import { Code, PasswordReset, PasswordResetRequest } from "../types/Verification";

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

  requestPasswordReset: async (request: PasswordResetRequest): Promise<string> => {
    let url = PASSWORD_RESET_REQUEST_URL();
    let response = await axios.post(url, request);
    return response.data;
  },

  passwordReset: async (request: PasswordReset): Promise<string> => {
    let url = PASSWORD_RESET_URL();
    let response = await axios.post(url, request);
    return response.data;
  },

  logout: () => {
    localStorage.removeItem("token")
  },

  tfactorauth: async (code: Code): Promise<string> => {
    let url = TWO_FACTOR_AUTH_URL();
    let response = await axios.post(url, code);
    return response.data;
  },

  getUserData: async (): Promise<User> => {
    let url = GET_USER_URL();
    let response = await axios.get(url);
    return response.data;
  },
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