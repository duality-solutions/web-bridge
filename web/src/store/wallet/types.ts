import { WalletSetupStatus } from "../../api/WalletStatus";
  
export const ENCRYPT_WALLET = 'ENCRYPT_WALLET'
export const IMPORT_MNEMONIC = 'IMPORT_MNEMONIC'

export interface WalletStatus {
    status: WalletSetupStatus
}

interface EncryptWalletAction {
  type: typeof ENCRYPT_WALLET
  payload: WalletStatus
}

interface ImportMnemonicAction {
  type: typeof IMPORT_MNEMONIC
  payload: WalletStatus
}

export type WalletActionTypes = EncryptWalletAction | ImportMnemonicAction