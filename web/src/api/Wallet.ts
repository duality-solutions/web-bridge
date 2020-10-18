import axios from 'axios';
import { RequestConfig }from "./Config";

export interface WalletAddressResponse {
    address: string;
};

export interface MnemonicResponse {
    hdseed: string;
    mnemonic: string;
    mnemonicpassphrase: string;
};

export interface ImportMnemonicRequest {
    mnemonic: string;
    language?: string;
    passphrase?: string;
};

export interface ImportMnemonicResponse {
    done: string;
}

export const GetMnemonic = async (): Promise<MnemonicResponse> => {
    return await axios.get<MnemonicResponse>("/wallet/mnemonic", RequestConfig).then(function (response) {
        return response.data;
    }).catch(function (error) {
        var errMessage = "GetMnemonics execute [Get] /wallet/mnemonic error: " + error;
        console.log(errMessage);
        var errResponse: MnemonicResponse = {
            hdseed: "",
            mnemonic: "",
            mnemonicpassphrase: ""
        }
        return errResponse;
    });
}

export const GetWalletAddresses = async (): Promise<WalletAddressResponse> => {
    return await axios.get<WalletAddressResponse>("/wallet/defaultaddress", RequestConfig).then(function (response) {
        return response.data;
    }).catch(function (error) {
        var errMessage = "GetWalletAddresses execute [Get] /wallet/defaultaddress error: " + error;
        var errResponse: WalletAddressResponse = {
            address: errMessage
        }
        return errResponse;
    });
}

export const RestoreMnemonic = async (mnemonic: ImportMnemonicRequest): Promise<ImportMnemonicResponse> => {
    return await axios.post<ImportMnemonicResponse>("/wallet/mnemonic", mnemonic, RequestConfig).then(function (response) {
        return response.data;
    }).catch(function (error) {
        var errMessage = "RestoreMnemonic execute [Post] /wallet/mnemonic error: " + error;
        console.log(errMessage);
        var errResponse: ImportMnemonicResponse = {
            done: "failed"
        }
        return errResponse;
    });
}