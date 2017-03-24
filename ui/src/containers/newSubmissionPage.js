import React, { Component } from 'react';
import Nav from '../components/nav.js';
import {inject} from "mobx-react";
import {GetFileRepresentativeImage } from '../utils/representativeImages.js';



@inject("MainStore")
class NewSubmissionPage extends Component {
  state = {files:[],showSuccessMessage:false,showErrorMessage:false}
  constructor(props){
    super(props)
    this.removeFile = this.removeFile.bind(this)
  }
  submitForm(e){
    e.stopPropagation()
    e.preventDefault()

    let formData = {}
    formData.submissionName = this.refs.submissionName.value
    formData.files = this.state.files

    console.log(formData)
    this.props.MainStore.submitFormToServer(formData,()=>{
      this.setState({showSuccessMessage:true})
    })

  }
  FileSelectHandler(e){
    var files = e.target.files || e.dataTransfer.files;
  	for (let i = 0;i<files.length ; i++) {
      let file = files[i]
      var reader = new FileReader();
      reader.onload = (()=> {
         return (e)=> {
           this.setState({files:this.state.files.concat({file:e.target.result,name:file.name,type:file.type})})
         };
       })();
      reader.readAsDataURL(file);
  	}
  }
  removeFile(key){
    console.log(key)
    console.log(this.state.removeKey)
    let arr = this.state.files
    let newFiles = arr.splice(key)
    this.setState({files:newFiles})
  }
  onDragEnterHandler(e){
    e.stopPropagation();
    e.preventDefault();
    this.refs.upload_file.style.backgroundColor = "#f9f9f9"
  }
  onMouseOutHandler(e){
    e.stopPropagation();
    e.preventDefault();
    this.refs.upload_file.style.backgroundColor = "transparent"
    console.log("Exit")
  }
  onDragOverHandler(e){
    e.stopPropagation();
    e.preventDefault();
  }
  onDropFileHandler(e){
    console.log(e)
    e.stopPropagation();
    e.preventDefault();

    this.FileSelectHandler(e)
  }

  render() {
    let {state} = this;
    let selectedFiles = state.files.map(function(file,key){
      console.log(key)
      return (<div className="pr1 pb1 w-25 fl" key={key}>
          <div className="h-100 ba b--black-20 tc pa2">
            <div>
              <a className="dib  link pointer navy" onClick={()=>{this.removeFile(key)}}>x</a>
            </div>
            <img src={GetFileRepresentativeImage(file)} className="w3 h3 dib v-mid" alt="file representative logo"/>
            <small className="db link pointer truncate pb1 navy pa1" >{file.name}</small>
          </div>
      </div>)
    })
    return (
      <section className="">
        <Nav/>
        <section className="tc pt5">
          <section className="pt5 dib w-100 w-70-m w-50-l ">
            <div className="pv3">
              <span className="navy w-100 f3">New Submission</span>
            </div>
            <form className="pv3" onSubmit={this.submitForm.bind(this)}>
              <div className="pv3">
                <input type="text" className="pa3 w-100" placeholder="Submission Name" ref="submissionName"/>
              </div>
              <div className="pv3">
                <label htmlFor="upload_file" ref="upload_file" className="pa5 w-100 border-box dib ba b--black-30 navy hover-bg-light-gray pointer" onDrop={this.onDropFileHandler.bind(this)} onDragEnter={this.onDragEnterHandler.bind(this)} onDragOver={this.onDragOverHandler.bind(this)} onDragEnd={this.onMouseOutHandler.bind(this)}>
                  <div>Drop File to Upload</div>
                  <div><small>( or click to select a file )</small></div>
                </label>
                <input type="file" className="dn" id="upload_file" onChange={this.FileSelectHandler.bind(this)}/>
              </div>
              <div className="cf">
                {selectedFiles}
              </div>
              <div className="pv3">
                {state.showSuccessMessage?<p className="pa3 ba">Submitted Successfully</p>:""}
                {state.showErrorMessage?<p className="pa3 ba">Error In Submission</p>:""}
              </div>
              <div className="pv3 tr">
                <button className="pa3 bg-navy grow shadow-4  bw0 white-80 hover-white" type="submit">Submit</button>
              </div>
            </form>
          </section>
        </section>
      </section>
    );
  }
}

export default NewSubmissionPage;
