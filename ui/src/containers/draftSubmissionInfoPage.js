import React, { Component } from 'react';
import Nav from '../components/nav.js';
// import FileSelect from '../components/fileSelect.js';
import {observer, inject} from "mobx-react";
import moment from "moment";


import Form from "react-jsonschema-form";

//This is a dirty and quick workaround, because using setState prevents the form from submitting.
var STATUS = ""

const log = (type) => console.log.bind(console, type);


function CustomFieldTemplate(props) {
  const {id, classNames, label, help, required, description, errors, children} = props;
  return (
    <div className={classNames+" pv2 tl"}>
      <label htmlFor={id} className="pv2 dib">{label}{required ? "*" : null}</label>
      {description}
      {children}
      {errors}
      {help}
    </div>
  );
}

@inject("MainStore") @observer
class DraftSubmissionInfoPage extends Component {
  state = {files:[],showSuccessMessage:false,showErrorMessage:false}
  constructor(props){
    super(props)
  }

  componentDidMount(){
    this.props.MainStore.getWorkspace(this.props.match.params.workspaceID).then(()=>{
      this.props.MainStore.getSubmissionInfo(this.props.match.params.submissionID)
    })

  }


  submitForm(data){
    this.setState({showSuccessMessage:false})
    console.log(this)
    console.log(data)
    let response = {};
    response.status = STATUS;

    response.lastModified = Date.now();

    response.formData = data.formData;

    console.log(response)
    console.log(JSON.stringify(response))

    this.props.MainStore.updateFormOnServer(this.props.match.params.submissionID,response,()=>{

      this.setState({showSuccessMessage:true})

    })

  }

  render() {
    let {state} = this;

    let {CurrentWorkspace,SubmissionInfo} = this.props.MainStore;

    let jsonschema = CurrentWorkspace.jsonschema;

    return (
      <section className="">
        <Nav/>
        <section className="tc pt5">
          <section className="pt5 dib w-100 w-70-m w-50-l tl">
            <div className="pv3">
              <h1 className="navy w-100 mv2">{SubmissionInfo.submissionName}</h1>
            </div>

            <div className="pv2">
              <strong>status: </strong>
              <span className="navy">{SubmissionInfo.status}</span>
            </div>
            <div className="pv2">
              <div className="w-100 w-50-ns dib ">
                <small>Created: {moment(SubmissionInfo.created).format("h:mma, MM-DD-YYYY")}</small>
              </div>
              <div className="w-100 w-50-ns dib ">
                <small>Modified: {moment(SubmissionInfo.lastModified).format("h:mma, MM-DD-YYYY")}</small>
              </div>
            </div>
            <Form
              schema={CurrentWorkspace.jsonschema}
              uiSchema={CurrentWorkspace.uischema}
              formData={SubmissionInfo.formData}
              onError={log("errors")}
              FieldTemplate={CustomFieldTemplate}
              onSubmit={this.submitForm.bind(this)}
              ref={(form) => {this.form=form;}} >

              <div className="pv3">
                {state.showSuccessMessage?
                  <p className="pa3 ba">
                    Submitted Successfully
                  </p>
                  :""}
                  {state.showErrorMessage?
                    <p className="pa3 ba">
                      Error In Submission
                    </p>
                    :""}
                  </div>

                  <input
                    type="submit"
                    ref={(btn) => {this.submitButton=btn;}}
                    className="hidden dn"/>
                </Form>

                <div className="pv3 tr">
                  <button
                    className="pa3 bg-transparent ba bw1 navy b--navy grow shadow-4  white-80 mh2 pointer"
                    onClick={()=>{STATUS="draft";this.submitButton.click()}}>
                    save as draft
                  </button>

                  <button
                    className="pa3 bg-navy grow shadow-4  bw0 white-80 hover-white ml2 pointer"
                    onClick={()=>{STATUS="publish";this.submitButton.click()}}>publish</button>
                </div>
              </section>
            </section>
          </section>
    );
  }
}

export default DraftSubmissionInfoPage;
