import { take } from "redux-saga/effects";
import { getType } from "typesafe-actions";
import { ManageWalletActions } from "../actions/manageWallet";

function* waitForMnemonicImport() {
  console.log("waitForReindex -- /WalletImportMnemonic");
  for (;;) {
    const { reindexComplete } = yield {
      reindexComplete: take(getType(ManageWalletActions.walletComplete))
    };
    if (reindexComplete) {
      console.log("waitForReindex complete");
      return;
    }
  }
}

export function* watchWalletStatusSaga() {
  yield waitForMnemonicImport();
}
