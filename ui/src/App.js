import React, {Component} from 'react';
import {BrowserRouter as Router, Route} from 'react-router-dom';
import NewSubmissionPage from './containers/newSubmissionPage.js';
import UploadsPage from './containers/uploadsPage.js';

import ListOfForms from './containers/listOfForms.js';
import NewFormPage from './containers/newFormPage';

import NewWorkspacePage from './containers/newWorkspacePage';
import ListOfWorkspaces from './containers/listOfWorkspaces.js';


import LoginPage from './containers/loginPage.js';
import ProtectedRoute from './components/protectedRoute.js';

import DraftSubmissionInfoPage from './containers/draftSubmissionInfoPage';
import PublishedSubmissionInfoPage
  from './containers/publishedSubmissionInfoPage';

import 'tachyons';

class App extends Component {
  render() {
    return (
      <Router>
        <section>
          <ProtectedRoute exact path="/" component={ListOfWorkspaces} />
          <ProtectedRoute
            exact
            path="/new_workspace"
            component={NewWorkspacePage}
          />


          <ProtectedRoute
            exact
            path="/workspaces/:workspaceID/"
            component={ListOfForms}
          />
          <ProtectedRoute
            exact
            path="/workspaces/:workspaceID/new_form"
            component={NewFormPage}
          />


          <ProtectedRoute
            path="/workspaces/:workspaceID/new_submission"
            component={NewSubmissionPage}
          />
          <ProtectedRoute
            path="/workspaces/:workspaceID/forms/:formID"
            component={UploadsPage}
          />
          <ProtectedRoute
            path="/workspaces/:workspaceID/submissions/draft/:submissionID"
            component={DraftSubmissionInfoPage}
          />
          <ProtectedRoute
            path="/workspaces/:workspaceID/submissions/published/:submissionID"
            component={PublishedSubmissionInfoPage}
          />


          <Route path="/login" component={LoginPage} />
        </section>
      </Router>
    );
  }
}

export default App;
