import React, {Component} from 'react';
import Nav from '../components/nav.js';
// import FileSelect from '../components/fileSelect.js';
import {inject, observer} from 'mobx-react';

@inject('MainStore')
@observer
class NewWorkspacePage extends Component {
  state = {};

  submitWorkspaceFormToServer() {
    this.setState({showSuccessMessage: false});
    let workspace = {};
    workspace.name = this.refs.workspaceName.value;

    this.props.MainStore.submitNewWorkspaceToServer(workspace, () => {
      this.setState({showSuccessMessage: true});
      this.refs.workspaceName.value = '';
    });
  }
  render() {
    let {state} = this;

    return (
      <section>
        <Nav />
        <section className="tc pt5">
          <section className="pt5 dib w-100 w-70-m w-50-l ">
            <div className="pv3 ">
              <span className="navy w-100 f3 db">
                New Workspace
              </span>
            </div>
            <div className="pv3 tl">
              <label className="pv2 dib">
                Workspace Name
              </label>
              <input
                type="text"
                className="form-control "
                ref="workspaceName"
              />
            </div>

            <div className="pv3">
              {state.showSuccessMessage
                ? <p className="pa3 ba">
                    Workspace Created Successfully
                  </p>
                : ''}
              {state.showErrorMessage
                ? <p className="pa3 ba">
                    Error In Workspace Creation
                  </p>
                : ''}
            </div>

            <div className="pv3 tr">
              <button
                className="pa3 bg-navy grow shadow-4  bw0 white-80 hover-white ml2 pointer"
                onClick={this.submitWorkspaceFormToServer.bind(this)}
              >
                create workspace
              </button>
            </div>
          </section>
        </section>
      </section>
    );
  }
}

export default NewWorkspacePage;
