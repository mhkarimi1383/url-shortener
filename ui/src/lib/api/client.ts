import router from '@/router';
import axios, { type AxiosError } from 'axios';

const apiBaseURL = '/api';
let userToken: undefined | string = undefined;

const client = axios.create({
  baseURL: apiBaseURL,
  headers: {
    'Content-Type': 'application/json',
  },
});

const loginInterceptor = function (error: AxiosError) {
  if (error.status === 401) {
    router.push('/user/login');
  }
  return Promise.reject(error);
};

client.interceptors.response.use(undefined, loginInterceptor);

export function setToken(token: string | null) {
  if (token === null) {
    client.defaults.headers.common.Authorization = undefined;
    userToken = undefined;
  } else {
    client.defaults.headers.common.Authorization = `Bearer ${token}`;
    userToken = token;
  }
}

export function getToken() {
  return userToken;
}

export default client;
