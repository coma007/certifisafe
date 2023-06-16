const API_URL = "http://localhost:8080/api/";

export const CERTIFICATES_URL = () => API_URL + "certificate";
export const CERTIFICATES_ID_URL = (id: number) => API_URL + "certificate/" + id;
export const CERTIFICATES_DOWNLOAD_URL = (id: number) => API_URL + "certificate/" + id + "/download";
export const CERTIFICATES_WITHDRAW_URL = (id: number) => API_URL + "certificate/" + id + "/withdraw";
export const CERTIFICATES_IS_VALID_ID = (id: number) => API_URL + "certificate/" + id + "/valid";
export const CERTIFICATES_IS_VALID_FILE = () => API_URL + "certificate/valid";

export const REQUESTS_URL = API_URL + "request";
export const REQUESTS_BY_USER_URL = API_URL + "request/user";
export const REQUESTS_BY_USER_SIGNING_URL = API_URL + "request/signing";

export const LOGIN_URL = () => API_URL + "login";
export const TWO_FACTOR_AUTH_URL = () => API_URL + "two-factor-auth";
export const PASSWORD_RESET_REQUEST_URL = () => API_URL + "password-recovery-request";
export const PASSWORD_RESET_URL = () => API_URL + "password-recovery";
export const REGISTER_URL = () => API_URL + "register";