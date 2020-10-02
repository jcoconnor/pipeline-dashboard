import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import * as serviceWorker from './serviceWorker';

import { Provider } from 'mobx-react'
import DataStore from './stores/dataStore'

// import 'bootstrap/dist/css/bootstrap.min.css';
import '@puppet/react-components/source/scss/library/ui.scss';

class RootStore {
  constructor() {
    this.dataStore           = new DataStore(this)
    this.dataStore.fetchProducts().then((store) => {
      // set state here to re-render when new data appeared
    });
    console.log(this)

  }
}


ReactDOM.render(
	<Provider rootStore={new RootStore()}>
		<App />
	</Provider>,
	document.getElementById('root')
)
// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
