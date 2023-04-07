/*
  - Copyright 2022 Sven Loesekann
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/
/// <reference lib="webworker" />
/* eslint-disable-next-line no-restricted-globals */
declare var self: DedicatedWorkerGlobalScope;
export { };

interface MsgData {
  jwtToken: string;
  newNotificationUrl: string;
}

interface UserResponse {
  Token?: string
  Message?: string
}

let jwtToken = '';
let tokenIntervalRef: ReturnType<typeof setInterval>;
const refreshToken = (myToken: string) => {
  if (!!tokenIntervalRef) {
    clearInterval(tokenIntervalRef);
  }
  jwtToken = myToken;
  if (!!jwtToken && jwtToken.length > 10) {
    tokenIntervalRef = setInterval(() => {
      const requestOptions = {
        method: 'GET',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${jwtToken}` },
      };
      fetch('/appuser/refreshtoken', requestOptions).then(response => response.json() as UserResponse).then(result => {
        if ((!result.Message && !!result.Token && result.Token.length > 10)) {
          //console.log('Token refreshed.');
          jwtToken = result.Token;
          /* eslint-disable-next-line no-restricted-globals */
          self.postMessage(result);
        } else {
          jwtToken = '';
          clearInterval(tokenIntervalRef);
        }
      });
    }, 45000);
  }
}

let firstNotRequest = true;
let notificationIntervalRef: ReturnType<typeof setInterval>;
/* eslint-disable-next-line no-restricted-globals */
self.addEventListener('message', (event: MessageEvent) => {
  const msgData = event.data as MsgData;
  refreshToken(msgData.jwtToken);  
  if (!!notificationIntervalRef) {
    clearInterval(notificationIntervalRef);
  }
  notificationIntervalRef = setInterval(() => {
    if (!jwtToken) {
      clearInterval(notificationIntervalRef);
      firstNotRequest = true;
    }
    const requestOptions = {
      method: 'GET',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${jwtToken}` },
    };
    /* eslint-disable-next-line no-restricted-globals */
    self.fetch(msgData.newNotificationUrl, requestOptions).then(result => result.json()).then(resultJson => {
      if (!!resultJson && resultJson?.length > 0) {
        /* eslint-disable-next-line no-restricted-globals */
        self.postMessage(resultJson);
        //Notification
        //console.log(Notification.permission);
        if (Notification.permission === 'granted' && !firstNotRequest) { 
          if(resultJson?.length > 0 && resultJson[0]?.Message?.length > 1 && resultJson[0]?.Title?.length > 1) {            
            for(let value of resultJson) {
            new Notification(value?.Title, {body: value?.Message});
            }
          }                
        } else if(!!firstNotRequest) {
          firstNotRequest = false;
        }
      }
    });
  }, 60000);
});

