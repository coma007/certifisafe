import axios from 'axios'
import { CERTIFICATES_DOWNLOAD_URL, CERTIFICATES_ID_URL, CERTIFICATES_URL, CERTIFICATES_WITHDRAW_URL } from 'api/index';

export const getAll = async () => {
    let url = CERTIFICATES_URL();
    const response = await axios.get(url);
    return response.data;
}

export const getById = async (id: number) => {
    let url = CERTIFICATES_ID_URL(id);
    const response = await axios.get(url);
    return response.data;
}

export const download = async (id: number) => {
    let url = CERTIFICATES_DOWNLOAD_URL(id);
    const response = await axios.get(url);
    return response.data;
}

export const withdraw = async (id: number) => {
    let url = CERTIFICATES_WITHDRAW_URL(id);
    const response = await axios.patch(url);
    return response.data;
}