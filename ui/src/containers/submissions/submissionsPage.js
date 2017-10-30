import React, { Component } from "react";
import Nav from "../../components/nav.js";
import { Link } from "react-router-dom";
import { inject, observer } from "mobx-react";
import Modal from "../../components/modals/formSettingsModal.js";
import SubmissionListItem from "./submissionListItem";

@inject("MainStore")
@observer
class SubmissionsPage extends Component {
	constructor() {
		super();
		this.state = {
			showModal: false,
			showDropdown: false
		};

		this.handleOpenModal = this.handleOpenModal.bind(this);
		this.handleCloseModal = this.handleCloseModal.bind(this);
	}

	componentDidMount() {
		let { workspaceID, formID } = this.props.match.params;

		let { MainStore } = this.props;
		MainStore.getAllForms(workspaceID);
		MainStore.getMySubmissions(workspaceID, formID);
	}

	handleOpenModal() {
		this.setState({ showModal: true });
		document.getElementById("body").style.overflow = "hidden";
	}

	handleCloseModal() {
		document.getElementById("body").style.overflow = "scroll";
		this.setState({ showModal: false });
	}

	componentWillUpdate(nextProps, nextState) {
		if (this.props.location.pathname !== nextProps.location.pathname) {
			let { workspaceID, formID } = nextProps.match.params;
			this.props.MainStore.getMySubmissions(workspaceID, formID);
		}
	}

	render() {
		let { workspaceID, formID } = this.props.match.params;
		// let formID = this.props.match.params.formID;
		let { MainStore } = this.props;

		let AllForms = MainStore.AllForms.map(function(form, key) {
			let formURL = "/workspaces/" + workspaceID + "/forms/" + form.id;
			return (
				<Link to={formURL} key={key} className="link navy">
					<div
						className={
							" grow pa2 " +
							(window.location.pathname.startsWith(formURL)
								? "bg-blue white-80"
								: "navy")
						}
					>
						<span className="  ">{form.name}</span>
					</div>
				</Link>
			);
		});

		let userSubmissions = MainStore.Submissions.map(function(fileData, i) {
			return <SubmissionListItem fileData={fileData} key={i} id={i} />;
		});

		return (
			<section>
				<Nav workspaceID={workspaceID} />
				<section className="tc ">
					<section className="pt4 dib w-100 tl cf">
						<div className="w-100 w-25-ns dib v-top ph2 ph3-ns pt4 pb3  pr3 bg-light-gray fixed vh-100">
							<h3 className="bb dib pa1">Forms</h3>
							{AllForms}
						</div>
						<div className="w-100 w-75-ns dib v-top fr pa3-ns mv5">
							<div className="pv3 cf">
								<div className="fr ">
									<Link
										to={
											"/workspaces/" +
											workspaceID +
											"/forms/" +
											formID +
											"/new_submission"
										}
										className="ph3 pv2 ba link navy dib grow"
									>
										New Submission
									</Link>
									<a
										href="#"
										className="dib link pa2 navy"
										onClick={this.handleOpenModal}
									>
										⚙ &nbsp;settings
									</a>
									<Modal
										openModal={this.state.showModal}
										closeModal={this.handleCloseModal}
									/>
								</div>
								<span className="navy w-100 v-btm">All Form Submissions</span>
							</div>

							<section className="pa3-ns">{userSubmissions}</section>
						</div>
					</section>
				</section>
			</section>
		);
	}
}

export default SubmissionsPage;
