import { ActionType, createAction } from "typesafe-actions";

export enum ManageWalletTypes {
  INIT = "@@wallet/INIT",
  IMPORT_MNEMONIC = "@@wallet/IMPORT_MNEMONIC",
  ENCRYPT = "@@wallet/ENCRYPT",
  ADD_USER = "@@wallet/ADD_USER",
  ADD_LINK = "@@wallet/ADD_LINK",
  COMPLETE = "@@wallet/COMPLETE",
  CANCEL_RESTORE = "@@wallet/CANCEL_RESTORE",
}

export const ManageWalletActions = {
  walletInit: createAction(ManageWalletTypes.INIT)<void>(),
  walletImportMnemonic: createAction(ManageWalletTypes.IMPORT_MNEMONIC)<string>(),
  walletEncrypt: createAction(ManageWalletTypes.ENCRYPT)<void>(),
  walletAddUser: createAction(ManageWalletTypes.ADD_LINK)<void>(),
  walletAddLink: createAction(ManageWalletTypes.ADD_LINK)<void>(),
  walletComplete: createAction(ManageWalletTypes.COMPLETE)<void>(),
  walletRestoreCancel: createAction(ManageWalletTypes.CANCEL_RESTORE)<void>(),
};

export type ManageWalletActions = ActionType<typeof ManageWalletActions>;
