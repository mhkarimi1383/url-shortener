import axios from 'axios';
import type { AxiosError, AxiosResponse } from 'axios';
import router from '@/router';
import { message } from 'ant-design-vue';

const apiBaseURL = '/api';
export const loginStateCookie = "loginState";
export const loginInfoCookie = "loginInfo";

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
    "Content-Type": "application/json",
  },
})

const loginInterceptor = function (response: AxiosError) {
  if (response.status === 401) {
    message.error((response.response?.data as errorResponse).message)
    router.push("/user/login");
  }
}

export function setToken(token: string | null) {
  if (token === null) {
    client.defaults.headers.common.Authorization = undefined;
  } else {
    client.defaults.headers.common.Authorization = `Bearer ${token}`;
  }
}

client.interceptors.response.use(undefined, loginInterceptor)

export async function login(info: loginInfo): Promise<loginResponse | errorResponse> {
  let retVal = <loginResponse | errorResponse>{};
  await client.post<loginResponse>("/user/login/", info).then(
    (resp: AxiosResponse) => {
      retVal = resp.data;
    }
  ).catch(
    (err: AxiosError) => {
      retVal = (err.response?.data as errorResponse) || <errorResponse>{
        message: "Unknown error",
      };
    }
  );
  return retVal;
}
