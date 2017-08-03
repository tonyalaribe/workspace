import { observable, action, runInAction, toJS } from "mobx";
import AuthService from "../utils/auth0.js";

class formBuilderStore {
	@observable propertiesOrder = [];
	@observable
	JSONSchema = {
		type: "object",
		properties: {}
	};
	@observable UISchema = {};
	@observable Kinds = new Map();
	@observable Checkboxes = new Map()

	@action
	submitNewFormToServer = async (workspaceID, callback) => {
		let authToken = AuthService.getToken();

		let form = {};
		form.name = this.JSONSchema.title;

		form.jsonschema = toJS(this.JSONSchema);
		form.uischema = toJS(this.UISchema);

		const response = await fetch(
			"/api/workspaces/" + workspaceID + "/new_form",
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

	@action addRow = ()=>{
		let count = this.propertiesOrder.length;
		this.propertiesOrder.push(count);
		this.JSONSchema.properties[count] = {
			type: "string",
			kind:"ShortAnswer"
		};
	}

	@action updateTitle = (propertyKey, title)=>{
		this.JSONSchema.properties[propertyKey].title = title;
	}

	@action
	onFormKindChange = (propertyKey, kind) => {
		this.JSONSchema.properties[propertyKey].kind = kind;
		this.Kinds.set(propertyKey, kind);
		console.log(this.Kinds);
		switch (kind) {
			case "ShortAnswer":
				this.JSONSchema.properties[propertyKey].type = "string";
				break;
			case "Paragraph":
				this.JSONSchema.properties[propertyKey].type = "string";
				this.UISchema[propertyKey] = {
					"ui:widget":"textarea"
				};
				break;
			case "FileUpload":
				this.JSONSchema.properties[propertyKey].type = "string";
				this.JSONSchema.properties[propertyKey].format = "data-url";

				this.JSONSchema.properties[propertyKey].items = {};
				break;
			case "Checkboxes":
				this.JSONSchema.properties[propertyKey].type = "array";
				this.JSONSchema.properties[propertyKey].items = {
					type:"string",
				};
				this.UISchema[propertyKey] = {
					"ui:widget":"checkboxes"
				};
				this.Checkboxes.set(propertyKey, [""])
				break;
			default:
				this.JSONSchema.properties[propertyKey].type = "string";
				break;
		}
	};

	@action
	toggleMultipleFilesUpload = propertyKey => {
		this.JSONSchema.properties[propertyKey].showMultiple = !this.JSONSchema
			.properties[propertyKey].showMultiple;
		if (this.JSONSchema.properties[propertyKey].showMultiple) {
			this.JSONSchema.properties[propertyKey].type = "array";

			this.JSONSchema.properties[propertyKey].items = {
				type: "string",
				format: "data-url"
			};
		} else {
			this.JSONSchema.properties[propertyKey].type = "string";
			this.JSONSchema.properties[propertyKey].format = "data-url";

			this.JSONSchema.properties[propertyKey].items = {};
		}
	};

	@action addCheckbox = (propertyKey)=>{
		let propertyCheckboxes = this.Checkboxes.get(propertyKey)
		propertyCheckboxes.push("")
		this.Checkboxes.set(propertyKey, propertyCheckboxes)
	};

	@action setCheckboxOption = (propertyKey, checkboxKey, checkboxValue)=>{
		console.log(checkboxValue)
		let propertyCheckboxes = this.Checkboxes.get(propertyKey)
		propertyCheckboxes[checkboxKey] = checkboxValue
		this.Checkboxes.set(propertyKey, propertyCheckboxes)
		console.log(this.Checkboxes.get(propertyKey))

	};


}

export const FormBuilderStore = new formBuilderStore();
