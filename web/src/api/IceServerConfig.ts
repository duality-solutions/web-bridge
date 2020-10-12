
export interface IceServerConfig {
    URL: string;
    UserName: string;
    Credential: string;
};

export interface ConfigurationIceResponse {
    IceServers: IceServerConfig[];
};