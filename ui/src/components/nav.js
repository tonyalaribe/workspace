import React, { Component } from 'react';


class Nav extends Component {
  render() {
    return (
      <nav className="bg-navy w-100 fixed shadow-4 pa3 ph4 white-80 dib">
        <span className="pa2 dib">
          Submitter
        </span>
        <div className="dib pa2 fr pr5">
          <a href="#" className="  white-70  link pa2 ba hover-white">
            New Submission
          </a>
        </div>
      </nav>
    );
  }
}

export default Nav;
