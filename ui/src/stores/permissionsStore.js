import { observable, action, runInAction } from "mobx";
import AuthService from "../utils/auth0.js";

class permissionsStore {
		@observable WorkspaceUsers = []
		@action getWorkspaceUsersAndRoles = async (workspaceID) => {
			let authToken = AuthService.getToken();
			const response = await fetch("/api/users_in_workspace?w="+workspaceID, {
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
				this.WorkspaceUsers.replace(data);
			});
		};

// @action submitNewWorkspaceToServer = async (workspace, callback) => {
// 	let authToken = AuthService.getToken();
//
// 	const response = await fetch("/api/new_workspace", {
// 		method: "POST",
// 		body: JSON.stringify(workspace),
// 		mode: "cors",
// 		headers: {
// 			"Content-type": "application/json",
// 			authorization: "Bearer " + authToken
// 		}
// 	});
// 	const data = await response.json();
// 	/* required in strict mode to be allowed to update state: */
// 	runInAction("update state after fetching data", () => {
// 		console.log(data);
// 		callback();
// 	});
// };

}

export const PermissionsStore = new permissionsStore();
