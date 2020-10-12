
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
