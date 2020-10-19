import { all, call } from "redux-saga/effects";
import { watchWalletStatusSaga } from "./manageWallet";

export const getRootSaga = () => {
    return [watchWalletStatusSaga]
}

// single entry point to start all Sagas at once
export default function* rootSagas() {
    yield all([call(watchWalletStatusSaga)]);
}
