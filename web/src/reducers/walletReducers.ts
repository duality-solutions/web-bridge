import { 
    AddressRes, 
    MnemonicRes, 
    ImportMnemonicReq, 
    ImportMnemonicRes,  
    UnlockWalletReq,
    UnlockWalletRes,
    EncryptWalletReq,
    EncryptWalletRes
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
 } from '../types/actions';

interface DefaultStateI {
    loading: boolean,
    walletAddressResponse?: AddressRes,
    mnemonicResponse?: MnemonicRes[],
    importMnemonicRequest?: ImportMnemonicReq[],
    importMnemonicResponse?: ImportMnemonicRes,
    unlockWalletRequest?: UnlockWalletReq[],
    unlockWalletResponse?: UnlockWalletRes,
    encryptWalletRequest?: EncryptWalletReq,
    encryptWalletResponse?: EncryptWalletRes
}

const defaultState: DefaultStateI = {
    loading: false
}

const walletReducer = (state: DefaultStateI = defaultState, action: WalletActionTypes): DefaultStateI => {
    switch(action.type) {
        case WALLET_ADDRESS_RESPONSE: 
            return {
                loading: false,
                walletAddressResponse: action.address_response
            }
        case MNEMONIC_RESPONSE:
            return {
                loading: false,
                mnemonicResponse: action.mnemonic_response
            }
        case IMPORT_MNEMONIC_REQUEST:
            return {
                loading: true,
                importMnemonicRequest: action.mnemonic_request
            }
        case IMPORT_MNEMONIC_RESPONSE:
            return {
                loading: false,
                importMnemonicResponse: action.restore_mnemonic_response
            }
        case UNLOCK_WALLET_REQUEST:
            return {
                loading: true,
                unlockWalletRequest: action.wallet_request
            }
        case UNLOCK_WALLET_RESPONSE:
            return {
                loading: false,
                unlockWalletResponse: action.wallet_response
            }
        case ENCRYPT_WALLET_REQUEST:
            return {
                loading: true,
                encryptWalletRequest: action.encrypt_request
            }
        case ENCRYPT_WALLET_RESPONSE:
            return {
                loading: false,
                encryptWalletResponse: action.encrypt_response
            }
        default: 
            return state
    }
}

export { walletReducer };