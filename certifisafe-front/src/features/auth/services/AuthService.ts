import { LOGIN_URL, OAUTH_URL, REGISTER_URL } from "api";
import axios from "axios";
import { useLocation } from "react-router-dom";
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