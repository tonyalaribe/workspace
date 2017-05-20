import React, {Component} from 'react';
import Nav from '../components/nav.js';
import {Link} from 'react-router-dom';
import AuthService from '../utils/auth0.js';
import {inject, observer} from 'mobx-react';
import moment from 'moment';
var fileImageRepresentation = require('../assets/files.png');

@inject('MainStore')
@observer
class SubmissionsPage extends Component {
  componentDidMount() {
    // console.log(this.props)
    let workspaceID = this.props.match.params.workspaceID
    let formID = this.props.match.params.formID

    this.props.MainStore.getMySubmissions(workspaceID, formID)
  }
  render() {

    let workspaceID = this.props.match.params.workspaceID
    let formID = this.props.match.params.formID

    let submissions = this.props.MainStore.Submissions;

    let userSubmissions = submissions.map(function(fileData, key) {

      return (
        <Link
          to={
            '/workspaces/' +
              workspaceID  +
              '/forms/'    +
              formID        +
              '/submissions/' +
              fileData.status +
              '/' +
              fileData.id
          }
          key={key}
          className="link navy"
        >
          <div className="shadow-4 grow mv2 h4">
            {/** Upload Item **/}
            <div className="dib w-30 v-top tc h-100 fl">
              <div className="h-100 flex flex-column  items-center justify-around">
                <img
                  src={fileImageRepresentation}
                  className="w3 h3 dib v-mid"
                  alt="file representative logo"
                />
              </div>
            </div><div className="dib w-70 h-100 v-top bl b--light-gray pa3">
              <div>
                <small className="fr pa2 bg-navy white-80">
                  {fileData.status}
                </small>
              </div>
              <h3 className="navy mv1 ">{fileData.submissionName} </h3>
              <div>
                <div className=" pv1">
                  <small>
                    created on:&nbsp;&nbsp;&nbsp;
                    {moment(fileData.created).format('h:mma, MM-DD-YYYY')}
                  </small>
                </div>
                <div className=" pv1">
                  <small>
                    modified on:&nbsp;&nbsp;
                    {moment(fileData.lastModified).format('h:mma, MM-DD-YYYY')}
                  </small>
                </div>
              </div>
            </div>
            {/** End Upload Item **/}
          </div>
        </Link>
      );
    });
    console.log(AuthService.loggedIn());
    return (
      <section className="">
        <Nav />
        <section className="tc pt5">
          <section className="pt4 dib w-100 w-70-m w-50-l tl">
            <div className="pv3 cf">
              <Link
                to={'/workspaces/'+workspaceID+'/forms/'+formID+'/new_submission'}
                className="ph3 pv2 ba fr link navy dib grow"
              >
                New Submission
              </Link>
              <span className="navy w-100 v-btm">All Form Submissions</span>
            </div>

            <section>
              {userSubmissions}
            </section>
          </section>
        </section>
      </section>
    );
  }
}

export default SubmissionsPage;