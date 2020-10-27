import { WalletStatus, WalletActionTypes, ENCRYPT_WALLET, IMPORT_MNEMONIC } from './types';

const initialState: WalletStatus = {
    status: {
        MnemonicBackup: false,
        WalletEncrypted: false,
        HasAccounts: false,
        HasLinks: false,
        HasTransactions: false,
        UnlockedUntil: 0,
    }
}
  
export function walleReducer(
    state = initialState,
    action: WalletActionTypes
): WalletStatus {
    switch (action.type) {
      case ENCRYPT_WALLET:
        return {
            status: state.status
        }
      case IMPORT_MNEMONIC:
        return {
            status: state.status
        }
      default:
        return state
    }
}