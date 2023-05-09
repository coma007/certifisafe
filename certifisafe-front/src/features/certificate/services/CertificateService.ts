import axios from 'axios'
import { CERTIFICATES_DOWNLOAD_URL, CERTIFICATES_ID_URL, CERTIFICATES_URL, CERTIFICATES_WITHDRAW_URL } from 'api/index';
import { Certificate } from '../types/Certificate';

export const CertificateService = {

    getAll: async (): Promise<Certificate[]> => {
        let url = CERTIFICATES_URL();
        const response = await axios.get(url);
        return response.data;
    },

    getById: async (id: number): Promise<Certificate> => {
        let url = CERTIFICATES_ID_URL(id);
        const response = await axios.get(url);
        return response.data;
    },

    download: async (id: number) => {
        let url = CERTIFICATES_DOWNLOAD_URL(id);
        const response = await axios.get(url, { responseType: 'blob' });

        const blob = new Blob([response.data], { type: response.headers['content-type'] });
        const urlObject = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = urlObject;
        link.download = `certificate_${id}.crt`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        window.URL.revokeObjectURL(urlObject);
    },

    withdraw: async (id: number): Promise<Certificate> => {
        let url = CERTIFICATES_WITHDRAW_URL(id);
        const response = await axios.patch(url);
        return response.data;
    }
}