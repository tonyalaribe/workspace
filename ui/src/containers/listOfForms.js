import React, {Component} from 'react';
import Nav from '../components/nav.js';
import {inject, observer} from 'mobx-react';
import {Link} from 'react-router-dom';

@inject('MainStore')
@observer
class ListOfForms extends Component {
  componentDidMount() {
    let workspaceID = this.props.match.params.workspaceID
    this.props.MainStore.getAllForms(workspaceID);
    this.props.MainStore.getAllWorkspaces();
  }
  componentWillUpdate(nextProps, nextState){
    if (this.props.location.pathname!==nextProps.location.pathname){
      this.props.MainStore.getAllForms(nextProps.match.params.workspaceID);
    }
  }
  render() {
    let {MainStore} = this.props;
    let workspaceID = this.props.match.params.workspaceID

    let AllWorkspaces = MainStore.AllWorkspaces.map(function(workspace, key) {
      let workspaceURL = '/workspaces/' + workspace.id
      return (
        <Link
          to={workspaceURL}
          key={key}
          className="link navy"
        >
          <div className={" grow pa1 "+(window.location.pathname.startsWith(workspaceURL)?"bg-light-gray":"")}>
            <span className="navy  ">{workspace.name}</span>
          </div>
        </Link>
      );
    });

    let AllForms = MainStore.AllForms.map(function(form, key) {
      return (
        <Link
          to={'/workspaces/'+workspaceID+'/forms/' + form.id}
          key={key}
          className="link navy"
        >
          <div className="shadow-4 grow mv2 pa3">

            <h3 className="navy mv1 ">{form.name}</h3>
            <div>
              <div className=" pv1">
                <small>created by:&nbsp;&nbsp;&nbsp;{form.creator}</small>
              </div>
            </div>
          </div>
        </Link>
      );
    });


    return (
      <section className="">
        <Nav workspaceID={workspaceID}/>
        <section className="tc pt5">

          <section className="pt4 dib w-100 w-80-m w-60-l tl">
            <div className="w-30 dib v-top pv3 pr3">
              <h3 className="bb dib pa1">Workspaces</h3>
              {AllWorkspaces}
            </div><div className="w-70 dib v-top">
              <div className="pv3 cf">
                <Link
                  to={'/workspaces/'+workspaceID+'/new_form'}
                  className="ph3 pv2 ba fr link navy dib grow"
                >
                  New Forms
                </Link>
                <span className="navy w-100 v-btm">All Available Forms</span>
              </div>
              <section>
                {AllForms}
              </section>
            </div>
          </section>
        </section>
      </section>
    );
  }
}

export default ListOfForms;
