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

export interface UnlockWalletRequest {
    passphrase: string;
    timeout: number;
};

export interface UnlockWalletResponse {
    result: string;
};

export interface EncryptWalletRequest {
    passphrase: string;
};

export interface EncryptWalletResponse {
    result: string;
};

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
    const wordCount = mnemonic.mnemonic.split(" ").length;
    if (
      wordCount === 12 ||
      wordCount === 13 ||
      wordCount === 24 ||
      wordCount === 25
    ) {
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
    } else {
        var response: ImportMnemonicResponse = {
            done: "Word count error"
        };
        return response;
    }
}

export const UnlockWallet = async (request: UnlockWalletRequest): Promise<UnlockWalletResponse> => {
    return await axios.patch<UnlockWalletResponse>("/wallet/unlock", request, RequestConfig).then(function (response) {
        return response.data;
    }).catch(function (error) {
        var errMessage = "UnlockWallet execute [Patch] /wallet/unlock error: " + error;
        console.log(errMessage);
        var errResponse: UnlockWalletResponse = {
            result: "failed"
        }
        return errResponse;
    });
}

export const EncryptWallet = async (request: EncryptWalletRequest): Promise<EncryptWalletResponse> => {
    return await axios.patch<EncryptWalletResponse>("/wallet/encrypt", request, RequestConfig).then(function (response) {
        return response.data;
    }).catch(function (error) {
        var errMessage = "EncryptWallet execute [Patch] /wallet/encrypt error: " + error;
        console.log(errMessage);
        var errResponse: EncryptWalletResponse = {
            result: "failed"
        }
        return errResponse;
    });
}