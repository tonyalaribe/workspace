import {observable, action, runInAction} from 'mobx';
import AuthService from '../utils/auth0.js';

class mainStore {
  @observable Submissions = [];
  @observable AllWorkspaces = [];
  @observable AllForms = [];

  @observable SubmissionInfo = {formData: []};
  @observable CurrentForm= {
    jsonschema: {
      properties: {},
    },
    uischema: {},
  };

  @action submitFormToServer = async (formData, callback) => {
    console.log(formData);

    let authToken = AuthService.getToken();

    const response = await fetch(
      '/api/workspaces/' + this.CurrentWorkspace.id + '/new_submission',
      {
        method: 'POST',
        body: JSON.stringify(formData),
        mode: 'cors',
        headers: {
          'Content-type': 'application/json',
          authorization: 'Bearer ' + authToken,
        },
      },
    );
    const data = await response.json();
    /* required in strict mode to be allowed to update state: */
    runInAction('update state after fetching data', () => {
      console.log(data);
      callback();
    });
  };

  @action updateFormOnServer = async (submissionID, formData, callback) => {
    console.log(formData);

    let authToken = AuthService.getToken();

    const response = await fetch(
      '/api/workspaces/' +
        this.CurrentWorkspace.id +
        '/submissions/' +
        submissionID,
      {
        method: 'PUT',
        body: JSON.stringify(formData),
        mode: 'cors',
        headers: {
          'Content-type': 'application/json',
          authorization: 'Bearer ' + authToken,
        },
      },
    );
    const data = await response.json();
    /* required in strict mode to be allowed to update state: */
    runInAction('update state after fetching data', () => {
      console.log(data);
      callback();
    });
  };
  @action submitNewFormToServer = async (workspaceID, form, callback) => {
    let authToken = AuthService.getToken();

    const response = await fetch('/api/workspaces/'+workspaceID+'/new_form', {
      method: 'POST',
      body: JSON.stringify(form),
      mode: 'cors',
      headers: {
        'Content-type': 'application/json',
        authorization: 'Bearer ' + authToken,
      },
    });
    const data = await response.json();
    /* required in strict mode to be allowed to update state: */
    runInAction('update state after fetching data', () => {
      console.log(data);
      callback();
    });
  };


  @action getAllForms = async (workspaceID) => {
    let authToken = AuthService.getToken();
    const response = await fetch('/api/workspaces/'+workspaceID+'/forms', {
      method: 'GET',
      mode: 'cors',
      headers: {
        'Content-type': 'application/json',
        authorization: 'Bearer ' + authToken,
      },
    });
    const data = await response.json();
    /* required in strict mode to be allowed to update state: */
    runInAction('update state after fetching data', () => {
      console.log(data);
      this.AllForms.replace(data);
    });
  };

  @action getWorkspace = async workspaceID => {
    this.CurrentWorkspace.id = workspaceID;

    let authToken = AuthService.getToken();

    const response = await fetch('/api/workspaces/' + workspaceID, {
      method: 'GET',
      mode: 'cors',
      headers: {
        'Content-type': 'application/json',
        authorization: 'Bearer ' + authToken,
      },
    });
    const data = await response.json();
    /* required in strict mode to be allowed to update state: */
    runInAction('update state after fetching data', () => {
      console.log(data);
      this.CurrentWorkspace = data;
    });
  };
  @action getAllWorkspaces = async () => {
    let authToken = AuthService.getToken();
    const response = await fetch('/api/workspaces', {
      method: 'GET',
      mode: 'cors',
      headers: {
        'Content-type': 'application/json',
        authorization: 'Bearer ' + authToken,
      },
    });
    const data = await response.json();
    /* required in strict mode to be allowed to update state: */
    runInAction('update state after fetching data', () => {
      console.log(data);
      this.AllWorkspaces.replace(data);
    });
  };


  @action submitNewWorkspaceToServer = async (workspace, callback) => {
    let authToken = AuthService.getToken();

    const response = await fetch('/api/new_workspace', {
      method: 'POST',
      body: JSON.stringify(workspace),
      mode: 'cors',
      headers: {
        'Content-type': 'application/json',
        authorization: 'Bearer ' + authToken,
      },
    });
    const data = await response.json();
    /* required in strict mode to be allowed to update state: */
    runInAction('update state after fetching data', () => {
      console.log(data);
      callback();
    });
  };

  @action getMySubmissions = async () => {
    let authToken = AuthService.getToken();
    const response = await fetch(
      '/api/workspaces/' + this.CurrentWorkspace.id + '/submissions',
      {
        method: 'GET',
        mode: 'cors',
        headers: {
          'Content-type': 'application/json',
          authorization: 'Bearer ' + authToken,
        },
      },
    );
    const data = await response.json();
    /* required in strict mode to be allowed to update state: */
    runInAction('update state after fetching data', () => {
      console.log(data);
      this.Submissions.replace(data);
    });
  };

  @action getSubmissionInfo = async submissionID => {
    let authToken = AuthService.getToken();

    const response = await fetch(
      '/api/workspaces/' +
        this.CurrentWorkspace.id +
        '/submissions/' +
        submissionID,
      {
        method: 'GET',
        mode: 'cors',
        headers: {
          'Content-type': 'application/json',
          authorization: 'Bearer ' + authToken,
        },
      },
    );
    const data = await response.json();
    /* required in strict mode to be allowed to update state: */
    runInAction('update state after fetching data', () => {
      console.log(data);
      this.SubmissionInfo = data;
    });
  };
}

export const MainStore = new mainStore();
