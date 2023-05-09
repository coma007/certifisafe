const API_URL = "http://localhost:8080/api/";

export const CERTIFICATES_URL = () => API_URL + "certificate";
export const CERTIFICATES_ID_URL = (id: number) => API_URL + "certificate/" + id;
export const CERTIFICATES_DOWNLOAD_URL = (id: number) => API_URL + "certificate/" + id + "/download";
export const CERTIFICATES_WITHDRAW_URL = (id: number) => API_URL + "certificate/" + id + "/withdraw";
export const CERTIFICATES_IS_VALID_ID = (id: number) => API_URL + "certificate/" + id + "/valid";
export const CERTIFICATES_IS_VALID_FILE = (id: number) => API_URL + "certificate/valid";

export const REQUESTS_URL = API_URL + "request";