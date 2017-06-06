import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import { Provider } from "mobx-react";
import {MainStore} from "./stores/mainStore.js";
import {PermissionsStore} from "./stores/permissionsStore.js";
import './index.css';
import './assets/animate.css';


ReactDOM.render(
  <Provider MainStore={MainStore} PermissionsStore={PermissionsStore}>
    <App />
  </Provider>,
  document.getElementById('root')
);
