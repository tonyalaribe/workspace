import React, { Component } from 'react';
import {inject, observer} from 'mobx-react';
import { withRouter } from 'react-router'
import iziToast from "izitoast";

@inject('MainStore','PermissionsStore')
@observer
class modal extends Component {
  constructor(props){
    super(props)
    this.AddCollaborator = this.AddCollaborator.bind(this)
  }
  componentDidMount(){
    this.props.PermissionsStore.getWorkspaceUsersAndRoles(this.props.match.params.workspaceID)
  }
  AddCollaborator(){
    let permissions = {}
    permissions.email = this.refs.email.value
    permissions.role = this.refs.role.value

    let {PermissionsStore,match, closeModal} = this.props

    PermissionsStore.updateUserWorkspacePermissions(match.params.workspaceID, permissions, ()=>{
      iziToast.success({
          title: 'Modify Collaborators',
          message: `Collaborator was added/modified successfully`,
          position: 'topRight',
      });
      closeModal()
    })
  }
  render(){
    let {openModal,closeModal} = this.props;
    let {WorkspaceUsers} = this.props.PermissionsStore;
    return (
      <section className={"vh-100 fixed w-100  justify-center items-center z-4 top-0 left-0 animated "+(openModal?"flex fadeIn":"dn fadeOut ")} style={{backgroundColor:"rgba(0,0,0,0.4)"}}>
        <div className="bg-white w-100 w-60-ns modal-shadow giorgia f6 " >
          <div className=" bg-light-gray pv2 ph3 shadow-btm ">
            <div className="pv1 cf">
              <strong className="dib v-mid fw4 pv2 ph3">Settings</strong>
              <button className="fr dib v-mid pv2 ph3 bg-navy white shadow-4 bw0 grow pointer" onClick={closeModal}>close</button>
            </div>
          </div>
          <div>
            <div className="w-20 dib fl br1 ph2 pv3">
              <a className="dib pv2 ph3 hover-bg-light-gray w-100">permissions</a>
            </div>
            <div className="w-80 dib fl ph2 pv3">
              <h4 className="ma0 pv2 pb3 fw4 ph2">who has access</h4>
              <div>

                {WorkspaceUsers.map(function(user, key){
                  let roleComponents = user.CurrentRoleString.split("-")
                  let roleLevel = roleComponents[roleComponents.length -1]
                  return (<div className="pv3 ph2 bt bb bw-tiny b--light-gray" key={key}>
                    <div className="dib">
                      <span className="db">
                        {user.Username}
                      </span>
                      <small className="db">
                        {user.Email}
                      </small>
                    </div>
                    <div className="dib fr ">
                      {
                        user.CurrentRoleString==="superadmin"?
                        <div>
                          <span>superadmin</span>
                        </div>
                        :
                        <select defaultValue={roleLevel}>
                          <option value="spectator"  >Spectator (can view)</option>
                          <option value="editor">Editor (can edit)</option>
                          <option value="supervisor">Supervisor (can approve)</option>
                          <option value="admin" >Admin (admin)</option>
                        </select>
                      }
                    </div>
                  </div>)
                })

              }

                <div className="pv3 ph2 mv3  bt bb bw-tiny b--light-gray bg-near-white cf">
                  <div className="dib">
                    <input type="email" className="pv2 ph3 w-100" placeholder="email" ref="email"/>
                  </div>
                  <div className="dib fr ">
                    <select ref="role">
                      <option value="spectator">Spectator (can view)</option>
                      <option value="editor">Editor (can edit)</option>
                      <option value="supervisor">Supervisor (can approve)</option>
                      <option value="admin">Admin (admin)</option>
                    </select>
                  </div>
                  <div className="pv2 tr">
                    <button className="bg-navy bw0 shadow-4 pv2 ph3 white-80 grow" onClick={this.AddCollaborator}>add collaborator</button>
                  </div>
                </div>

              </div>
            </div>
          </div>
        </div>
      </section>
    )
  }
}

var Modal = withRouter(modal)
export default Modal;
