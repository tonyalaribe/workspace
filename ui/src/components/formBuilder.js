import React, { Component } from "react";
import { observer, inject } from "mobx-react";
import {toJS} from "mobx";



@inject("MainStore", "FormBuilderStore")
@observer
class FormBuilder extends Component {
	state = {
		kinds:{},
	}

	render() {
		let { FormBuilderStore } = this.props;
		let {state} = this;

		let otherOptions = function(key, kind){
		 switch (kind) {

			 case "FileUpload":
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
			case "Checkboxes":
				let checkboxes = FormBuilderStore.Checkboxes.get(key)
				console.log(checkboxes)
				let options = checkboxes.map((value, i)=>{
					return (
						<div className="cf pa2" key={i}>
						<input type="text" placeholder={"eg. option "+(i+1)} className="w-80 pv2 ph3 dib fl" defaultValue={value} onInput={(e)=>{FormBuilderStore.setCheckboxOption(key, i, e.target.value)}}/>
						<div className="dib w-20 fl">
							<button className="pv2 ph3 bg-white ba b--light-gray shadow-4" >
								✖
							</button>
						</div>
					</div>
				)
				})
				return (
					<div>
						<div className="pv3">
							<strong>Options</strong>
							<section className="pa3 ">
								{options}
								<div className="cf pa2">
									<div className="dib  fl">
										<button className="pv2 ph3 bg-white ba b--light-gray shadow-4"  onClick={()=>FormBuilderStore.addCheckbox(key)}>+ add</button>
									</div>
								</div>
							</section>
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
									onChange={(e) => {
										FormBuilderStore.updateTitle(key, e.target.value)
									}}
								/>
							</div>
							<div className="w-40 dib fl pa2">
								<select
									className="pv2 ph3 w-100"
									value={FormBuilderStore.JSONSchema.properties[key].kind}
									onChange={e => {
										FormBuilderStore.onFormKindChange(key, e.target.value)
									}}
								>
									<option className="pv2 ph3 w-100" value="ShortAnswer">Short answer</option>
									<option className="pv2 ph3 w-100" value="Paragraph">Paragraph</option>
									<option className="pv2 ph3 w-100" value="Checkboxes">Checkboxes</option>
									<option className="pv2 ph3 w-100" value="list">List</option>
									<option className="pv2 ph3 w-100" value="FileUpload">File upload</option>
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
						onClick={() => FormBuilderStore.addRow()}
					>
						Add Row
					</a>
				</div>

				<section>
					<div className="code pv3">
						{JSON.stringify(toJS(FormBuilderStore.JSONSchema))}
					</div>

					<div className="code pv3">
						{JSON.stringify(toJS(FormBuilderStore.UISchema))}
					</div>

				</section>

			</section>
		);
	}
}

export default FormBuilder;
