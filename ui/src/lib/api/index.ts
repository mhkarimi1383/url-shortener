import axios from 'axios';
import type { AxiosError, AxiosResponse } from 'axios';
import router from '@/router';

const apiBaseURL = '/api';
export const loginStateCookie = 'loginState';
export const loginInfoCookie = 'loginInfo';

export interface errorResponse {
  message: string;
}

export interface loginInfo {
  Username: string;
  Password: string;
}

export interface userInfo {
  Id: number;
  Admin: boolean;
  Username: string;
  Version: number;
  CreatedAt: string;
  UpdatedAt: string;
}

export interface loginResponse {
  Info: userInfo;
  Token: string;
}

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

export function setToken(token: string | null) {
  if (token === null) {
    client.defaults.headers.common.Authorization = undefined;
  } else {
    client.defaults.headers.common.Authorization = `Bearer ${token}`;
  }
}

client.interceptors.response.use(undefined, loginInterceptor);

export async function login(info: loginInfo): Promise<loginResponse | errorResponse> {
  let retVal = <loginResponse | errorResponse>{};
  await client
    .post<loginResponse>('/user/login/', info)
    .then((resp: AxiosResponse) => {
      retVal = resp.data;
    })
    .catch((err: AxiosError) => {
      retVal =
        (err.response?.data as errorResponse) ||
        <errorResponse>{
          message: 'Unknown error',
        };
    });
  return retVal;
}

export async function register(info: loginInfo): Promise<null | errorResponse> {
  let retVal = <null | errorResponse>{};
  await client
    .post<null>('/user/register/', info)
    .then((resp: AxiosResponse) => {
      retVal = resp.data;
    })
    .catch((err: AxiosError) => {
      retVal =
        (err.response?.data as errorResponse) ||
        <errorResponse>{
          message: 'Unknown error',
        };
    });
  return retVal;
}