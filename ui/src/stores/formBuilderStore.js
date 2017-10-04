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
			"/api/workspaces/" + workspaceID + "/forms",
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
			console.log(data)
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

	@action deleteRow = (propertyKey)=>{
		this.propertiesOrder.splice(this.propertiesOrder.indexOf(propertyKey),1);
		delete this.JSONSchema.properties[propertyKey];
	}


	@action updateTitle = (propertyKey, title)=>{
		this.JSONSchema.properties[propertyKey].title = title;
	}

	@action
	onFormKindChange = (propertyKey, kind) => {
		this.JSONSchema.properties[propertyKey].kind = kind;
		this.Kinds.set(propertyKey, kind);
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
					enum:[],
				};
				this.JSONSchema.properties[propertyKey].uniqueItems = true
				this.UISchema[propertyKey] = {
					"ui:widget":"checkboxes"
				};
				this.Checkboxes.set(propertyKey, [""])
				break;
			case "List":
				this.JSONSchema.properties[propertyKey].type = "array";
				this.JSONSchema.properties[propertyKey].items = {
					type:"string",
				};
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
			this.JSONSchema.properties[propertyKey].items.enum.push("")
		}
	};

	@action addCheckbox = (propertyKey)=>{
		let propertyCheckboxes = this.Checkboxes.get(propertyKey)
		propertyCheckboxes.push("")
		this.Checkboxes.set(propertyKey, propertyCheckboxes)

		this.JSONSchema.properties[propertyKey].items.enum.push("a")
		console.log(this.JSONSchema.properties[propertyKey].items.enum)
	};

	@action deleteCheckbox = (propertyKey, checkboxKey) =>{
		let propertyCheckboxes = this.Checkboxes.get(propertyKey)
		propertyCheckboxes.splice(checkboxKey, 1)
		this.Checkboxes.set(propertyKey, propertyCheckboxes)
		this.JSONSchema.properties[propertyKey].items.enum.splice(checkboxKey, 1)
	}

	@action setCheckboxOption = (propertyKey, checkboxKey, checkboxValue)=>{
		let propertyCheckboxes = this.Checkboxes.get(propertyKey)
		propertyCheckboxes[checkboxKey] = checkboxValue
		this.Checkboxes.set(propertyKey, propertyCheckboxes)

		this.JSONSchema.properties[propertyKey].items.enum[checkboxKey] = checkboxValue

	};


}

export const FormBuilderStore = new formBuilderStore();
