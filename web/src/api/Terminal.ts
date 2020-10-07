export interface JsonRpc {
    jsonrpc: string;
    method: string;
    params: (string | number | boolean)[];
    id: string;
};