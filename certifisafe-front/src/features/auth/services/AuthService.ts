
import axios from "axios";
import { OAUTH_URL, LOGIN_URL, REGISTER_URL, TWO_FACTOR_AUTH_URL  } from "api";
import { useLocation } from "react-router-dom";
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

  oauth: async () => {
    // fetch('http://localhost:8080/api/oauth')
    //   .then(response => {
    //     if (response.ok) {
    //       // Redirect the user to the authentication URL received from the backend
    //       window.location.href = response.url;
    //     } else {
    //       console.error('Failed to initiate OAuth process');
    //     }
    //   })
    //   .catch(error => {
    //     console.error('Error:', error);
    //   });
    // let url = OAUTH_URL();
    // console.log(OAUTH_URL())
    // window.open(OAUTH_URL(), "_self")
    // let response = await axios.get(url);


    window.open(OAUTH_URL(), "_self");
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