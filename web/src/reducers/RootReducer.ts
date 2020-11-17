import { combineReducers } from 'redux';
import { walletReducer } from './walletReducers';

const RootReducer = combineReducers({
    wallet: walletReducer
});

export default RootReducer;