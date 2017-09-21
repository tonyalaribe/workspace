import React, { Component } from "react";
import { inject, observer } from "mobx-react";
import { withRouter } from "react-router";
import iziToast from "izitoast";

@inject("MainStore", "IntegrationsStore")
@observer
class integrationForm extends Component {
	constructor(props) {
		super(props);
		this.AddIntegration = this.AddIntegration.bind(this);
	}
	AddIntegration() {
		let Result = {};
		console.log(this.refs);
		Result.URL = this.refs.URL.value;
		Result.SecretToken = this.refs.SecretToken.value;
		Result.NewSubmission = this.refs.NewSubmission.checked;
		Result.UpdateSubmission = this.refs.UpdateSubmission.checked;
		Result.DeleteSubmission = this.refs.DeleteSubmission.checked;
		Result.ApproveSubmission = this.refs.ApproveSubmission.checked;

		let { IntegrationsStore, match } = this.props;

		IntegrationsStore.updateFormIntegrationSettings(
			match.params.workspaceID,
			match.params.formID,
			Result,
			() => {
				iziToast.success({
					title: "Update Integrations",
					message: `"Integrations" was updated successfully`,
					position: "topRight"
				});
				if (this.props.clear) {
					this.refs.URL.value = "";
					this.refs.SecretToken.value = "";
					this.refs.NewSubmission.checked = false;
					this.refs.UpdateSubmission.checked = false;
					this.refs.DeleteSubmission.checked = false;
					this.refs.ApproveSubmission.checked = false;
				}

				// closeModal();
				if (this.props.onSave) {
					this.props.onSave();
				}
			}
		);
	}
	render() {
		let { integration, show } = this.props;
		return (
			<section className={"cf mb3 " + (show ? "" : "dn")}>
				<div className="mv2">
					<label className="ma0 pv2 pb3 fw6 ph2">URL</label>
					<input
						type="text"
						className="pv2 ph3 w-100  mv1"
						ref="URL"
						defaultValue={integration.URL}
					/>
				</div>
				<div className="mv2">
					<label className="ma0 pv2 pb3 fw6 ph2">Secret Token</label>
					<input
						type="text"
						className="pv2 ph3 w-100  mv1"
						ref="SecretToken"
						defaultValue={integration.SecretToken}
					/>
					<p className="gray">
						This Token will be sent with the request in the X-Workspace-Token
						HTTP header.
					</p>
				</div>
				<div className="mv2">
					<h4 className="ma0 pv2 pb3 fw6 ph2">Trigger</h4>
					<div>
						<div className="pv2">
							<div>
								<input
									type="checkbox"
									className="mr2 dib "
									ref="NewSubmission"
									id="NewSubmission"
									defaultChecked={integration.NewSubmission}
								/>
								<label className="dib" htmlFor="NewSubmission">
									New Submission
								</label>
							</div>
							<p className="gray mt1 mb2 pl3">
								A message will be sent to this URL when a new submission is made
							</p>
						</div>
						<div className="pv2">
							<div>
								<input
									type="checkbox"
									className="mr2 dib "
									ref="UpdateSubmission"
									id="UpdateSubmission"
									defaultChecked={integration.UpdateSubmission}
								/>
								<label className="dib" htmlFor="UpdateSubmission">
									Update Submission
								</label>
							</div>
							<p className="gray mt1 mb2 pl3">
								A message will be sent to this URL when a submission is updated
							</p>
						</div>
						<div className="pv2">
							<div>
								<input
									type="checkbox"
									className="mr2 dib "
									ref="DeleteSubmission"
									defaultChecked={integration.DeleteSubmission}
								/>
								<label className="dib">Delete Submission</label>
							</div>
							<p className="gray mt1 mb2 pl3">
								A message will be sent to this URL when a submission is deleted
							</p>
						</div>
						<div className="pv2">
							<div>
								<input
									type="checkbox"
									className="mr2 dib "
									ref="ApproveSubmission"
									id="ApproveSubmission"
									defaultChecked={integration.ApproveSubmission}
								/>
								<label className="dib" htmlFor="ApproveSubmission">
									Approve Submission
								</label>
							</div>
							<p className="gray mt1 mb2 pl3">
								A message will be sent to this URL when a submission is approved
							</p>
						</div>
					</div>
				</div>
				<div className="pv2 fr">
					<button
						className="bg-green grow pv2 ph3 shadow-4 bw0 white-80"
						onClick={() => this.AddIntegration()}
					>
						Save Integration
					</button>
				</div>
			</section>
		);
	}
}

export default withRouter(integrationForm);
