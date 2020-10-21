import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from "react-redux";
import * as serviceWorker from './serviceWorker';
import App from './App';
import 'semantic-ui-css/semantic.min.css'
import './index.css';
import configureStore from "./state/store";
import { RestBaseUrl } from './api/Config';
import axios from 'axios';

// set default baseURL
axios.defaults.baseURL = RestBaseUrl;

// setup redux store
const initialState = {};
const store = configureStore(initialState);

ReactDOM.render(
  <React.StrictMode>
    <Provider store={store}>
      <App />
    </Provider>
  </React.StrictMode>,
  document.getElementById("root")
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
