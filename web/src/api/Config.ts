import axios,  { AxiosRequestConfig } from 'axios';
import { IceServerConfig } from './IceServerConfig';
import { WebServerConfig } from './WebServerConfig';

// TODO: read settings file to get rest web server URL and port instead of using RestBaseUrl constant variable 
export const RestBaseUrl: string = "http://localhost:35350/api/v1";

export interface ConfigurationResponse {
    IceServers: IceServerConfig[];
    WebServer: WebServerConfig;
};

export const RequestConfig: AxiosRequestConfig = {
    headers: {
        'Access-Control-Allow-Origin': '*',
        'Cache-Control': 'no-cache',
        'Content-Type': 'text/plain',
        'Accept': 'application/json'
    },
    responseType: 'json'
};

export const GetConfigSettings = async (): Promise<ConfigurationResponse> => {
    return await axios.get<ConfigurationResponse>("/config", RequestConfig).then(function (response) {
        return response.data;
    }).catch(function (error) {
        console.log("GetConfigSettings execute [Get] /config error: " + error);
        var nullConfig: ConfigurationResponse = {
        IceServers: [],
        WebServer: {
                BindAddress: "",
                AllowCIDR: "",
                ListenPort: 0
            }
        };
        return nullConfig;
    });
}
