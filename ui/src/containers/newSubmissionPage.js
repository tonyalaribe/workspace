import React, { Component } from 'react';
import Nav from '../components/nav.js';
import {inject, observer} from "mobx-react";

var fileImageRepresentation = require("../assets/files.png")
function GetFileRepresentativeImage(file){

  switch (file.type){
    case "image/png":
      return file.file
    case "image/jpeg":
      return file.file
    case "image/jpg":
      return file.file
    case "image/gif":
      return file.file
    default:
      return fileImageRepresentation
  }
}

@inject("MainStore") @observer
class NewSubmissionPage extends Component {
  state = {files:[]}
  constructor(props){
    super(props)
    this.removeFile = this.removeFile.bind(this)
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
    let selectedFiles = this.state.files.map(function(file,key){
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
              <span className="navy w-100 f3">New Submission{this.props.MainStore.xyz}</span>
            </div>
            <section className="pv3">
              <div className="pv3">
                <input type="text" className="pa3 w-100" placeholder="Submission Name"/>
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
              <div className="pv3 tr">
                <button className="pa3 bg-navy grow shadow-4  bw0 white-80 hover-white">Submit</button>
              </div>
            </section>
          </section>
        </section>
      </section>
    );
  }
}

export default NewSubmissionPage;
