import { combineReducers } from "redux";
import { manageWallet } from "./manageWallet";

export let rootReducer = combineReducers({
  ...manageWallet
});

export default function createReducer(injectedReducers = {}) {
  rootReducer = combineReducers({
    ...manageWallet,
    ...injectedReducers
  });

  return rootReducer;
}

export type RootState = ReturnType<typeof rootReducer>;
