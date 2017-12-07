import React, { Component } from "react";
import { inject, observer } from "mobx-react";
import { Link } from "react-router-dom";


@inject("MainStore")
@observer
class ListOfFormsSideNav extends Component {
	componentDidMount() {
    let { workspaceID } = this.props;

		let { MainStore } = this.props;
		MainStore.getAllForms(workspaceID);
  }
	render() {
    let { MainStore } = this.props;
    let { workspaceID, children } = this.props;
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
						<span>{form.name}</span>
					</div>
				</Link>
			);
		});
		return (
			<section className="tc ">
				<section className="pt4 dib w-100 tl cf">
					<div className="w-100 w-25-ns dib v-top ph2 ph3-ns pt4 pb3  pr3 bg-light-gray fixed vh-100">
						<h3 className="bb dib pa1">Forms</h3>
						{AllForms}
					</div>
					<div className="w-100 w-75-ns dib v-top fr pa3-ns mv5">
						{children}
					</div>
				</section>
			</section>
		);
	}
}

export default ListOfFormsSideNav;
