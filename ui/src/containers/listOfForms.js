import React, {Component} from 'react';
import Nav from '../components/nav.js';
import {inject, observer} from 'mobx-react';
import {Link} from 'react-router-dom';

@inject('MainStore')
@observer
class ListOfForms extends Component {
  componentDidMount() {
    let workspaceID = this.props.match.params.workspaceID
    this.props.MainStore.getAllForms(workspaceID);
  }
  render() {
    let {MainStore} = this.props;

    let workspaceID = this.props.match.params.workspaceID

    let AllForms = MainStore.AllForms.map(function(form, key) {
      return (
        <Link
          to={'/workspaces/'+workspaceID+'/forms/' + form.id}
          key={key}
          className="link navy"
        >
          <div className="shadow-4 grow mv2 pa3">

            <h3 className="navy mv1 ">{form.name}</h3>
            <div>
              <div className=" pv1">
                <small>created by:&nbsp;&nbsp;&nbsp;{form.creator}</small>
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
                to={'/workspaces/'+workspaceID+'/new_form'}
                className="ph3 pv2 ba fr link navy dib grow"
              >
                New Forms
              </Link>
              <span className="navy w-100 v-btm">All Available Forms</span>
            </div>
            <section>
              {AllForms}
            </section>
          </section>
        </section>
      </section>
    );
  }
}

export default ListOfForms;
