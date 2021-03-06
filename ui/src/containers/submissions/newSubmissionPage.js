import React, { Component } from "react";
import Nav from "../../components/nav.js";
import { inject, observer } from "mobx-react";
import { toJS } from "mobx";
import iziToast from "izitoast";
import Form from "react-jsonschema-form";
import { CustomFieldTemplate, widgets, log } from "./JSONSchemaFormsHelper.js";
import ListOfFormsSideNav from "../ListOfFormsSideNav.js";

//This is a dirty and quick workaround, because using setState prevents the form from submitting.
var STATUS = "";

@inject("MainStore")
@observer
class NewSubmissionPage extends Component {
	state = { files: [] };

	componentDidMount() {
		this.props.MainStore.getFormInfo(
			this.props.match.params.workspaceID,
			this.props.match.params.formID
		);
	}

	submitForm(data) {
		let { workspaceID, formID } = this.props.match.params;

		this.setState({ showSuccessMessage: false });
		let response = {};
		response.status = STATUS;
		response.submissionName = this.refs.submissionName.value;
		response.created = Date.now();
		response.lastModified = Date.now();
		response.formData = data.formData;

		this.props.MainStore.submitFormToServer(
			workspaceID,
			formID,
			response,
			resp => {
				this.setState({ showSuccessMessage: true, files: [] });
				iziToast.success({
					title: "New Submission",
					message: `"${response.submissionName}" was created successfully`,
					position: "topRight"
				});

				this.refs.submissionName.value = "";
				setTimeout(() => {
					window.requestAnimationFrame(() => {
						this.props.history.push(
							"/workspaces/" +
								workspaceID +
								"/forms/" +
								formID +
								"/submissions/" +
								response.status +
								"/" +
								resp.submission.id
						);
					});
				}, 1000);
			}
		);
	}

	render() {
		let { state } = this;
		let { CurrentForm } = this.props.MainStore;
		let { workspaceID } = this.props.match.params;
		return (
			<section>
				<Nav workspaceID={workspaceID} />
				<ListOfFormsSideNav workspaceID={workspaceID}>
					<section className="tc ">
						<section className="pt5 dib w-100 w-80-ns ">
							<div className="pv3">
								<span className="navy w-100 f3 db">New Submission</span>
								<span className="db">
									{CurrentForm.name ? "(" + CurrentForm.name + ")" : ""}
								</span>
							</div>
							<div className="pv3 tl">
								<label className="pv2 dib">Submission Name</label>
								<input
									type="text"
									className="form-control"
									ref="submissionName"
								/>
							</div>
							<Form
								schema={toJS(CurrentForm.jsonschema)}
								uiSchema={toJS(CurrentForm.uischema)}
								onError={log("errors")}
								FieldTemplate={CustomFieldTemplate}
								onSubmit={this.submitForm.bind(this)}
								widgets={widgets}
								ref={form => {
									this.form = form;
								}}
							>
								<div className="pv3">
									{state.showSuccessMessage ? (
										<p className="pa3 ba">Submitted Successfully</p>
									) : (
										""
									)}
									{state.showErrorMessage ? (
										<p className="pa3 ba">Error In Submission</p>
									) : (
										""
									)}
								</div>

								<input
									type="submit"
									ref={btn => {
										this.submitButton = btn;
									}}
									className="hidden dn"
								/>
							</Form>

							<div className="pv3 tr">
								<button
									className="pa3 bg-transparent ba bw1 navy b--navy grow shadow-4  white-80 mh2 pointer"
									onClick={() => {
										STATUS = "draft";
										this.submitButton.click();
									}}
								>
									save as draft
								</button>

								<button
									className="pa3 bg-navy grow shadow-4  bw0 white-80 hover-white ml2 pointer"
									onClick={() => {
										STATUS = "published";
										this.submitButton.click();
									}}
								>
									publish
								</button>
							</div>
						</section>
					</section>
				</ListOfFormsSideNav>
			</section>
		);
	}
}

export default NewSubmissionPage;
