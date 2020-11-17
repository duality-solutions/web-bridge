import { createStore,  combineReducers, applyMiddleware } from 'redux';
import thunk,  { ThunkMiddleware } from 'redux-thunk';
import { walletReducer } from "../reducers/walletReducers";
import { AppActions } from '../types/actions';

export const rootReducer = combineReducers({
    wallets: walletReducer
})

export type AppState = ReturnType<typeof rootReducer>;

export const store = createStore(
    rootReducer,
    applyMiddleware(thunk)
)