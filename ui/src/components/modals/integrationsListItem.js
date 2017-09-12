import React, { Component } from "react";
import { inject, observer } from "mobx-react";
import { withRouter } from "react-router";
import IntegrationForm from "./integrationForm.js";


@inject("MainStore", "IntegrationsStore")
@observer
class integrationsListItem extends Component {
  state ={
    Edit: false,
  }
  render(){
    let {integration,IntegrationsStore  } = this.props;
    let { workspaceID, formID } = this.props.match.params;
    console.log(integration)
    return (
      <div className="pa2 mv2 ba b--light-gray  " >
        <div className="db cf">
          <strong className="f5 fw5 db ">
            {integration.URL}
          </strong>
        </div>
        <div className="cf pv2">
          <a
            className="ba b--light-gray navy bg-transparent pv1 ph2 link pointer "
            onClick={()=>{
              this.setState({"Edit":!this.state.Edit})
            }}
          >
            Edit
          </a>
          <div className="di">
            <button className="pv1 ph2 ba b--light-gray navy bg-transparent pv1 ph2 link  pointer"
              onClick={()=>IntegrationsStore.testFormIntegration(workspaceID, formID, integration)}>
              Test
            </button>
          </div>
          <a
            data-confirm="Are you sure?"
            className=" link bg-transparent b--light-gray navy pv1 ph2 ba pointer"
            rel="nofollow"
            onClick={()=>IntegrationsStore.DeleteFormIntegration(workspaceID, formID, integration)}
          >
            <span className="">Remove</span>
          </a>
        </div>
        <IntegrationForm integration={integration} show={this.state.Edit} onSave={()=>this.setState({"Edit":!this.state.Edit})}/>
      </div>
    )
  }
}

export default withRouter(integrationsListItem);
