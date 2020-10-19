import { ActionType } from "typesafe-actions";
import { ManageWalletActions } from "./manageWallet";

export const RootActions = {
  ...ManageWalletActions
};

export type RootActions = ActionType<typeof RootActions>;
