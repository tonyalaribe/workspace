import React, { Component } from 'react';
import {GetFileRepresentativeImage } from '../utils/representativeImages.js';


class FileSelect extends Component {
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
           let newFile = {
             file:e.target.result,
             name:file.name,
             type:file.type
           }

           this.props.newFile(newFile)
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

  render(){
    let {props} = this;
    let selectedFiles = props.files.map(function(file,key){
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
      <section>
        <div className="pv3">
          <label htmlFor="upload_file" ref="upload_file" className="pa5 w-100 border-box dib ba b--black-30 navy hover-bg-light-gray pointer" onDrop={this.onDropFileHandler.bind(this)} onDragEnter={this.onDragEnterHandler.bind(this)} onDragOver={this.onDragOverHandler.bind(this)} onDragEnd={this.onMouseOutHandler.bind(this)}>
            <div>Drop File to Upload</div>
            <div><small>( or click to select a file )</small></div>
          </label>
          <input type="file" className="dn" id="upload_file" onChange={this.FileSelectHandler.bind(this)} multiple/>
        </div>
        <div className="cf">
          {selectedFiles}
        </div>
      </section>
    )
  }
}

export default FileSelect;
