import client from './client';
import { type AxiosError, type AxiosResponse } from 'axios';

export const loginInfoCookie = 'loginInfo';
export const loginStateCookie = 'loginState';
export { setToken, getToken } from './client';

const unknownError = 'Unknown error';
const limitQueryParam = 'Limit';
const offsetQueryParam = 'Offset';

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

export interface entityCreateRequest {
  Name: string;
  Description: string;
}

export interface entity {
  Id: number;
  Name: string;
  Description: string;
  CreatedAt: string;
  UpdatedAt: string;
  Version: number;
  Creator: userInfo;
}

export interface listEntitiesResponse {
  MetaData: metaData;
  Result: entity[];
}

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
        [limitQueryParam]: limit,
        [offsetQueryParam]: offset,
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

export async function createEntity(entity: entityCreateRequest): Promise<null | errorResponse> {
  let retVal = <null | errorResponse>{};
  await client
    .post<null>('/entity/', entity)
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

export async function deleteEntity(Id: number): Promise<null | errorResponse> {
  let retVal = <null | errorResponse>{};
  await client
    .delete<null>('/entity/' + Id.toString() + '/')
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

export async function listEntities(
  limit: number,
  offset: number,
): Promise<listEntitiesResponse | errorResponse> {
  let retVal = <listEntitiesResponse | errorResponse>{};
  await client
    .get<listEntitiesResponse>('/entity/', {
      params: {
        [limitQueryParam]: limit,
        [offsetQueryParam]: offset,
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
