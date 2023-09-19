import axios from 'axios';
import router from '@/router';
import type { AxiosError, AxiosResponse } from 'axios';

const apiBaseURL = '/api';
export const loginStateCookie = 'loginState';
export const loginInfoCookie = 'loginInfo';
const unknownError = 'Unknown error';

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

interface metaData {
  Count: number;
}

export interface listUsersResponse {
  MetaData: metaData;
  Result: userInfo[];
}

export interface changeUserPasswordRequest {
  Password: string;
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
          message: unknownError,
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
          message: unknownError,
        };
    });
  return retVal;
}

export async function adminChangeUserPassword(
  userId: number,
  info: changeUserPasswordRequest,
): Promise<null | errorResponse> {
  let retVal = <null | errorResponse>{};
  await client
    .put<null>('/user/change-password/' + userId.toString() + '/', info)
    .then((resp: AxiosResponse) => {
      retVal = resp.data;
    })
    .catch((err: AxiosError) => {
      retVal =
        (err.response?.data as errorResponse) ||
        <errorResponse>{
          message: unknownError,
        };
    });
  return retVal;
}

export async function changeUserPassword(
  info: changeUserPasswordRequest,
): Promise<null | errorResponse> {
  let retVal = <null | errorResponse>{};
  await client
    .put<null>('/user/change-password/', info)
    .then((resp: AxiosResponse) => {
      retVal = resp.data;
    })
    .catch((err: AxiosError) => {
      retVal =
        (err.response?.data as errorResponse) ||
        <errorResponse>{
          message: unknownError,
        };
    });
  return retVal;
}

export async function listUsers(
  limit: number,
  offset: number,
): Promise<listUsersResponse | errorResponse> {
  let retVal = <listUsersResponse | errorResponse>{};
  await client
    .get<listUsersResponse>('/user/', {
      params: {
        limit: limit,
        offset: offset,
      },
    })
    .then((resp: AxiosResponse) => {
      retVal = resp.data;
    })
    .catch((err: AxiosError) => {
      retVal =
        (err.response?.data as errorResponse) ||
        <errorResponse>{
          message: unknownError,
        };
    });
  return retVal;
}
