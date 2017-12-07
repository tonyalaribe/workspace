import Auth0Lock from 'auth0-lock';
import mitt from 'mitt';
import jwtDecode from 'jwt-decode';
//import { browserHistory } from 'react-router'

 class AuthService {
  constructor(clientId, domain) {
    // Configure Auth0
    this.lock = new Auth0Lock(clientId, domain, {
      auth: {
        redirectUrl: window.location.origin+'/login',
        responseType: 'token'
      }
    })
    // Add callback for lock `authenticated` event
    this.lock.on('authenticated', this._doAuthentication.bind(this))
    // binds login functions to keep this context
    this.login = this.login.bind(this)

    this.emitter = mitt()
  }

  _doAuthentication(authResult) {
    // Saves the user token
    this.setToken(authResult.idToken)
    console.log("supposed to redirect")
    // navigate to the home route
    //browserHistory.replace('/')
    this.lock.getProfile(authResult.idToken, (error, profile) => {
      if (error) {
        console.log('Error loading the Profile', error)
      } else {
        this.setProfile(profile)
      }
    })
    this.emitter.emit('authenticated',true)
  }
  setProfile(profile) {
    // Saves profile data to local storage
    localStorage.setItem('profile', JSON.stringify(profile))
    // Triggers profile_updated event to update the UI
    this.emitter.emit('profile_updated',profile)

  }

  getProfile() {
    // Retrieves the profile data from local storage
    const profile = localStorage.getItem('profile')

    return profile ? JSON.parse(profile) : {}
  }

  login() {
    // Call the show method to display the widget.
    this.lock.show()
  }

  loggedIn() {
    // Checks if there is a saved token and it's still valid
    let token = this.getToken()
    if (token && token.length>0){
      console.log(token)
      const authMeta = jwtDecode(token);
      console.log(authMeta)
      const current_time = new Date().getTime() / 1000;
      if (current_time > authMeta.exp) {
        /* expired */
        this.logout();
        return false;
      }
      return true
    }
    this.logout();
    return false
  }

  setToken(idToken) {
    // Saves user token to local storage
    localStorage.setItem('id_token', idToken)

  }

  getToken() {
    // Retrieves the user token from local storage
    return localStorage.getItem('id_token')
  }

  logout() {
    // Clear user token and profile data from local storage
    localStorage.removeItem('id_token');
    localStorage.removeItem('profile');
  }

}

const authService = new AuthService(process.env.AUTH0_SECRET, process.env.AUTH0_ENDPOINT)

export default authService
