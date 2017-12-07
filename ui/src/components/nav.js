import React, { Component } from "react";
import { Link } from "react-router-dom";
import AuthService from "../utils/auth0.js";
import { observer, inject } from "mobx-react";

@inject("MainStore")
@observer
class Nav extends Component {
	constructor(props, context) {
		super(props, context);
		this.state = {
			showDropdown: false
		};


	}
	componentDidMount() {
		this.props.MainStore.getAllWorkspaces();
	}
	render() {
		let { MainStore, workspaceID } = this.props;
		let currentWorkspace = {};

		if (workspaceID && MainStore.AllWorkspaces.length > 0) {
			currentWorkspace = MainStore.AllWorkspaces.find(function(workspace) {
				return workspace.id === workspaceID;
			});
		}

		let AllWorkspaces = MainStore.AllWorkspaces.map(function(workspace, key) {
			let workspaceURL = "/workspaces/" + workspace.id;
			return (
				<Link
					to={workspaceURL}
					key={key}
					className={
						"db ph4 pv2 link navy hover-bg-light-gray " +
							(window.location.pathname.startsWith(workspaceURL)
								? "bg-light-gray"
								: "")
					}
				>
					<span className="navy  ">{workspace.name}</span>
				</Link>
			);
		});

		console.log(MainStore.profile)
		return (
			<nav className="bg-custom-green w-100 fixed shadow-4 pa3 ph4 white-80 dib z-3">
				<Link to={currentWorkspace.id ? "/workspaces/" + currentWorkspace.name : "/"} className="pa2 dib link white-80 hover-white">
					{currentWorkspace.name ? currentWorkspace.name : "Workspace"}
				</Link>
				<div className="dib fr  mr4">
					<Link
						className="dib pv2 ph4 pointer white-80 link hover-white"
						to="/"
					>

						home
					</Link>
					<div className="dib relative">
						<a
							className="db pa2 pointer"
							onClick={() =>
								this.setState({ showDropdown: !this.state.showDropdown })}
						>
							☰ {(MainStore.profile?MainStore.profile.user_metadata.given_name:"") + " "+(MainStore.profile?MainStore.profile.user_metadata.family_name:"")}
						</a>
						<div
							className={
								"bg-white absolute shadow-4 pv3 right-0 " +
									(this.state.showDropdown ? "dib" : "dn")
							}
							style={{ width: "12rem" }}
						>
							{AllWorkspaces}

							<Link
								to="/"
								onClick={AuthService.logout}
								className="db pv2 ph4 link navy hover-bg-light-gray mt3"
							>
								logout
							</Link>
						</div>
					</div>

				</div>
			</nav>
		);
	}
}

export default Nav;
