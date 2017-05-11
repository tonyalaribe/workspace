import React, {Component} from 'react';
import Nav from '../components/nav.js';
import {inject, observer} from 'mobx-react';
import {Link} from 'react-router-dom';

@inject('MainStore')
@observer
class ListOfWorkspaces extends Component {
  componentDidMount() {
    this.props.MainStore.getAllWorkspaces();
  }
  render() {
    let {MainStore} = this.props;

    let AllWorkspaces = MainStore.AllWorkspaces.map(function(workspace, key) {
      return (
        <Link
          to={'/workspaces/' + workspace.id}
          key={key}
          className="link navy"
        >
          <div className="shadow-4 grow mv2 pa3">

            <h3 className="navy mv1 ">{workspace.name}</h3>
            <div>
              <div className=" pv1">
                <small>created by:&nbsp;&nbsp;&nbsp;{workspace.creator}</small>
              </div>
            </div>
          </div>
        </Link>
      );
    });

    return (
      <section className="">
        <Nav />
        <section className="tc pt5">
          <section className="pt4 dib w-100 w-70-m w-50-l tl">
            <div className="pv3 cf">
              <Link
                to="/new_workspace"
                className="ph3 pv2 ba fr link navy dib grow"
              >
                New Workspace
              </Link>
              <span className="navy w-100 v-btm">All Available Workspaces</span>
            </div>
            <section>
              {AllWorkspaces}
            </section>
          </section>
        </section>
      </section>
    );
  }
}

export default ListOfWorkspaces;
