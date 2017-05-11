import React, {Component} from 'react';
import {Link} from 'react-router-dom';
import AuthService from '../utils/auth0.js';
import {observer, inject} from 'mobx-react';

@inject('MainStore')
@observer
class Nav extends Component {
  constructor(props, context) {
    super(props, context);
    this.state = {
      profile: AuthService.getProfile(),
    };

    // listen to profile_updated events to update internal state
    AuthService.emitter.on('profile_updated', newProfile => {
      console.log('profile updated');
      this.setState({profile: newProfile});
    });
  }
  render() {
    return (
      <nav className="bg-navy w-100 fixed shadow-4 pa3 ph4 white-80 dib z-3">
        <Link to="/" className="pa2 dib link white-80 hover-white">
          Workspace
        </Link>
        <div className="dib pa2 fr pr5">
          <span className="pa1 pl4 dib">
            {this.state.profile.username}
            :
            <Link
              to="/"
              onClick={AuthService.logout}
              className="white-90 hover-white dib ph1 link"
            >
              logout
            </Link>
          </span>
        </div>
      </nav>
    );
  }
}

export default Nav;
