import { observable, action, runInAction } from "mobx";
import AuthService from "../utils/auth0.js";

class integrationsStore {
	@observable Integrations = []

	@action getFormIntegrationSettings = async (workspaceID, formID) => {
		console.log(workspaceID)
		console.log(formID)
		let authToken = AuthService.getToken();
		const response = await fetch(`/api/workspaces/${workspaceID}/forms/${formID}/integrations`, {
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
			this.Integrations.replace(data);
		});
	};

	@action updateFormIntegrationSettings = async (workspaceID,formID, result, callback) => {
		let authToken = AuthService.getToken();

		const response = await fetch(`/api/workspaces/${workspaceID}/forms/${formID}/integrations`, {
			method: "POST",
			body: JSON.stringify(result),
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
}

export const IntegrationsStore = new integrationsStore();
