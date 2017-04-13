import React, { Component } from 'react';
import {
  BrowserRouter as Router,
  Route,
} from 'react-router-dom';
import UploadsPage from './containers/uploadsPage.js';
import NewSubmissionPage from './containers/newSubmissionPage.js';
import LoginPage from './containers/loginPage.js';
import ProtectedRoute from './components/protectedRoute.js';
import DraftSubmissionInfoPage from './containers/draftSubmissionInfoPage';
import PublishedSubmissionInfoPage from './containers/publishedSubmissionInfoPage';

import 'tachyons';

class App extends Component {
  render() {
    return (
      <Router>
        <section>
          <ProtectedRoute exact path="/" component={UploadsPage}/>
          <ProtectedRoute path="/new_submission" component={NewSubmissionPage}/>
          <ProtectedRoute path="/submissions/draft/:submissionID" component={DraftSubmissionInfoPage}/>
          <ProtectedRoute path="/submissions/published/:submissionID" component={PublishedSubmissionInfoPage}/>
          <Route path="/login" component={LoginPage}/>
        </section>
      </Router>
    );
  }
}

export default App;
