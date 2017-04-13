import React, { Component } from 'react';
import Nav from '../components/nav.js';
// import FileSelect from '../components/fileSelect.js';
import {observer, inject} from "mobx-react";
import moment from "moment";


import Form from "react-jsonschema-form";

//This is a dirty and quick workaround, because using setState prevents the form from submitting.
var STATUS = ""

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

const uiSchema = {
  "firstName": {
    "ui:autofocus": true
  },
  "age": {
    "ui:widget": "updown"
  },
  "bio": {
    "ui:widget": "textarea"
  },
  "password": {
    "ui:widget": "password",
    "ui:help": "Hint: Make it strong!"
  },
  "date": {
    "ui:widget": "alt-datetime"
  }
};

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
    this.props.MainStore.getSubmissionInfo(this.props.match.params.submissionID)
  }


  submitForm(data){
    this.setState({showSuccessMessage:false})
    console.log(this)
    console.log(data)
    let response = {};
    response.status = STATUS;
    response.submissionName = this.refs.submissionName.value;
    response.created = Date.now();
    response.lastModified = Date.now();

    response.formData = data.formData;

    console.log(response)
    console.log(JSON.stringify(response))

    this.props.MainStore.submitFormToServer(response,()=>{

      this.setState({showSuccessMessage:true,files:[]})
      this.refs.submissionName.value = ""

    })

  }

  render() {
    let {state} = this;

    let submissionInfo = this.props.MainStore.SubmissionInfo;
    console.log(submissionInfo)
    console.log(submissionInfo.submissionName)

    return (
      <section className="">
        <Nav/>
        <section className="tc pt5">
          <section className="pt5 dib w-100 w-70-m w-50-l tl">
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
            <Form
              schema={schema}
              uiSchema={uiSchema}
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
                    className="hidden"/>
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
