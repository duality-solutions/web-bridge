import { AxiosRequestConfig } from 'axios';

export interface IceServerConfig {
    URL: string;
    UserName: string;
    Credential: string;
};

export interface WebServerConfig {
    BindAddress: string;
    ListenPort: number;
    AllowCIDR: string;
};

export interface ConfigurationResponse {
    result: {
        IceServers: IceServerConfig[];
        WebServer: WebServerConfig;
    }
};

export interface ConfigurationIceResponse {
    result: {
        IceServers: IceServerConfig[];
    }
};

export interface ConfigurationWebResponse {
    result: {
        WebServer: WebServerConfig;
    }
};

export const RestUrl: string = "http://localhost:35350/api/v1/";

export const RequestConfig: AxiosRequestConfig = {
    headers: {
        'Access-Control-Allow-Origin': '*',
        'Cache-Control': 'no-cache',
        'Content-Type': 'text/plain',
        'Accept': 'application/json'
    },
    responseType: 'json'
};