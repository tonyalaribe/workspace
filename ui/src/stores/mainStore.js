import { observable, action, runInAction } from "mobx";
import AuthService from "../utils/auth0.js";

class mainStore {
	@observable Submissions = [];
	@observable AllWorkspaces = [];
	@observable AllForms = [];

	@observable SubmissionInfo = {};
	@observable Changelog = [];
	@observable
	CurrentForm = {
		jsonschema: {
			properties: {}
		},
		uischema: {}
	};

	@action
	updateFormOnServer = async (
		workspaceID,
		formID,
		submissionID,
		formData,
		callback
	) => {
		let authToken = AuthService.getToken();
		const response = await fetch(
			"/api/workspaces/" +
				workspaceID +
				"/forms/" +
				formID +
				"/submissions/" +
				submissionID,
			{
				method: "PUT",
				body: JSON.stringify(formData),
				mode: "cors",
				headers: {
					"Content-type": "application/json",
					authorization: "Bearer " + authToken
				}
			}
		);
		const data = await response.json();
		/* required in strict mode to be allowed to update state: */
		runInAction("update state after fetching data", () => {
			console.log(data);
			callback();
		});
	};

	@action
	submitFormToServer = async (workspaceID, formID, formData, callback) => {
		let authToken = AuthService.getToken();
		const response = await fetch(
			"/api/workspaces/" + workspaceID + "/forms/" + formID + "/new_submission",
			{
				method: "POST",
				body: JSON.stringify(formData),
				mode: "cors",
				headers: {
					"Content-type": "application/json",
					authorization: "Bearer " + authToken
				}
			}
		);

		const data = await response.json();
		/* required in strict mode to be allowed to update state: */
		runInAction("update state after fetching data", () => {
			console.log(data);
			callback();
		});
	};
	@action
	getAllForms = async workspaceID => {
		let authToken = AuthService.getToken();
		const response = await fetch("/api/workspaces/" + workspaceID + "/forms", {
			method: "GET",
			mode: "cors",
			headers: {
				"Content-type": "application/json",
				authorization: "Bearer " + authToken
			}
		});
		const data = await response.json();
		/* required in strict mode to be allowed to update state: */
		runInAction("update state after fetching data", () => {
			console.log(data);
			this.AllForms.replace(data);
		});
	};

	@action
	getFormInfo = async (workspaceID, formID) => {
		this.CurrentForm.id = workspaceID;

		let authToken = AuthService.getToken();

		const response = await fetch(
			"/api/workspaces/" + workspaceID + "/forms/" + formID,
			{
				method: "GET",
				mode: "cors",
				headers: {
					"Content-type": "application/json",
					authorization: "Bearer " + authToken
				}
			}
		);
		const data = await response.json();
		/* required in strict mode to be allowed to update state: */
		runInAction("update state after fetching data", () => {
			console.log(data);
			this.CurrentForm = data;
		});
	};

	@action
	getAllWorkspaces = async () => {
		let authToken = AuthService.getToken();
		const response = await fetch("/api/workspaces", {
			method: "GET",
			mode: "cors",
			headers: {
				"Content-type": "application/json",
				authorization: "Bearer " + authToken
			}
		});
		const data = await response.json();
		/* required in strict mode to be allowed to update state: */
		runInAction("update state after fetching data", () => {
			this.AllWorkspaces.replace(data);
		});
	};

	@action
	submitNewWorkspaceToServer = async (workspace, callback) => {
		let authToken = AuthService.getToken();

		const response = await fetch("/api/new_workspace", {
			method: "POST",
			body: JSON.stringify(workspace),
			mode: "cors",
			headers: {
				"Content-type": "application/json",
				authorization: "Bearer " + authToken
			}
		});
		const data = await response.json();
		/* required in strict mode to be allowed to update state: */
		runInAction("update state after fetching data", () => {
			console.log(data);
			callback();
		});
	};

	@action
	getMySubmissions = async (workspaceID, formID) => {
		let authToken = AuthService.getToken();
		const response = await fetch(
			"/api/workspaces/" + workspaceID + "/forms/" + formID + "/submissions",
			{
				method: "GET",
				mode: "cors",
				headers: {
					"Content-type": "application/json",
					authorization: "Bearer " + authToken
				}
			}
		);
		const data = await response.json();
		/* required in strict mode to be allowed to update state: */
		runInAction("update state after fetching data", () => {
			console.log(data);
			this.Submissions.replace(data);
		});
	};

	@action
	getSubmissionInfo = async (workspaceID, formID, submissionID) => {
		let authToken = AuthService.getToken();

		const response = await fetch(
			"/api/workspaces/" +
				workspaceID +
				"/forms/" +
				formID +
				"/submissions/" +
				submissionID,
			{
				method: "GET",
				mode: "cors",
				headers: {
					"Content-type": "application/json",
					authorization: "Bearer " + authToken
				}
			}
		);
		const data = await response.json();
		console.log(data);
		/* required in strict mode to be allowed to update state: */
		runInAction("update state after fetching data", () => {
			this.SubmissionInfo = data;
		});
	};
	@action
	getSubmissionChangelog = async (workspaceID, formID, submissionID) => {
		let authToken = AuthService.getToken();

		const response = await fetch(
			"/api/workspaces/" +
				workspaceID +
				"/forms/" +
				formID +
				"/submissions/" +
				submissionID + "/changelog",
			{
				method: "GET",
				mode: "cors",
				headers: {
					"Content-type": "application/json",
					authorization: "Bearer " + authToken
				}
			}
		);
		const data = await response.json();
		console.log(data);
		/* required in strict mode to be allowed to update state: */
		runInAction("update state after fetching data", () => {
			this.Changelog = data;
			console.log(data)
		});
	};
}

export const MainStore = new mainStore();
