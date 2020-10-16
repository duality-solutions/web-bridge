import axios from 'axios';
import { RequestConfig } from "./Config";

export interface IceServerConfig {
    URL: string;
    UserName: string;
    Credential: string;
};

export interface ConfigurationIceResponse {
    IceServers: IceServerConfig[];
};

export const UpdateIceConfig = async (iceservers: IceServerConfig[]): Promise<ConfigurationIceResponse> => {
    return await axios.post<ConfigurationIceResponse>("/config/ice", iceservers, RequestConfig).then(function (response) {
        return response.data;
    }).catch(function (error) {
        console.log("UpdateIceConfig execute [Post] /config/ice error: " + error);
        var nullIce: ConfigurationIceResponse = {
          IceServers: []
        }
        return nullIce;
    });
}
