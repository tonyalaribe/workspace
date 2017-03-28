import React, { Component } from 'react';
import Nav from '../components/nav.js';
import {inject,observer} from "mobx-react";
import {GetRepresentativeImageByFileExtension } from '../utils/representativeImages.js';


@inject("MainStore") @observer
class SubmissionInfoPage extends Component {
  componentDidMount(){
    this.props.MainStore.getSubmissionInfo(this.props.match.params.submissionID)
  }
  render() {
    let submissionInfo = this.props.MainStore.SubmissionInfo
    console.log(submissionInfo)

    let userFilesCard = submissionInfo.files.map(function(fileData, key){
      console.log(fileData)
      return (<div className="shadow-4 mv2 h4" key={key}>
        {/** Upload Item **/}
        <div className="dib w-30 v-top tc h-100 fl">
          <div className="h-100 flex flex-column  items-center justify-around">
            <img src={GetRepresentativeImageByFileExtension(fileData.name)} className="w3 h3 dib v-mid" alt="file representative logo"/>
          </div>
        </div><div className="dib w-70 h-100 v-top bl b--light-gray pa3">
          <a className="link underline-hover" href={"/"+fileData.path} target="_blank" >
            <h3 className="navy mv1 " >{fileData.name} </h3>
          </a>
          <div><small>Uploaded on: {fileData.uploadDate}</small></div>
        </div>
        {/** End Upload Item **/}
      </div>)
    })

    return (
      <section className="">
        <Nav/>
        <section className="tc pt5">
          <section className="pt4 dib w-100 w-70-m w-50-l tl">
            <div className="pv3">
              <h1 className="navy w-100 mv2">{submissionInfo.submissionName}</h1>
            </div>
            <div className="pv2">
              <strong>status: </strong>
              <span className="navy">{submissionInfo.status}</span>
            </div>
            <div className="pv2">
              <div className="w-100 w-50-ns dib ">
                <small>Created Date: {submissionInfo.uploadDate}</small>
              </div>
              <div className="w-100 w-50-ns dib ">
                <small>Created By: {submissionInfo.createdBy}</small>
              </div>
            </div>
            <section className=" mt5 bt bw1 b--navy">
              <div className="navy tc">
                <h4 className="pv3">Uploaded Files</h4>
              </div>
              {userFilesCard}
            </section>
          </section>
        </section>
      </section>
    );
  }
}

export default SubmissionInfoPage;
