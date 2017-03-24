import {observable, action, runInAction} from 'mobx';
import AuthService from '../utils/auth0.js';

class  mainStore {
  @observable Submissions = []

  @action submitFormToServer = async(formData,callback)=>{
    console.log(formData)

    let authToken = AuthService.getToken()

    const response = await fetch("/api/new_submission",{
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

    const response = await fetch("/api/submissions",{
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
}


export const MainStore = new mainStore();
