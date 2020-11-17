import React, {useState} from 'react';
import { MainFrame } from './components/MainFrame';
import './App.css';
import {useDispatch, useSelector} from "react-redux";
import {RootStore} from "./Store";
import { GetWalletAddress } from '../src/api/Wallet'

function App() {

  return (
    <div>
        <MainFrame currentPage="home" />
    </div>
    

  );
}

export default App;