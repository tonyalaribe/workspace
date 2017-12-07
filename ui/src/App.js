import React, { Component } from "react";
import { BrowserRouter as Router, Route } from "react-router-dom";

import NewSubmissionPage from "./containers/submissions/newSubmissionPage.js";
import SubmissionsPage from "./containers/submissions/submissionsPage.js";
import DraftSubmissionInfoPage from "./containers/submissions/draftSubmissionInfoPage";
import PublishedSubmissionInfoPage from "./containers/submissions/publishedSubmissionInfoPage";

import ListOfForms from "./containers/listOfForms.js";
import NewFormPage from "./containers/newFormPage";

import NewWorkspacePage from "./containers/newWorkspacePage";
import ListOfWorkspaces from "./containers/listOfWorkspaces.js";

import LoginPage from "./containers/loginPage.js";
import ProtectedRoute from "./components/protectedRoute.js";
import AuthService from "./utils/auth0.js";
import {inject } from "mobx-react";
import "tachyons";

@inject("MainStore")
class App extends Component {
	componentDidMount(){
		let profile = AuthService.getProfile()
		console.log(profile)
		this.props.MainStore.loadProfile(profile)
		// listen to profile_updated events to update internal state
		AuthService.emitter.on("profile_updated", newProfile => {
			AuthService.setProfile(newProfile)
			this.props.MainStore.loadProfile(newProfile)
		});
	}
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
						exact
						path="/workspaces/:workspaceID/forms/:formID/new_submission"
						component={NewSubmissionPage}
					/>

					<ProtectedRoute
						exact
						path="/workspaces/:workspaceID/forms/:formID"
						component={SubmissionsPage}
					/>

					<ProtectedRoute
						path="/workspaces/:workspaceID/forms/:formID/submissions/draft/:submissionID"
						component={DraftSubmissionInfoPage}
					/>
					<ProtectedRoute
						path="/workspaces/:workspaceID/forms/:formID/submissions/published/:submissionID"
						component={PublishedSubmissionInfoPage}
					/>

					<Route path="/login" component={LoginPage} />
				</section>
			</Router>
		);
	}
}

export default App;
