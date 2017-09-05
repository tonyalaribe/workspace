import React, { Component } from "react";
import { inject, observer } from "mobx-react";
import { withRouter } from "react-router";
import IntegrationsListItem from "./integrationsListItem.js";
import IntegrationForm from "./integrationForm.js";

@inject("MainStore", "IntegrationsStore")
@observer
class modal extends Component {

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
							<IntegrationForm integration={IntegrationsStore.CurrentIntegration} show={true}/>
							<section className="pv2 db">
								<div>
									<h3>web hooks</h3>
								</div>
								<div>
									<div>
										{IntegrationsStore.Integrations.map(function(item, i) {
											return (
												<IntegrationsListItem  integration={item} key={i} />
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
