import React, { Component } from "react";
import { Link } from "react-router-dom";
import { inject, observer } from "mobx-react";
import moment from "moment";
var fileImageRepresentation = require("../../assets/files.png");
import { withRouter } from "react-router";

@inject("MainStore")
@observer
class submissionListItem extends Component {
	constructor() {
		super();
		this.state = {
			showModal: false,
			showDropdown: false
		};
	}

	render() {
		let { fileData, MainStore, id } = this.props;
		let { workspaceID,formID } = this.props.match.params;

		return (
			<Link
				to={
					"/workspaces/" +
					workspaceID +
					"/forms/" +
					formID +
					"/submissions/" +
					fileData.status +
					"/" +
					fileData.id
				}
				key={id}
				className="link navy"
			>
				<div className="shadow-4 hover-grow mv2 " style={{height:"10rem"}}>
					{/** Upload Item **/}
					<div className="dib w-30 v-top tc h-100 fl">
						<div className="h-100 flex flex-column  items-center justify-around">
							<img
								src={fileImageRepresentation}
								className="w3 h3 dib v-mid"
								alt="file representative logo"
							/>
						</div>
					</div>
					<div className="dib w-70 h-100 v-top bl b--light-gray pa3">
						<div>
							<small className="fr pa2 bg-navy white-80">
								{fileData.status}
							</small>
						</div>
						<h3 className="navy mv1 ">
							{fileData.submissionName}{" "}
						</h3>
						<div>
							<div className=" pv1">
								<small>
									created on:&nbsp;&nbsp;&nbsp;
									{moment(fileData.created).format("h:mma, MM-DD-YYYY")}
								</small>
							</div>
							<div className=" pv1">
								<small>
									modified on:&nbsp;&nbsp;
									{moment(fileData.lastModified).format("h:mma, MM-DD-YYYY")}
								</small>
							</div>
						</div>
            <div className="db mt3">
  						<a className="pa2 dib ba b--light-gray fr">
  							<span className="dib " onClick={
										(e)=>{
											e.preventDefault()
											MainStore.deleteSubmission(workspaceID, formID, fileData.id,id,function(){})
										}
									}>delete</span>
  						</a>
            </div>

					</div>
					{/** End Upload Item **/}
				</div>
			</Link>
		);
	}
}


var SubmissionListItem = withRouter(submissionListItem);
export default SubmissionListItem;
