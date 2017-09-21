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
          <div className="bb b--light-gray grow ph3 pv2">

            <h3 className="navy mv1 fw4">{workspace.name}</h3>
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
      <section >
        <Nav />
        <section className="tc pt5 vh-100 flex flex-column items-center justify-center">
          <section className="pt4 dib w-100 w-60-m w-40-l tl ">
            <div className="pv3 cf">
              <strong className="dib navy v-btm f4">Select a  Workspace</strong>
              <div className="fr  dib">
                <Link
                  to="/new_workspace"
                  className="ph3 pv2 ba link navy dib grow"
                >
                  New Workspace
                </Link>

              </div>
            </div>
            <section className="shadow-4  pa3 bt bw2 b--custom-green">
              {AllWorkspaces}
            </section>
          </section>
        </section>
      </section>
    );
  }
}

export default ListOfWorkspaces;
