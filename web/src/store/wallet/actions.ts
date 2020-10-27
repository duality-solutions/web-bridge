import { WalletStatus, ENCRYPT_WALLET, IMPORT_MNEMONIC, WalletActionTypes } from './types'

export function encryptWallet(newStatus: WalletStatus): WalletActionTypes {
  return {
    type: ENCRYPT_WALLET,
    payload: newStatus
  }
}

export function importMnemonic(newStatus: WalletStatus): WalletActionTypes {
  return {
    type: IMPORT_MNEMONIC,
    payload: newStatus
  }
}