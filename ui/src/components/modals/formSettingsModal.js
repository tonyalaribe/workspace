import React, { Component } from "react";
import { inject, observer } from "mobx-react";
import { withRouter } from "react-router";

@inject("MainStore", "IntegrationsStore")
@observer
class modal extends Component {
	constructor(props) {
		super(props);
		this.AddIntegration = this.AddIntegration.bind(this);
	}
	state = {};
	componentDidMount() {
		// this.props.IntegrationsStore.getWorkspaceUsersAndRoles(
		// 	this.props.match.params.workspaceID
		// );
		this.props.IntegrationsStore.getFormIntegrationSettings(
			this.props.match.params.workspaceID,
			this.props.match.params.formID
		)
	}
	AddIntegration() {
		let Result = {};
    console.log(this.refs)
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
				this.refs.URL.value= "";
				this.refs.SecretToken.value = "";
				this.refs.NewSubmission.checked = false;
				this.refs.UpdateSubmission.checked = false;
				this.refs.DeleteSubmission.checked = false;
				this.refs.ApproveSubmission.checked = false;
				// closeModal();
			}
		);
	}

	render() {
		let { openModal, closeModal, IntegrationsStore } = this.props;
		return (
			<section
				className={
					"vh-100 fixed w-100  justify-center items-center z-4 top-0 left-0 animated " +
					(openModal ? "flex fadeIn" : "dn fadeOut ")
				}
				style={{ backgroundColor: "rgba(0,0,0,0.4)" }}
			>
				<div className="bg-white w-100 w-60-ns modal-shadow giorgia f6 " style={{height:"85%"}} >
					<div className=" bg-light-gray pv2 ph3 shadow-btm ">
						<div className="pv1 cf">
							<strong className="dib v-mid fw4 pv2 ph3">Settings</strong>
							<button
								className="fr dib v-mid pv2 ph3 bg-navy white shadow-4 bw0 grow pointer"
								onClick={closeModal}
							>
								close
							</button>
						</div>
					</div>
					<div className=" cf overflow-y-scroll" style={{height:"85%"}}>
						<div className="w-20 dib fl br1 ph2 pv3">
							<a className="dib pv2 ph3 hover-bg-light-gray w-100">
								integrations
							</a>
						</div>
						<div className="w-80 dib fl pl2 pv3 pr4">
							<section className="cf mb3">
								<div className="mv2">
									<label className="ma0 pv2 pb3 fw6 ph2">URL</label>
									<input type="text" className="pv2 ph3 w-100  mv1" ref="URL" defaultValue={IntegrationsStore.CurrentIntegration.URL}/>
								</div>
								<div className="mv2">
									<label className="ma0 pv2 pb3 fw6 ph2">Secret Token</label>
									<input
										type="text"
										className="pv2 ph3 w-100  mv1"
										ref="SecretToken"
										defaultValue={IntegrationsStore.CurrentIntegration.SecretToken}
									/>
									<p className="gray">
										This Token will be sent with the request in the
										X-Workspace-Token HTTP header.
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
													checked={IntegrationsStore.CurrentIntegration.NewSubmission}

												/>
											<label className="dib" htmlFor="NewSubmission">New Submission</label>
											</div>
											<p className="gray mt1 mb2 pl3">
												A message will be sent to this URL when a new submission
												is made
											</p>
										</div>
										<div className="pv2">
											<div>
												<input
													type="checkbox"
													className="mr2 dib "
													ref="UpdateSubmission"
													id="UpdateSubmission"
													checked={IntegrationsStore.CurrentIntegration.UpdateSubmission}
												/>
											<label className="dib" htmlFor="UpdateSubmission">Update Submission</label>
											</div>
											<p className="gray mt1 mb2 pl3">
												A message will be sent to this URL when a submission is
												updated
											</p>
										</div>
										<div className="pv2">
											<div>
												<input
													type="checkbox"
													className="mr2 dib "
													ref="DeleteSubmission"
													checked={IntegrationsStore.CurrentIntegration.DeleteSubmission}
												/>
												<label className="dib">Delete Submission</label>
											</div>
											<p className="gray mt1 mb2 pl3">
												A message will be sent to this URL when a submission is
												deleted
											</p>
										</div>
										<div className="pv2">
											<div>
												<input
													type="checkbox"
													className="mr2 dib "
													ref="ApproveSubmission"
													id="ApproveSubmission"
													checked={IntegrationsStore.CurrentIntegration.ApproveSubmission}
												/>
											<label className="dib" htmlFor="ApproveSubmission">Approve Submission</label>
											</div>
											<p className="gray mt1 mb2 pl3">
												A message will be sent to this URL when a submission is
												approved
											</p>
										</div>
									</div>
								</div>
								<div className="pv2 fr">
									<button
										className="bg-green grow pv2 ph3 shadow-4 bw0 white-80"
										onClick={() => this.AddIntegration()}
									>
										Add Integration
									</button>
								</div>
							</section>
							<section className="pv2 db">
								<div>
									<h3>web hooks</h3>
								</div>
								<div>
									<div>
										{IntegrationsStore.Integrations.map(function(integration, key) {
											console.log(integration)
											return (
										<div className="pa2 mv2 ba b--light-gray  " key={key}>
											<div className="db cf">
												<strong className="f5 fw5 db ">
													{integration.URL}
												</strong>
											</div>
											<div className="cf pv2">
												<a
													className="ba b--light-gray navy bg-transparent pv1 ph2 link pointer "
													onClick={()=>IntegrationsStore.selectIntegration(integration.ID)}
												>
													Edit
												</a>
												<div className="di">
													<button className="pv1 ph2 ba b--light-gray navy bg-transparent pv1 ph2 link  pointer" >
														Test
													</button>
												</div>

												<a
													data-confirm="Are you sure?"
													className=" link bg-transparent b--light-gray navy pv1 ph2 ba pointer"
													rel="nofollow"
												>
													<span className="">Remove</span>
												</a>
											</div>
										</div>
									)})}
									</div>
								</div>
							</section>
						</div>
					</div>
				</div>
			</section>
		);
	}
}

var Modal = withRouter(modal);
export default Modal;
