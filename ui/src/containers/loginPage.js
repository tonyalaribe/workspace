import React, { Component } from 'react';
import {Redirect} from 'react-router-dom';
import AuthService from '../utils/auth0.js';

class LoginPage extends Component {
  constructor(props){
    super(props)
    AuthService.emitter.on("authenticated",()=>{
      console.log("authenticated")
      this.setState({authenticated:true})
      AuthService.lock.hide()
    })
  }
  state = {}
  componentDidMount(){
    AuthService.login()
  }
  render() {
    return this.state.authenticated?(<Redirect to="/"/>):(
      <section className="">
      </section>
    );
  }
}

export default LoginPage;
