import { combineReducers } from 'redux';
import { walletReducer } from './wallet/reducers'

const rootReducer = combineReducers({
  wallet: walletReducer
})

export type RootState = ReturnType<typeof rootReducer>