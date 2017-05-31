import React, {Component} from 'react';
import Nav from '../components/nav.js';
import {inject, observer} from 'mobx-react';
import {Link} from 'react-router-dom';
import Modal from '../components/modal.js';


@inject('MainStore')
@observer
class ListOfWorkspaces extends Component {
  constructor () {
    super();
    this.state = {
      showModal: false
    };

    this.handleOpenModal = this.handleOpenModal.bind(this);
    this.handleCloseModal = this.handleCloseModal.bind(this);
  }

  componentDidMount() {
    this.props.MainStore.getAllWorkspaces();
  }
  handleOpenModal () {
    this.setState({ showModal: true });
  }

  handleCloseModal () {
    this.setState({ showModal: false });
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
      <section >
        <Nav />
        <section className="tc pt5">
          <section className="pt4 dib w-100 w-70-m w-50-l tl">
            <div className="pv3 cf">
              <span className="dib navy v-btm">All Available Workspaces</span>
              <div className="fr  dib">
                <Link
                  to="/new_workspace"
                  className="ph3 pv2 ba link navy dib grow"
                >
                  New Workspace
                </Link>
                <a href="#" className="dib link pa2 navy" onClick={this.handleOpenModal}>⚙ &nbsp;settings</a>
                <Modal openModal={this.state.showModal} closeModal={this.handleCloseModal}/>
              </div>
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
