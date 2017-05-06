import React, {Component} from 'react';
import Nav from '../components/nav.js';
import {inject, observer} from 'mobx-react';
import moment from 'moment';

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
      console.log(jsonschema.properties[current]);
      console.log(SubmissionInfo)
      switch (jsonschema.properties[current].type) {
        case 'string':
          console.log('type string');
          console.log(jsonschema.properties[current].format);
          switch (jsonschema.properties[current].format) {
            case 'data-url':
              value = (
                <a
                  target="_blank"
                  href={'/' + SubmissionInfo.formData[current]}
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
          value = SubmissionInfo.formData[current].map(function(item) {
            console.log(item);
            console.log(jsonschema.properties[current]);
            switch (jsonschema.properties[current].items.type) {
              case 'string':
                switch (jsonschema.properties[current].items.format) {
                  case 'data-url':
                    return <a target="_blank" href={'/' + item}>xxx</a>;
                  default:
                    return item;
                }
              default:
                return item;
            }
          });
        default:
          console.log('type unknown');
          console.log(jsonschema.properties[current]);
          value = SubmissionInfo.formData[current];
          break;
      }
      previous.push(
        <div className="pv2" key={current}>
          <strong className="pa1 dib">
            {jsonschema.properties[current].title}: &nbsp;&nbsp;
          </strong>
          <span className="pa1 dib">{value}</span>
        </div>,
      );
      return previous;
    }, []);

    return (
      <section className="">
        <Nav />
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
