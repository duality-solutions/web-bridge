import { AxiosRequestConfig } from 'axios';
import { IceServerConfig } from './IceServerConfig';
import { WebServerConfig } from './WebServerConfig';

export interface ConfigurationResponse {
    IceServers: IceServerConfig[];
    WebServer: WebServerConfig;
};

// TODO: read settings file to get rest web server URL and port instead of using RestBaseUrl constant variable 
export const RestBaseUrl: string = "http://localhost:35350/api/v1";

export const RequestConfig: AxiosRequestConfig = {
    headers: {
        'Access-Control-Allow-Origin': '*',
        'Cache-Control': 'no-cache',
        'Content-Type': 'text/plain',
        'Accept': 'application/json'
    },
    responseType: 'json'
};