import React, { Component } from 'react';
import Nav from '../components/nav.js';

class UploadsPage extends Component {

  render() {
    let userFilesCard = [{},{},{}].map(function(fileData, key){
      return (<div className="shadow-4 mv2 h4" key={key}>
        {/** Upload Item **/}
        <div className="dib w-30 v-top tc h-100 fl">
          <div className="h-100 flex flex-column  items-center justify-around">
            <img src="//placehold.it/200x200" className="w3 h3 dib v-mid" alt="file representative logo"/>
          </div>
        </div><div className="dib w-70 h-100 v-top bl b--light-gray pa3">
          <h3 className="navy mv1 ">File name </h3>
          <div><small>Uploaded By: bla bla</small></div>
        </div>
        {/** End Upload Item **/}
      </div>)

    })
    return (
      <section className="">
        <Nav/>
        <section className="tc pt5">
          <section className="pt4 dib w-100 w-70-m w-50-l tl">
            <div className="pv3">
              <span className="navy w-100">Today</span>
            </div>
            <section>
              {userFilesCard}
            </section>
          </section>
        </section>
      </section>
    );
  }
}

export default UploadsPage;
