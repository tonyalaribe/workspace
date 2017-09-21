import React, {Component} from 'react';
import Nav from '../../components/nav.js';
import {inject, observer} from 'mobx-react';
import moment from 'moment';
import iziToast from "izitoast";

@inject('MainStore')
@observer
class PublishedSubmissionInfoPage extends Component {
  state = {};

  componentDidMount() {
    let {workspaceID, formID, submissionID} = this.props.match.params;
    this.props.MainStore.getSubmissionInfo(workspaceID, formID, submissionID).then(()=>{
      this.props.MainStore.getFormInfo(workspaceID, formID)
    })
  }
  render() {
    let {CurrentForm, SubmissionInfo} = this.props.MainStore;
    let jsonschema = CurrentForm.jsonschema;

    let formFields = Object.keys(
      jsonschema.properties,
    ).reduce((previous, current) => {
      let value;
      switch (jsonschema.properties[current].type) {
        case 'string':
          switch (jsonschema.properties[current].format) {
            case 'data-url':
              value = (
                <a
                  target="_blank" className="db link pa3 mv1 shadow-4 navy underline-hover overflow-hidden"
                  href={ SubmissionInfo.formData[current]}
                >
                  {SubmissionInfo.formData[current]}
                </a>
              );
              break;
            default:
              value = SubmissionInfo.formData[current];
              break;
          }
          break;
        case 'array':
          if (SubmissionInfo.formData[current]){
            value = SubmissionInfo.formData[current].map(function(item,i) {
              switch (jsonschema.properties[current].items.type) {
                case 'string':
                  switch (jsonschema.properties[current].items.format) {
                    case 'data-url':
                      return <a target="_blank" className="db link pa3 mv1 shadow-4 navy underline-hover overflow-hidden" href={item} key={i}>{item}</a>;
                    default:
                      return item;
                  }
                default:
                  return item;
              }
            });
          }
          break;
        default:
          value = SubmissionInfo.formData[current];
          break;
      }
      previous.push(
        <div className="pv2" key={current}>
          <strong className="pa1 dib">
            {jsonschema.properties[current].title}: &nbsp;&nbsp;
          </strong>
          <div className="pa1 dib w-100">{value}</div>
        </div>,
      );
      return previous;
    }, []);

    let {workspaceID} = this.props.match.params;
    return (
      <section className="">
        <Nav workspaceID={workspaceID} />
        <section className="tc pt5">
          <section className="pt4 dib w-100 w-70-m w-50-l tl">
            <div className="pv3">
              <h1 className="navy w-100 mv2">
                {SubmissionInfo.submissionName}
              </h1>
            </div>

            <div className="pv2">
              <strong>status: </strong>
              <span className="navy">{SubmissionInfo.status}</span>
            </div>
            <div className="pv2">
              <div className="w-100 w-50-ns dib ">
                <small>
                  Created:
                  {' '}
                  {moment(SubmissionInfo.created).format('h:mma, MM-DD-YYYY')}
                </small>
              </div>
              <div className="w-100 w-50-ns dib ">
                <small>
                  Modified:
                  {' '}
                  {moment(SubmissionInfo.lastModified).format(
                    'h:mma, MM-DD-YYYY',
                  )}
                </small>
              </div>
            </div>

            <section className="mt5">
              <div className="navy tc bb bw1 b--navy ">
                <h4 className="mv3">Form Data</h4>
              </div>
              <div className=" ph2 pv3 ">
                {formFields}
              </div>
            </section>

          </section>
        </section>
      </section>
    );
  }
}

export default PublishedSubmissionInfoPage;
