
export interface WebServerRestartRequest {
    restart_epoch: number
};

export interface RestartResponse {
    result: string
};

export interface WebServerConfig {
    BindAddress: string;
    ListenPort: number;
    AllowCIDR: string;
};

export interface ConfigurationWebResponse {
    result: {
        WebServer: WebServerConfig;
    }
};
