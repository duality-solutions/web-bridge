import axios from 'axios';
import { Dispatch } from 'redux';
import { RequestConfig }from "./Config";
import { 
    AddressRes,
    MnemonicRes,
    ImportMnemonicReq,
    ImportMnemonicRes,
    UnlockWalletReq,
    UnlockWalletRes,
    EncryptWalletReq,
    EncryptWalletRes,
} from '../types/Wallet';
import { 
    WALLET_ADDRESS_RESPONSE,
    MNEMONIC_RESPONSE,
    IMPORT_MNEMONIC_REQUEST,
    IMPORT_MNEMONIC_RESPONSE,
    UNLOCK_WALLET_REQUEST,
    UNLOCK_WALLET_RESPONSE,
    ENCRYPT_WALLET_REQUEST,
    ENCRYPT_WALLET_RESPONSE,
    WalletActionTypes,
    AppActions
 } from '../types/actions';
 import { AppState } from '../store/configureStore';

export const getMnemonicResponse = (mnemonic_response: MnemonicRes[]): AppActions => ({
    type: MNEMONIC_RESPONSE,
    mnemonic_response
})

export const getWalletAdress = (address_response: AddressRes ): AppActions => ({
    type: WALLET_ADDRESS_RESPONSE,
    address_response
})

export const MnemonicRequest = (mnemonic_request: ImportMnemonicReq[]): AppActions => ({
    type: IMPORT_MNEMONIC_REQUEST,
    mnemonic_request
}) 

export const RestoreMnemonicResponse = (restore_mnemonic_response: ImportMnemonicRes ): AppActions => ({
    type: IMPORT_MNEMONIC_RESPONSE,
    restore_mnemonic_response
})

export const UnlockRequest = (wallet_request: UnlockWalletReq[]): AppActions => ({
    type: UNLOCK_WALLET_REQUEST,
    wallet_request
})

export const UnlockResponse = (wallet_response: UnlockWalletRes): AppActions => ({
    type: UNLOCK_WALLET_RESPONSE,
    wallet_response
})

export const EncryptRequest = (encrypt_request: EncryptWalletReq): AppActions => ({
    type: ENCRYPT_WALLET_REQUEST,
    encrypt_request
})

export const EncryptResponse = (encrypt_response: EncryptWalletRes): AppActions => ({
    type: ENCRYPT_WALLET_RESPONSE,
    encrypt_response
})

export const GetWalletAddress = () => async (dispatch: Dispatch<WalletActionTypes>) => {
    try {
        const response = await axios.get("/wallet/defaultaddress", RequestConfig);
        console.log(response, "response inside Wallet Address Action. ts")
        dispatch({
            type: WALLET_ADDRESS_RESPONSE,
            address_response: response.data
        })
    } catch(error) {
        var errMessage = "GetWalletAddresses execute [Get] /wallet/defaultaddress error: " + error;
        var errResponse: WalletAddressResponse = {
            address: errMessage
        }
        return errResponse;
    }

}

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