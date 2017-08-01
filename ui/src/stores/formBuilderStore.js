import { observable, action, runInAction,toJS } from "mobx";
import AuthService from "../utils/auth0.js";

class formBuilderStore {
	@observable propertiesOrder = [];
	@observable JSONSchema = {
		type: "object",
		properties:{}
	};
	@observable UISchema = {}
	@observable Kinds = new Map()


	@action
	submitFormToServer = async (workspaceID, formID, callback) => {

		let form = {};
		form.name = this.JSONSchema.title
		form.jsonschema = JSON.parse(toJS(this.JSONSchema));
		form.uischema = JSON.parse(toJS(this.UISchema));


		let authToken = AuthService.getToken();
		const response = await fetch(
			"/api/workspaces/" + workspaceID + "/forms/" + formID + "/new_submission",
			{
				method: "POST",
				body: JSON.stringify(form),
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

	@action onFormKindChange = (propertyKey, kind)=>{
		this.JSONSchema.properties[propertyKey].kind =
			kind;
		this.Kinds.set(propertyKey, kind)
		console.log(this.Kinds)
		switch (kind) {
			case "Short answer":
				this.JSONSchema.properties[propertyKey].type =
					"string";
				break;
			case "Paragraph":
				this.JSONSchema.properties[propertyKey].type =
					"string";
					this.UISchema[propertyKey]={}
				this.UISchema[propertyKey]["ui:widget"] =
					"textarea";
				break;
			case "File upload":
				this.JSONSchema.properties[propertyKey].type = "string"
				this.JSONSchema.properties[propertyKey].format = "data-url"

				this.JSONSchema.properties[propertyKey].items = {}
				break
			default:
				this.JSONSchema.properties[propertyKey].type =
					"string";
				break;
		}
	}


	@action toggleMultipleFilesUpload = (propertyKey)=>{

		this.JSONSchema.properties[propertyKey].showMultiple = !this.JSONSchema.properties[propertyKey].showMultiple
		if (this.JSONSchema.properties[propertyKey].showMultiple){
			this.JSONSchema.properties[propertyKey].type = "array"

			this.JSONSchema.properties[propertyKey].items =  {
				type: "string",
				format: "data-url"
			}
		}else{
			this.JSONSchema.properties[propertyKey].type = "string"
			this.JSONSchema.properties[propertyKey].format = "data-url"

			this.JSONSchema.properties[propertyKey].items = {}
		}
	}
}

export const FormBuilderStore = new formBuilderStore();
