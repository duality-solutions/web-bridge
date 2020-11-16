import { 
    AddressRes, 
    MnemonicRes, 
    ImportMnemonicReq, 
    ImportMnemonicRes,  
    UnlockWalletReq,
    UnlockWalletRes,
    EncryptWalletReq,
    EncryptWalletRes
} from './Wallet';

export const WALLET_ADDRESS_RESPONSE = "WALLET_ADDRESS_RESPONSE";
export const MNEMONIC_RESPONSE = "MNEMONIC_RESPONSE";
export const IMPORT_MNEMONIC_REQUEST = "IMPORT_MNEMONIC_REQUEST";
export const IMPORT_MNEMONIC_RESPONSE = "IMPORT_MNEMONIC_RESPONSE";
export const UNLOCK_WALLET_REQUEST = "UNLOCK_WALLET_REQUEST";
export const UNLOCK_WALLET_RESPONSE = "UNLOCK_WALLET_RESPONSE";
export const ENCRYPT_WALLET_REQUEST = "ENCRYPT_WALLET_REQUEST";
export const ENCRYPT_WALLET_RESPONSE = "ENCRYPT_WALLET_RESPONSE";

export interface GetWalletAddress {
    type: typeof WALLET_ADDRESS_RESPONSE;
    payload: AddressRes; 
}

export interface GetMnemonic {
    type: typeof MNEMONIC_RESPONSE;
    mnemonic_response: MnemonicRes[]
}

export interface GetWalletAddresses {
    type: typeof WALLET_ADDRESS_RESPONSE;
    address_response: AddressRes; 
}

export interface MnemonicRequest {
    type: typeof IMPORT_MNEMONIC_REQUEST;
    mnemonic_request: ImportMnemonicReq[]
}

export interface RestoreMnemonic {
    type: typeof IMPORT_MNEMONIC_RESPONSE;
    restore_mnemonic_response: ImportMnemonicRes;
}

export interface UnlockRequest {
    type: typeof UNLOCK_WALLET_REQUEST;
    wallet_request: UnlockWalletReq[];
}

export interface UnlockWallet {
    type: typeof UNLOCK_WALLET_RESPONSE;
    wallet_response: UnlockWalletRes
}

export interface EncryptRequest {
    type: typeof ENCRYPT_WALLET_REQUEST;
    encrypt_request: EncryptWalletReq;
}

export interface EncryptWallet {
    type: typeof ENCRYPT_WALLET_RESPONSE
    encrypt_response: EncryptWalletRes
}

export type WalletActionTypes = 
    | GetMnemonic
    | GetWalletAddresses
    | MnemonicRequest
    | RestoreMnemonic
    | UnlockRequest
    | UnlockWallet
    | EncryptRequest
    | EncryptWallet;


export type AppActions = WalletActionTypes;