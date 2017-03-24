import React, {Component} from 'react';
import {
  Route,
  Redirect
} from 'react-router-dom'
import AuthService from '../utils/auth0.js';

class ProtectedRoute extends Component {
  render(){
    let {component, ...rest} = this.props;
    let loggedIn = AuthService.loggedIn()
    console.log(loggedIn)
    return (
      <Route {...rest} render={props => (
         loggedIn? (
          React.createElement(component, props)
        ) : (
          <Redirect to={{
            pathname: '/login',
            state: { from: props.location }
          }}/>
          )
      )}/>
    )
  }
}

export default ProtectedRoute;
