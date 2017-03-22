import React, { Component } from 'react';
import {
  BrowserRouter as Router,
  Route,
} from 'react-router-dom'
import UploadsPage from './containers/uploadsPage.js';
import NewSubmissionPage from './containers/newSubmissionPage.js';
import 'tachyons';

class App extends Component {
  render() {
    return (
      <Router>
        <section>
          <Route exact path="/" component={UploadsPage}/>
          <Route path="/new_submission" component={NewSubmissionPage}/>
        </section>
      </Router>
    );
  }
}

export default App;
