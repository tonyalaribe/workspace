import React, { Component } from "react";
import Nav from "../../components/nav.js";
import { inject, observer } from "mobx-react";
import moment from "moment";
import { Tab, Tabs, TabList, TabPanel } from "react-tabs";
import { GetRepresentativeImageByFileExtension } from "../../utils/representativeImages.js";

@inject("MainStore")
@observer
class PublishedSubmissionInfoPage extends Component {
	state = {};

	componentDidMount() {
		let { workspaceID, formID, submissionID } = this.props.match.params;
		this.props.MainStore.getSubmissionInfo(
			workspaceID,
			formID,
			submissionID
		).then(() => {
			this.props.MainStore.getFormInfo(workspaceID, formID).then(() => {
				this.props.MainStore.getSubmissionChangelog(
					workspaceID,
					formID,
					submissionID
				);
			});
		});
	}

	render() {
		let { CurrentForm, SubmissionInfo, Changelog } = this.props.MainStore;
		let jsonschema = CurrentForm.jsonschema;

		let formFields = Object.keys(jsonschema.properties).reduce(
			(previous, current) => {
				let value;
				switch (jsonschema.properties[current].type) {
					case "string":
						switch (jsonschema.properties[current].format) {
							case "data-url":
								let currentURI = SubmissionInfo.formData[current];
								let currentURISplit = currentURI.split("/");
								let currentURIFileName =
									currentURISplit[currentURISplit.length - 1];
								value = (
									<div className="pr1 pb1 w-100 ">
										<div className="h-100 ba b--black-20  pa2">
											<a
												className="h-100 dib link pointer  pb1 navy pa1 w-25 tc br b--light-gray"
												href={"/uploads/" + currentURI}
												target="_blank"
											>
												<img
													src={GetRepresentativeImageByFileExtension(currentURIFileName)}
													className="w3 h3 dib v-mid"
													alt="file representative logo"
												/>
											</a>
											<a
												className="fr mt4 dib link pointer truncate pb1  navy pa1 pl3 w-75"
												href={"/uploads/" + currentURI}
												target="_blank"
											>
												{currentURIFileName}
											</a>
										</div>
									</div>
								);
								break;
							default:
								value = SubmissionInfo.formData[current];
								break;
						}
						break;
					case "array":
						if (SubmissionInfo.formData[current]) {
							value = SubmissionInfo.formData[current].map(function(item, i) {
								switch (jsonschema.properties[current].items.type) {
									case "string":
										switch (jsonschema.properties[current].items.format) {
											case "data-url":
												let currentURI = item;
												let currentURISplit = currentURI.split("/");
												let currentURIFileName =
													currentURISplit[currentURISplit.length - 1];
												return (
													<div className="pr1 pb1 w-100 " key={i}>
														<div className="h-100 ba b--black-20  pa2">
															<a
																className="h-100 dib link pointer  pb1 navy pa1 w-25 tc br b--light-gray"
																href={"/uploads/" + currentURI}
																target="_blank"
															>
																<img
																	src={GetRepresentativeImageByFileExtension(
																		currentURIFileName
																	)}
																	className="w3 h3 dib v-mid"
																	alt="file representative logo"
																/>
															</a>
															<a
																className="fr mt4 dib link pointer truncate pb1  navy pa1 pl3 w-75"
																href={"/uploads/" + currentURI}
																target="_blank"
															>
																{currentURIFileName}
															</a>
														</div>
													</div>
												);
											default:
												return item;
										}
									default:
										return item;
								}
							});
						}
						break;
					default:
						value = SubmissionInfo.formData[current];
						break;
				}
				previous.push(
					<div className="pv2" key={current}>
						<strong className="pa1 dib">
							{jsonschema.properties[current].title}: &nbsp;&nbsp;
						</strong>
						<div className="pa1 dib w-100">{value}</div>
					</div>
				);
				return previous;
			},
			[]
		);

		let { workspaceID } = this.props.match.params;
		return (
			<section className="">
				<Nav workspaceID={workspaceID} />
				<section className="tc pt5">
					<section className="pt4 dib w-100 w-70-m w-50-l tl">
						<div className="pv3">
							<h1 className="navy w-100 mv2">
								{SubmissionInfo.submissionName}
							</h1>
						</div>

						<div className="pv2">
							<strong>status: </strong>
							<span className="navy">{SubmissionInfo.status}</span>
						</div>
						<div className="pv2">
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

						<section className="mt5">
							<Tabs className="pt4">
								<TabList>
									<Tab>Form</Tab>
									<Tab>Changelog</Tab>
								</TabList>
								<TabPanel>
									<div className="navy tc bb bw1 b--light-gray pv3">
										<h4 className="mv3">Form Data</h4>
									</div>
									<div className=" ph2 pv3 ">{formFields}</div>
								</TabPanel>
								<TabPanel>
									<section className="pv2">
										{Changelog.map(function(changelogItem,i) {
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
			</section>
		);
	}
}

export default PublishedSubmissionInfoPage;
