import axios from 'axios';
import { RequestConfig } from "./Config";

interface JsonRpc {
    jsonrpc: string;
    method: string;
    params: (string | number | boolean)[];
    id: string;
};

function isBoolean(value: string): boolean {
    return ((value != null) && (value !== '') && (value.toLocaleLowerCase() === 'false' || value.toLocaleLowerCase() === 'true'));
}

function isNumber(value: string): boolean {
    return ((value != null) && (value !== '') && !isNaN(Number(value.toString())));
}

function stringToParamsArray(s: string[]): (string | number | boolean)[] {
    var objArray: (string | number | boolean)[] = [];
    s.forEach(element => {
        if (isNumber(element)){
            objArray.push(Number(element));
        } else if (isBoolean(element)) {
            if (element.toLocaleLowerCase() === 'true') {
                objArray.push(true);
            } else {
                objArray.push(false);
            }
        } else {
            objArray.push(element);
        }
    });
    return objArray;
}

export const ExecCommand = async (cmd: string) => {
    var parsed: string[] = cmd.split(',');
    if (parsed.length > 0) {
        let method: string = parsed[0];
        let params: string[] = parsed.slice(1, parsed.length);
        let paramsObj: (string | number | boolean)[] = stringToParamsArray(params);
        let command: JsonRpc = {
            jsonrpc: "2.0",
            method: method,
            params: paramsObj,
            id: "123" // TODO: create unique id
        };
        await axios.post<object>("/blockchain/jsonrpc", command, RequestConfig).then(function (response) {
          console.log(JSON.stringify(response.data, null, 2));
        }).catch(function (error) {
          console.log("ExecCommand execute [Post] /blockchain/jsonrpc error: " + error);
        });
    }
}
