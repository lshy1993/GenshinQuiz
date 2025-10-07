
import Axios, { type AxiosRequestConfig, type AxiosRequestHeaders } from 'axios';
import { DateTime } from 'luxon';
import qs from 'qs';


const axios = Axios.create({
    baseURL: '/api/',
    timeout: 5000,
    headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json',
    } as AxiosRequestHeaders,
    withCredentials: false,
    paramsSerializer: (params) => qs.stringify(params, { arrayFormat: 'repeat', skipNulls: true }),
    transformRequest: Axios.defaults.transformRequest ? Array.isArray(Axios.defaults.transformRequest) ? [...Axios.defaults.transformRequest] : [Axios.defaults.transformRequest] : [],
});

export const Fetcher = async<T>(config: AxiosRequestConfig): Promise<T> => {
    // if (getJWT() === '') {
    //     return Promise.reject('No JWT token set');
    // }
    const promise = await axios.request<T>(config)
        .then(res => res.data)
        .catch(err => {
            if (err.response) {
                if (err.response.status === 401) {
                    // 401 Unauthorized
                    // setJWT('');
                }
                return Promise.reject(err.response.data);
            }
            return Promise.reject(err);
        });
    return promise;
}
export default Fetcher;

export const getJWT = () => {
    return axios.defaults.headers.common['Authorization']?.toString() ?? '';
}

export const setJWT = (token: string) => {
    if (token) {
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    } else {
        delete axios.defaults.headers.common['Authorization'];
    }
}