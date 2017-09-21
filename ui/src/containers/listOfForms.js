import React, { Component } from "react";
import Nav from "../components/nav.js";
import { inject, observer } from "mobx-react";
import { Link } from "react-router-dom";
import Modal from '../components/modals/workspacePermissionsModal.js';

@inject("MainStore")
@observer
class ListOfForms extends Component {
	constructor () {
		super();
		this.state = {
			showModal: false
		};

		this.handleOpenModal = this.handleOpenModal.bind(this);
		this.handleCloseModal = this.handleCloseModal.bind(this);
	}

	componentDidMount() {
		let {workspaceID} = this.props.match.params;

		let {MainStore} = this.props;
    MainStore.getAllForms(workspaceID);
		MainStore.getAllWorkspaces();


	}
	handleOpenModal () {
		this.setState({ showModal: true });
	}

	handleCloseModal () {
		this.setState({ showModal: false });
	}

	componentWillUpdate(nextProps, nextState) {
		if (this.props.location.pathname !== nextProps.location.pathname) {
			this.props.MainStore.getAllForms(nextProps.match.params.workspaceID);
		}
	}
	render() {
		let { MainStore } = this.props;
		let {workspaceID} = this.props.match.params;

		let AllWorkspaces = MainStore.AllWorkspaces.map(function(workspace, key) {
			let workspaceURL = "/workspaces/" + workspace.id;
			return (
				<Link to={workspaceURL} key={key} className="link navy">
					<div
						className={
							" grow pa2 " +
								(window.location.pathname.startsWith(workspaceURL)
									? "bg-blue white-80 "
									: "navy")
						}
					>
						<span className="  ">{workspace.name}</span>
					</div>
				</Link>
			);
		});

		let AllForms = MainStore.AllForms.map(function(form, key) {
			return (
				<Link
					to={"/workspaces/" + workspaceID + "/forms/" + form.id}
					key={key}
					className="link navy"
				>
					<div className="shadow-4 grow mv2 pa3 bg-white">
						<h3 className="navy mv1 ">{form.name}</h3>
						<div>
							<div className=" pv1">
								<small>created by:&nbsp;&nbsp;&nbsp;{form.creator}</small>
							</div>
						</div>
					</div>
				</Link>
			);
		});

		return (
			<section className="">
				<Nav workspaceID={workspaceID} />
				<section className="tc ">
					<section className="pt4 dib w-100 tl cf">
						<div className="w-100 w-25-ns dib v-top ph2 ph3-ns pt4 pb3  pr3 bg-light-gray fixed vh-100">
							<h3 className="bb dib pa1">Workspaces</h3>
							{AllWorkspaces}
						</div><div className="w-100 w-75-ns dib v-top fr pa3-ns mv5">
							<div className="pv3 cf">
								<div className=" fr ">
										<Link
										to={"/workspaces/" + workspaceID + "/new_form"}
										className="ph3 pv2 ba link navy dib grow"
									>
										New Forms
									</Link>
									<a href="#" className="dib link pa2 navy" onClick={this.handleOpenModal}>⚙ &nbsp;settings</a>
									<Modal openModal={this.state.showModal} closeModal={this.handleCloseModal}/>
								</div>
								<span className="navy w-100 v-btm">All Available Forms</span>

							</div>
							<section class="pa3-ns">
								{AllForms}
							</section>
						</div>
					</section>
				</section>
			</section>
		);
	}
}

export default ListOfForms;
