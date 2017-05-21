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
			profile: AuthService.getProfile(),
			showDropdown: false
		};
		console.log(this.state);

		// listen to profile_updated events to update internal state
		AuthService.emitter.on("profile_updated", newProfile => {
			console.log("profile updated");
			this.setState({ profile: newProfile });
		});
	}
	componentDidMount() {
		this.props.MainStore.getAllWorkspaces();
	}
	render() {
		let { MainStore,workspaceID } = this.props;
    let currentWorkspace = {};

		console.log(this.props);
    if (workspaceID&&MainStore.AllWorkspaces.length>0){
      currentWorkspace = MainStore.AllWorkspaces.find(function(workspace){
        console.log(workspace)
        console.log(workspaceID)
        if (workspace.id === workspaceID){
          return workspace
        }
      })
    }
    console.log(currentWorkspace)

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
    console.log(currentWorkspace)

		return (
			<nav className="bg-navy w-100 fixed shadow-4 pa3 ph4 white-80 dib z-3">
				<Link to="/" className="pa2 dib link white-80 hover-white">
					{currentWorkspace.name?currentWorkspace.name:"Workspace"}
				</Link>
				<div className="dib  fr  w5">
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
							☰ {this.state.profile.username}
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
