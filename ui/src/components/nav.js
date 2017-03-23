import React, { Component } from 'react';
import {Link} from 'react-router-dom';


class Nav extends Component {
  render() {
    return (
      <nav className="bg-navy w-100 fixed shadow-4 pa3 ph4 white-80 dib">
        <Link to="/" className="pa2 dib link white-80 hover-white">
          Submitter
        </Link>
        <div className="dib pa2 fr pr5">
          <Link to="/new_submission" className="white-70  link pa2 ba hover-white">
            New Submission
          </Link>
        </div>
      </nav>
    );
  }
}

export default Nav;
