import {observable, action, runInAction} from 'mobx';

class  mainStore {
  xyz="ewwe"

  @action submitFormToServer = async(formData,callback)=>{
    console.log(formData) 
    const response = await fetch("/api/new_submission",{
      method: 'POST',
      body: JSON.stringify(formData),
      mode: 'cors',
    });
    const data = await response.json()
    /* required in strict mode to be allowed to update state: */
    runInAction("update state after fetching data", () => {
        console.log(data)
        console.log("form submitted Successfully")
        callback()
    })
  }
}


export const MainStore = new mainStore();
