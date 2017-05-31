import React, { Component } from 'react';
import ReactDOM from 'react-dom';

export default class Modal extends Component {

  render(){
    let {openModal,closeModal} = this.props;
    return (
      <section className={"vh-100 fixed w-100  justify-center items-center z-4 top-0 left-0 animated "+(openModal?"flex fadeIn":"dn fadeOut ")} style={{backgroundColor:"rgba(0,0,0,0.4)"}}>
        <div className="bg-white w-100 w-60-ns modal-shadow giorgia f6 " >
          <div className=" bg-light-gray pv2 ph3 shadow-btm ">
            <div className="pv1 cf">
              <strong className="dib v-mid fw4 pv2 ph3">Settings</strong>
              <button className="fr dib v-mid pv2 ph3 bg-navy white shadow-4 bw0 grow pointer" onClick={closeModal}>close</button>
            </div>
          </div>
          <div>
            <div className="w-20 dib fl br1 ph2 pv3">
              <a className="dib pv2 ph3 hover-bg-light-gray w-100">permissions</a>

            </div>
            <div className="w-80 dib fl ph2 pv3">
              <h4 className="ma0 pv2 pb3 fw4">who has access</h4>
              <div>

                <div className="pv3 ph2 bt bb bw-tiny b--light-gray">
                  <div className="dib">
                    <span className="db">
                      Anthony Alaribe
                    </span>
                    <small className="db">
                      anthonyalaribe@gmail.com
                    </small>
                  </div>
                  <div className="dib fr ">
                    <select>
                      <option>Spectator (can view)</option>
                      <option>can edit</option>
                    </select>
                  </div>
                </div>

                <div className="pv3 ph2 bt bb bw-tiny b--light-gray">
                  <div className="dib">
                    <span className="db">
                      Anthony Alaribe
                    </span>
                    <small className="db">
                      anthonyalaribe@gmail.com
                    </small>
                  </div>
                  <div className="dib fr ">
                    <select>
                      <option>Spectator (can view)</option>
                      <option>can edit</option>
                    </select>
                  </div>
                </div>

                <div className="pv3 ph2 bt bb bw-tiny b--light-gray">
                  <div className="dib">
                    <span className="db">
                      Anthony Alaribe
                    </span>
                    <small className="db">
                      anthonyalaribe@gmail.com
                    </small>
                  </div>
                  <div className="dib fr ">
                    <select>
                      <option>Spectator (can view)</option>
                      <option>can edit</option>
                    </select>
                  </div>
                </div>

                <div className="pv3 ph2 bt bb bw-tiny b--light-gray">
                  <div className="dib">
                    <span className="db">
                      Anthony Alaribe
                    </span>
                    <small className="db">
                      anthonyalaribe@gmail.com
                    </small>
                  </div>
                  <div className="dib fr ">
                    <select>
                      <option>Spectator (can view)</option>
                      <option>can edit</option>
                    </select>
                  </div>
                </div>

                <div className="pv3 ph2 mv3  bt bb bw-tiny b--light-gray bg-near-white">
                  <div className="dib">
                    <input type="" className="pv2 ph3 w-100" />
                  </div>
                  <div className="dib fr ">
                    <select>
                      <option>Spectator (can view)</option>
                      <option>can edit</option>
                    </select>
                  </div>
                </div>

              </div>
            </div>
          </div>
        </div>
      </section>
    )
  }
}
