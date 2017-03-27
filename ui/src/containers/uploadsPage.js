import React, { Component } from 'react';
import Nav from '../components/nav.js';
import AuthService from '../utils/auth0.js';
import {inject,observer} from "mobx-react";
import {GetRepresentativeImageByFileExtension } from '../utils/representativeImages.js';

@inject("MainStore") @observer
class UploadsPage extends Component {
  componentDidMount(){
    this.props.MainStore.getMySubmissions()
  }
  render() {
    let submissions = this.props.MainStore.Submissions

    let userFilesCard = submissions.map(function(fileData, key){
      console.log(fileData)
      return (<div className="shadow-4 mv2 h4" key={key}>
        {/** Upload Item **/}
        <div className="dib w-30 v-top tc h-100 fl">
          <div className="h-100 flex flex-column  items-center justify-around">
            <img src={GetRepresentativeImageByFileExtension(fileData.FileName)} className="w3 h3 dib v-mid" alt="file representative logo"/>
          </div>
        </div><div className="dib w-70 h-100 v-top bl b--light-gray pa3">
          <h3 className="navy mv1 ">{fileData.SubmissionName} </h3>
          <div><small>File Name: {fileData.FileName}</small></div>
          <div><small>status: {fileData.Status}</small></div>
          <div><small>Uploaded By: {fileData.CreatedBy}</small></div>

        </div>
        {/** End Upload Item **/}
      </div>)

    })
    console.log(AuthService.loggedIn())
    return (
      <section className="">
        <Nav/>
        <section className="tc pt5">
          <section className="pt4 dib w-100 w-70-m w-50-l tl">
            <div className="pv3">
              <span className="navy w-100">All Uploaded Files</span>
            </div>
            <section>
              {userFilesCard}
            </section>
          </section>
        </section>
      </section>
    );
  }
}

export default UploadsPage;
