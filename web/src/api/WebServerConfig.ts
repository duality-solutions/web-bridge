import axios from 'axios';
import { RequestConfig }from "./Config";

export interface WebServerRestartRequest {
    RestartEpoch: number
};

export interface RestartResponse {
    Result: string
};

export interface WebServerConfig {
    BindAddress: string;
    ListenPort: number;
    AllowCIDR: string;
};

export interface ConfigurationWebResponse {
    WebServer: WebServerConfig;
};

export const UpdateWebServerConfig = async (webserver: WebServerConfig): Promise<ConfigurationWebResponse> => {
    return await axios.post<ConfigurationWebResponse>("/config/web", webserver, RequestConfig).then(function (response) {
        return response.data;
    }).catch(function (error) {
        console.log("UpdateWebServerConfig execute [Post] /config/web error: " + error);
        var nullweb: ConfigurationWebResponse = {
            WebServer: {
                BindAddress: "",
                AllowCIDR: "",
                ListenPort: 0
            }
        }
        return nullweb;
    });
}

export const RestartWebServer = async (postData: WebServerRestartRequest): Promise<RestartResponse> => {
    return await axios.put<RestartResponse>("/config/web/restart", postData, RequestConfig).then(function (response) {
        return response.data;
    }).catch(function (error) {
        var errMessage = "RestartWebServer execute [Put] /config/web/restart error: " + error;
        var nullRestart: RestartResponse = {
            Result: errMessage
        }
        return nullRestart;
    });
}

