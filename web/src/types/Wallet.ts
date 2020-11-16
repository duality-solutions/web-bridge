export interface AddressRes {
    address: string;
};

export interface MnemonicRes {
    hdseed: string;
    mnemonic: string;
    mnemonicpassphrase: string;
};

export interface ImportMnemonicReq {
    mnemonic: string;
    language?: string;
    passphrase?: string;
};

export interface ImportMnemonicRes {
    done: string;
}

export interface UnlockWalletReq {
    passphrase: string;
    timeout: number;
};

export interface UnlockWalletRes {
    result: string;
};

export interface EncryptWalletReq {
    passphrase: string;
};

export interface EncryptWalletRes {
    result: string;
};