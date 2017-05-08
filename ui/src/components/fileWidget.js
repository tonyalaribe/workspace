import React, {Component} from 'react';
import {GetFileRepresentativeImage} from '../utils/representativeImages.js';
//import {  setState } from "../../utils";

function addNameToDataURL(dataURL, name) {
  return dataURL.replace(';base64', `;name=${name};base64`);
}

export function setState(instance, state, callback) {
  const {safeRenderCompletion} = instance.props;
  if (safeRenderCompletion) {
    instance.setState(state, callback);
  } else {
    instance.setState(state);
    setImmediate(callback);
  }
}

function processFile(file) {
  const {name, size, type} = file;
  return new Promise((resolve, reject) => {
    const reader = new window.FileReader();
    reader.onload = event => {
      resolve({
        dataURL: addNameToDataURL(event.target.result, name),
        name,
        size,
        type,
      });
    };
    reader.readAsDataURL(file);
  });
}

function processFiles(files) {
  return Promise.all([].map.call(files, processFile));
}

function FilesInfo(props) {
  const {filesInfo} = props;
  console.log(props)
  console.log(filesInfo)
  if (filesInfo.length === 0) {
    return null;
  }
  return (
    <div className="file-info cf">
      {filesInfo.map((fileInfo, key) => {
        const {name} = fileInfo;
        return (
          <div className="pr1 pb1 w-25 fl" key={key}>
            <div className="h-100 ba b--black-20 tc pa2">
              <div>
                <a className="dib  link pointer navy">x</a>
              </div>
              <img
                src={GetFileRepresentativeImage(name)}
                className="w3 h3 dib v-mid"
                alt="file representative logo"
              />
              <small className="db link pointer truncate pb1 navy pa1">
                {name}
              </small>
            </div>
          </div>
        );
      })}
    </div>
  );
}

export function dataURItoBlob(dataURI) {
  if (dataURI.split('base64,').length < 2) {
    return {
      blob: dataURI,
      name: dataURI,
    };
    // console.log("not a base64 file")
    // console.log({
    //   blob:dataURI,
    //   name:dataURI
    // })
  }
  // Split metadata from data
  const splitted = dataURI.split(',');
  // Split params
  const params = splitted[0].split(';');
  // Get mime-type from params
  const type = params[0].replace('data:', '');
  // Filter the name property from params
  const properties = params.filter(param => {
    return param.split('=')[0] === 'name';
  });
  // Look for the name and use unknown if no name property.
  let name;
  if (properties.length !== 1) {
    name = 'unknown';
  } else {
    // Because we filtered out the other property,
    // we only have the name case here.
    name = properties[0].split('=')[1];
  }

  // Built the Uint8Array Blob parameter from the base64 string.
  const binary = atob(splitted[1]);
  const array = [];
  for (let i = 0; i < binary.length; i++) {
    array.push(binary.charCodeAt(i));
  }
  // Create the blob object
  const blob = new window.Blob([new Uint8Array(array)], {type});

  return {blob, name};
}

function extractFileInfo(dataURLs) {
  // console.log(dataURLs)
  return dataURLs
    .filter(dataURL => typeof dataURL !== 'undefined')
    .map(dataURL => {
      const {blob, name} = dataURItoBlob(dataURL);

      console.log({
        name: name,
        size: blob.size,
        type: blob.type,
      });
      return {
        name: name,
        size: blob.size,
        type: blob.type,
      };
    });
}

class FileWidget extends Component {
  defaultProps = {
    multiple: false,
  };

  constructor(props) {
    super(props);
    const {value} = props;
    console.log("constructor props", props)
    console.log("constructor value", value)
    const values = Array.isArray(value) ? value : [value];
    this.state = {values, filesInfo: extractFileInfo(values)};
  }

  // shouldComponentUpdate(nextProps, nextState) {
  //   return shouldRender(this, nextProps, nextState);
  // }

  onChange = event => {
    const {multiple, onChange} = this.props;
    processFiles(event.target.files).then(filesInfo => {
      const state = {
        values: filesInfo.map(fileInfo => fileInfo.dataURL),
        filesInfo,
      };
      setState(this, state, () => {
        if (multiple) {
          onChange(state.values);
        } else {
          onChange(state.values[0]);
        }
      });
    });
  };

  render() {
    const {multiple, id, readonly, disabled, autofocus} = this.props;
    const {filesInfo} = this.state;
    return (
      <div className="">
        <div className="">
          <label
            htmlFor={id}
            ref="upload_file"
            className="pa5 w-100 border-box dib ba b--black-30 navy hover-bg-light-gray pointer"
          >
            <div>Drop File to Upload</div>
            <div><small>( or click to select a file )</small></div>
          </label>
          <input
            ref={ref => this.inputRef = ref}
            id={id}
            type="file"
            disabled={readonly || disabled}
            onChange={this.onChange}
            defaultValue=""
            autoFocus={autofocus}
            multiple={multiple}
            className="dn"
          />
      </div>
        <FilesInfo filesInfo={filesInfo} />
      </div>
    );
  }
}

export default FileWidget;
