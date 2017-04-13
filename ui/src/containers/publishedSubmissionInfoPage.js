import React, { Component } from 'react';
import Nav from '../components/nav.js';
import {inject,observer} from "mobx-react";
import {GetRepresentativeImageByFileExtension } from '../utils/representativeImages.js';
import moment from "moment";

const schema = {
  "type": "object",
  "required": [
    "firstName",
    "lastName"
  ],
  "properties": {
    "firstName": {
      "type": "string",
      "title": "First name"
    },
    "lastName": {
      "type": "string",
      "title": "Last name"
    },
    "age": {
      "type": "integer",
      "title": "Age"
    },
    "bio": {
      "type": "string",
      "title": "Bio"
    },
    "password": {
      "type": "string",
      "title": "Password",
      "minLength": 3
    }
  }
};


@inject("MainStore") @observer
class PublishedSubmissionInfoPage extends Component {
  state = {}

  constructor(props){
    super(props)
  }

  componentDidMount(){
    this.props.MainStore.getSubmissionInfo(this.props.match.params.submissionID)
  }
  render() {
    let {state} = this;
    let submissionInfo = this.props.MainStore.SubmissionInfo;

    let formFields = Object.keys(schema.properties).reduce((previous,current)=>{
      previous.push(<div className="pv2" key={current}>
        <strong className="pa1 dib">{schema.properties[current].title}: &nbsp;&nbsp;</strong>
        <span className="pa1 dib">{submissionInfo.formData[current]}</span>
      </div>)
      return previous
    },[]);

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
                <small>Created: {moment(submissionInfo.created).format("h:mma, MM-DD-YYYY")}</small>
              </div>
              <div className="w-100 w-50-ns dib ">
                <small>Modified: {moment(submissionInfo.lastModified).format("h:mma, MM-DD-YYYY")}</small>
              </div>
            </div>

            <section className="mt5">
              <div className="navy tc bb bw1 b--navy ">
                <h4 className="mv3">Form Data</h4>
              </div>
              <div className=" ph2 pv3 ">
                {formFields}
              </div>
            </section>

          </section>
        </section>
      </section>
    );
  }
}

export default PublishedSubmissionInfoPage;
