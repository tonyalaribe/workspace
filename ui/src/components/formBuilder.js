import React, { Component } from "react";
import { observer, inject } from "mobx-react";



@inject("MainStore", "FormBuilderStore")
@observer
class FormBuilder extends Component {
	state = {kinds:{}}
	AddRow() {
		console.log("add row");
		let { FormBuilderStore } = this.props;
		let count = FormBuilderStore.propertiesOrder.length;
		FormBuilderStore.propertiesOrder.push(count);
		FormBuilderStore.JSONSchema.properties[count] = { type: "string",kind:"Short answer" };


	}

	render() {
		let { FormBuilderStore } = this.props;
		let {state} = this;

		let otherOptions = function(key, kind){
		 switch (kind) {

			 case "File upload":
				 return (
					 <div>
						 <div>
							 <label className="dib pointer pa2 navy" htmlFor={"multiple_files_"+key}>
								 <input
									 type="checkbox"
									 name={"multiple_files_"+key}
									 id={"multiple_files_"+key}
									 checked={FormBuilderStore.JSONSchema.properties[key].showMultiple}
									 onChange={(e)=>{
										 console.log("multiple_files_"+key)
										 console.log(e)
										 FormBuilderStore.toggleMultipleFilesUpload(key)
									 }}
								 />{" "}
								 multiple files
							 </label>
						 </div>
					 </div>
				 )
			 default:
				 return ("")
		 }
	 }

		let formFields = FormBuilderStore.propertiesOrder.map(function(key, i) {
			// let property = FormBuilderStore.JSONSchema.properties[key];
			return (
				<div className="ba bw1 b--light-gray hover-grow mv3" key={key}>
					<div>
						<div className="cf w-100">
							<div className="w-60 dib fl pa2">
								<input
									type="text"
									className="pv2 ph3 w-100"
									placeholder="Title"
									defaultValue={
										FormBuilderStore.JSONSchema.properties[key].title
									}
									onChange={e => {
										console.log(e.target.value);
										FormBuilderStore.JSONSchema.properties[key].title =
											e.target.value;
									}}
								/>
							</div>
							<div className="w-40 dib fl pa2">
								<select
									className="pv2 ph3 w-100"
									defaultValue={
										FormBuilderStore.JSONSchema.properties[key].kind
									}
									onChange={e => {
										console.log(e.target.value);
										console.log(state)
										FormBuilderStore.onFormKindChange(key, e.target.value)
									}}
								>
									<option className="pv2 ph3 w-100">Short answer</option>
									<option className="pv2 ph3 w-100">Paragraph</option>
									<option className="pv2 ph3 w-100">Multiple Choice</option>
									<option className="pv2 ph3 w-100">Checkboxes</option>
									<option className="pv2 ph3 w-100">Dropdown</option>
									<option className="pv2 ph3 w-100">File upload</option>
								</select>
							</div>
						</div>
					</div>
					<div>
						{otherOptions(key, FormBuilderStore.Kinds.get(key))}
					</div>
					<div className="bt b--light-gray bw1 cf pa2">
						<div className="fr">
							<a href="#" className="dib link pa2 navy">
								delete
							</a>
						</div>
					</div>
				</div>
			);
		});

		return (
			<section >
				<section>
					{formFields}
				</section>
				<div className="pv2 cf">
					<a
						href="#"
						className="pv2 ph3  link grow bg-light-gray shadow-4 black-90 fr"
						onClick={() => this.AddRow()}
					>
						Add Row
					</a>
				</div>
			</section>
		);
	}
}

export default FormBuilder;
