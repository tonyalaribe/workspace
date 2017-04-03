import React, { Component } from 'react';
import Nav from '../components/nav.js';
import FileSelect from '../components/fileSelect.js';
import {inject} from "mobx-react";



@inject("MainStore")
class NewSubmissionPage extends Component {
  state = {files:[],showSuccessMessage:false,showErrorMessage:false}
  constructor(props){
    super(props)
    this.newFile = this.newFile.bind(this)
  }

  newFile(file){
    this.setState(
      {
        files:this.state.files.concat(file)
      }
    )
  }

  publishForm(e){
    e.stopPropagation()
    e.preventDefault()

    let formData = {}
    formData.status = "published"
    formData.submissionName = this.refs.submissionName.value
    formData.files = this.state.files


    console.log(formData)
    this.props.MainStore.submitFormToServer(formData,()=>{
      this.setState({showSuccessMessage:true})
      this.refs.submissionName.value = ""
    })

  }
  saveAsDraft(e){
    e.stopPropagation()
    e.preventDefault()

    let formData = {}
    formData.status = "draft"
    formData.submissionName = this.refs.submissionName.value
    formData.files = this.state.files


    console.log(formData)
    this.props.MainStore.submitFormToServer(formData,()=>{

      this.setState({showSuccessMessage:true,files:[]})
      this.refs.submissionName.value = ""

    })

  }


  render() {
    let {state} = this;

    return (
      <section className="">
        <Nav/>
        <section className="tc pt5">
          <section className="pt5 dib w-100 w-70-m w-50-l ">
            <div className="pv3">
              <span className="navy w-100 f3">New Submission</span>
            </div>
            <form className="pv3" >

              <div className="pv3">
                <input type="text" className="pa3 w-100" placeholder="Submission Name" ref="submissionName"/>
              </div>

              <FileSelect newFile={this.newFile} files={state.files}/>

              <div className="pv3">
                {state.showSuccessMessage?<p className="pa3 ba">Submitted Successfully</p>:""}
                {state.showErrorMessage?<p className="pa3 ba">Error In Submission</p>:""}
              </div>

              <div className="pv3 tr">

                <button className="pa3 bg-transparent ba bw1 navy b--navy grow shadow-4  white-80 mh2 pointer"  onClick={this.saveAsDraft.bind(this)}>save as draft</button>

                <button className="pa3 bg-navy grow shadow-4  bw0 white-80 hover-white ml2 pointer"  onClick={this.publishForm.bind(this)}>publish</button>

              </div>

            </form>
          </section>
        </section>
      </section>
    );
  }
}

export default NewSubmissionPage;
