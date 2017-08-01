import React, { Component } from "react";
import Nav from "../components/nav.js";
import FormBuilder from "../components/formBuilder.js";
import { inject, observer } from "mobx-react";
import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import 'react-tabs/style/react-tabs.css';
import { toJS } from "mobx";

import FileWidget from "../components/fileWidget.js";
import Form from "react-jsonschema-form";


const log = type => console.log.bind(console, type);
const widgets = {
	FileWidget: FileWidget
};

function CustomFieldTemplate(props) {
	const {
		id,
		classNames,
		label,
		help,
		required,
		description,
		errors,
		children
	} = props;
	return (
		<div className={classNames + " pv2 tl"}>
			<label htmlFor={id} className="pv2 dib">
				{label}{required ? "*" : null}
			</label>
			{description}
			{children}
			{errors}
			{help}
		</div>
	);
}


@inject("MainStore", "FormBuilderStore")
@observer
class NewFormPage extends Component {
	state = {};

	submitFormToServer() {
		console.log("submitFormToServer 1 ")
		let workspaceID = this.props.match.params.workspaceID;
		this.setState({ showSuccessMessage: false });
		console.log("submitNewFormToServer")

		this.props.FormBuilderStore.submitNewFormToServer(workspaceID,() => {
			this.setState({ showSuccessMessage: true });
			this.refs.formName.value = "";
			this.refs.jsonSchema.value = "";
			this.refs.uiSchema.value = "";
			setTimeout(() => {
				window.requestAnimationFrame(() => {
					this.props.history.push("/workspaces/" + workspaceID);
				});
			}, 1000);
		});
	}
	render() {
		let { state } = this;
		let { workspaceID } = this.props.match.params;
		let {FormBuilderStore} = this.props;

		let jsonschema = toJS(FormBuilderStore.JSONSchema)
		let uischema = toJS(FormBuilderStore.UISchema)
		console.log(uischema)
		return (
			<section>
				<Nav workspaceID={workspaceID} />
				<section className="tc pt5">
					<section className="pt5 dib w-100 w-70-m w-50-l ">
						<div className="pt3 pb5 ">
							<span className="navy w-100 f2 db">New Form</span>
						</div>

						<section>
							<Tabs>
						    <TabList>
						      <Tab>Builder</Tab>
						      <Tab>Preview</Tab>
						    </TabList>
						    <TabPanel>
									<section>
										<div className="pv3 tl">
											<label className="pv2 dib">Form Name</label>
											<input type="text" className="form-control " ref="formName" defaultValue={FormBuilderStore.JSONSchema.title} onChange={(e)=>{
													console.log(e.target.value)
													FormBuilderStore.JSONSchema.title = e.target.value
												}}/>
										</div>
							      <FormBuilder />
									</section>
						    </TabPanel>
						    <TabPanel>
									<section>
											<Form
												schema={jsonschema}
												uiSchema={uischema}
												onError={log("errors")}
												FieldTemplate={CustomFieldTemplate}
												widgets={widgets}
												ref={form => {
													this.form = form;
												}}
											>
												<input
													type="submit"
													ref={btn => {
														this.submitButton = btn;
													}}
													className="hidden dn"
												/>
											</Form>
										</section>
						    </TabPanel>
						  </Tabs>

						</section>
						<div className="pv3">
							{state.showSuccessMessage
								? <p className="pa3 ba">Submitted Successfully</p>
								: ""}
							{state.showErrorMessage
								? <p className="pa3 ba">Error In Submission</p>
								: ""}
						</div>

						<div className="pv3 tr">
							<button
								className="pa3 bg-navy grow shadow-4  bw0 white-80 hover-white ml2 pointer"
								onClick={()=>{console.log("DF");this.submitFormToServer}}
							>
								Create Form
							</button>
						</div>
					</section>
				</section>
			</section>
		);
	}
}

export default NewFormPage;
