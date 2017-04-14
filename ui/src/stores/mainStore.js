import {observable, action, runInAction} from 'mobx';
import AuthService from '../utils/auth0.js';

class  mainStore {
  @observable Submissions = []
  @observable SubmissionInfo = {formData:[]}
  @observable CurrentWorkspace = {
    jsonschema:{
      properties:{}
    },
    uischema:{}
  }
  @action getWorkspace = async(workspaceID) =>{
    this.CurrentWorkspace.id = workspaceID

    let authToken = AuthService.getToken()

    const response = await fetch("/api/workspaces/"+workspaceID,{
      method: 'GET',
      mode: 'cors',
      headers: {
        "Content-type":"application/json",
        "authorization":"Bearer "+authToken,
      }
    });
    const data = await response.json()
    /* required in strict mode to be allowed to update state: */
    runInAction("update state after fetching data", () => {
        console.log(data)
        console.log("form submitted Successfully")
        this.CurrentWorkspace = data
    })

  }

  @action submitFormToServer = async(formData,callback)=>{
    console.log(formData)

    let authToken = AuthService.getToken()

    const response = await fetch("/api/workspaces/"+this.CurrentWorkspace.id+"/new_submission",{
      method: 'POST',
      body: JSON.stringify(formData),
      mode: 'cors',
      headers: {
        "Content-type":"application/json",
        "authorization":"Bearer "+authToken,
      }
    });
    const data = await response.json()
    /* required in strict mode to be allowed to update state: */
    runInAction("update state after fetching data", () => {
        console.log(data)
        console.log("form submitted Successfully")
        callback()
    })
  }

  @action updateFormOnServer = async(submissionID,formData,callback)=>{
    console.log(formData)

    let authToken = AuthService.getToken()

    const response = await fetch("/api/workspaces/"+this.CurrentWorkspace.id+"/submissions/"+submissionID,{
      method: 'POST',
      body: JSON.stringify(formData),
      mode: 'cors',
      headers: {
        "Content-type":"application/json",
        "authorization":"Bearer "+authToken,
      }
    });
    const data = await response.json()
    /* required in strict mode to be allowed to update state: */
    runInAction("update state after fetching data", () => {
        console.log(data)
        console.log("form submitted Successfully")
        callback()
    })
  }

  @action getMySubmissions = async()=>{
    let authToken = AuthService.getToken()
    const response = await fetch("/api/workspaces/"+this.CurrentWorkspace.id+"/submissions",{
      method: 'GET',
      mode: 'cors',
      headers: {
        "Content-type":"application/json",
        "authorization":"Bearer "+authToken,
      }
    });
    const data = await response.json()
    /* required in strict mode to be allowed to update state: */
    runInAction("update state after fetching data", () => {
        console.log(data)
        this.Submissions.replace(data)
    })
  }

  @action getSubmissionInfo = async(submissionID)=>{
    let authToken = AuthService.getToken()

    const response = await fetch("/api/workspaces/"+this.CurrentWorkspace.id+"/submissions/"+submissionID,{
      method: 'GET',
      mode: 'cors',
      headers: {
        "Content-type":"application/json",
        "authorization":"Bearer "+authToken,
      }
    });
    const data = await response.json()
    /* required in strict mode to be allowed to update state: */
    runInAction("update state after fetching data", () => {
        console.log(data)
        this.SubmissionInfo = data
    })
  }


}


export const MainStore = new mainStore();
