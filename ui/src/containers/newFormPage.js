import React, {Component} from 'react';
import Nav from '../components/nav.js';
// import FileSelect from '../components/fileSelect.js';
import {inject, observer} from 'mobx-react';

@inject('MainStore')
@observer
class NewFormPage extends Component {
  state = {};

  submitFormToServer() {
    let workspaceID = this.props.match.params.workspaceID
    this.setState({showSuccessMessage: false});
    let form = {};
    form.name = this.refs.formName.value;
    form.jsonschema = JSON.parse(this.refs.jsonSchema.value);
    form.uischema = JSON.parse(this.refs.uiSchema.value);
    console.log(form);

    this.props.MainStore.submitNewFormToServer(workspaceID, form, () => {
      this.setState({showSuccessMessage: true});
      this.refs.formName.value = '';
      this.refs.jsonSchema.value = '';
      this.refs.uiSchema.value = '';
    });
  }
  render() {
    let {state} = this;

    return (
      <section className="">
        <Nav />
        <section className="tc pt5">
          <section className="pt5 dib w-100 w-70-m w-50-l ">
            <div className="pv3 ">
              <span className="navy w-100 f3 db">
                New Form
              </span>
            </div>
            <div className="pv3 tl">
              <label className="pv2 dib">
                Form Name
              </label>
              <input
                type="text"
                className="form-control "
                ref="formName"
              />
            </div>
            <div className="pv3 tl z-1">
              <label className="pv2 dib">
                JSON Schema
              </label>
              <textarea
                type="text"
                className="form-control w-100 h5  z-1"
                ref="jsonSchema"
              />
            </div>
            <div className="pv3 tl z-1">
              <label className="pv2 dib">
                UI Schema
              </label>
              <textarea
                type="text"
                className="form-control w-100 h5 z-1"
                ref="uiSchema"
              />
            </div>

            <div className="pv3">
              {state.showSuccessMessage
                ? <p className="pa3 ba">
                    Submitted Successfully
                  </p>
                : ''}
              {state.showErrorMessage
                ? <p className="pa3 ba">
                    Error In Submission
                  </p>
                : ''}
            </div>

            <div className="pv3 tr">
              <button
                className="pa3 bg-navy grow shadow-4  bw0 white-80 hover-white ml2 pointer"
                onClick={this.submitFormToServer.bind(this)}
              >
                publish
              </button>
            </div>
          </section>
        </section>
      </section>
    );
  }
}

export default NewFormPage;
