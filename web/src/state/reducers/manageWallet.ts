import { ManageWalletActions } from "../actions/manageWallet";
import { getType } from "typesafe-actions";

export interface ManageWalletState {
  complete: boolean;
}

const defaultState: ManageWalletState = {
  complete: false
};

export const manageWallet = (
  state: ManageWalletState = defaultState,
  action: ManageWalletActions
): ManageWalletState => {
  switch (action.type) {
    case getType(ManageWalletActions.walletInit):
      return { ...state, complete: true };
    case getType(ManageWalletActions.walletImportMnemonic):
      return { ...state, complete: true };
    case getType(ManageWalletActions.walletEncrypt):
      return { ...state, complete: true };
    case getType(ManageWalletActions.walletAddUser):
      return { ...state, complete: true };
    case getType(ManageWalletActions.walletAddLink):
      return { ...state, complete: true };
    case getType(ManageWalletActions.walletComplete):
      return { ...state, complete: true };
    default:
      return state;
  }
};
