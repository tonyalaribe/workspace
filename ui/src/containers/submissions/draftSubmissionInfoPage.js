import React, { Component } from "react";
import Nav from "../../components/nav.js";
import FileWidget from "../../components/fileWidget.js";
import { toJS } from "mobx";
import { observer, inject } from "mobx-react";
import moment from "moment";
import { Tab, Tabs, TabList, TabPanel } from "react-tabs";
import iziToast from "izitoast";
import Form from "react-jsonschema-form";

//This is a dirty and quick workaround, because using setState prevents the form from submitting.
var STATUS = "";

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
				{label}
				{required ? "*" : null}
			</label>
			{description}
			{children}
			{errors}
			{help}
		</div>
	);
}

@inject("MainStore")
@observer
class DraftSubmissionInfoPage extends Component {
	state = {
		files: [],
		showSuccessMessage: false,
		showErrorMessage: false,
		submissionNotes: ""
	};

	componentDidMount() {
		let { workspaceID, formID, submissionID } = this.props.match.params;
		this.props.MainStore.getFormInfo(workspaceID, formID).then(() => {
			this.props.MainStore.getSubmissionInfo(workspaceID, formID, submissionID);
			this.props.MainStore.getSubmissionChangelog(
				workspaceID,
				formID,
				submissionID
			);
		});
	}

	submitForm(data) {
		let { workspaceID, formID, submissionID } = this.props.match.params;

		this.setState({ showSuccessMessage: false });

		let req = {};
		req.status = STATUS;
		req.lastModified = Date.now();
		req.formData = data.formData;
		req.submissionNotes = this.state.submissionNotes;

		this.props.MainStore.SubmissionInfo = req; //To prevent reverting to old value on form submit.

		this.props.MainStore.updateFormOnServer(
			workspaceID,
			formID,
			submissionID,
			req,
			() => {
				this.setState({ showSuccessMessage: true });
				iziToast.success({
					title: "Update Submission",
					message: `Submission was updated successfully`,
					position: "topRight"
				});
				setTimeout(() => {
					this.props.MainStore.getFormInfo(workspaceID, formID).then(() => {
						this.props.MainStore.getSubmissionInfo(workspaceID, formID, submissionID);
						this.props.MainStore.getSubmissionChangelog(
							workspaceID,
							formID,
							submissionID
						);
					});

					window.requestAnimationFrame(() => {
						this.props.history.push(
							"/workspaces/" +
								workspaceID +
								"/forms/" +
								formID +
								"/submissions/" +
								req.status +
								"/" +
								submissionID
						);
					});
				}, 1000);
			}
		);
	}

	render() {
		let { state } = this;

		let { CurrentForm, SubmissionInfo, Changelog } = this.props.MainStore;
		let { workspaceID } = this.props.match.params;
		return (
			<section className="">
				<Nav workspaceID={workspaceID} />
				<section className="tc pt5">
					<section className="pt5 dib w-100 w-70-m w-50-l tl">
						<div className="pv3">
							<h1 className="navy w-100 mv2">
								{SubmissionInfo.submissionName}
							</h1>
						</div>

						<div className="pv2">
							<strong>status: </strong>
							<span className="navy">{SubmissionInfo.status}</span>
						</div>
						<div className="pv2 ">
							<div className="w-100 w-50-ns dib ">
								<small>
									Created:{" "}
									{moment(SubmissionInfo.created).format("h:mma, MM-DD-YYYY")}
								</small>
							</div>
							<div className="w-100 w-50-ns dib ">
								<small>
									Modified:{" "}
									{moment(SubmissionInfo.lastModified).format(
										"h:mma, MM-DD-YYYY"
									)}
								</small>
							</div>
						</div>

						<Tabs className="pt4">
							<TabList>
								<Tab>Form</Tab>
								<Tab>Changelog</Tab>
							</TabList>
							<TabPanel>
								<section>
									{/* Form does not rerender when after the formdata is retrieved from the server, so its now only being rendered after the data is retrieved from the server */}
									{SubmissionInfo.formData ? (
										<Form
											schema={toJS(CurrentForm.jsonschema)}
											uiSchema={toJS(CurrentForm.uischema)}
											formData={toJS(SubmissionInfo.formData)}
											onError={log("errors")}
											FieldTemplate={CustomFieldTemplate}
											onSubmit={this.submitForm.bind(this)}
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
									) : (
										""
									)}

									<div className="pv2">
										<label className="pv2">Submission Notes</label>
										<textarea
											rows="5"
											className="w-100 mv2 "
											onChange={e =>
												this.setState({ submissionNotes: e.target.value })}
										/>
									</div>

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
							</TabPanel>
							<TabPanel>
								<section className="pv2">
									{Changelog.map(function(changelogItem, i) {
										return (
											<div className="w-100 shadow-4 pa3 mv2" key={i}>
												<p className="navy mv1 ">{changelogItem.note}</p>
												<div>
													<div className=" pv1">
														<small>
															created on:&nbsp;&nbsp;&nbsp;
															{moment(changelogItem.created).format(
																"h:mma, MM-DD-YYYY"
															)}
														</small>
													</div>
													<div className=" pv1">
														<small>
															workspace: {changelogItem.workspaceID}
														</small>
													</div>
												</div>
											</div>
										);
									})}
								</section>
							</TabPanel>
						</Tabs>
					</section>
				</section>
			</section>
		);
	}
}

export default DraftSubmissionInfoPage;
