/// <reference lib="webworker" />
declare var self: DedicatedWorkerGlobalScope;
export {};

interface MsgData {
  jwtToken: string;
  newNotificationUrl: string;  
}

export interface UserResponse {
  Token?:  string
	Message?: string
}

let jwtToken = '';
let tokenIntervalRef: ReturnType<typeof setInterval>;
const refreshToken = (myToken: string) => {
  if(!!tokenIntervalRef) {
    clearInterval(tokenIntervalRef);
  }
  jwtToken = myToken;
  tokenIntervalRef = setInterval(() => {
  const requestOptions = {
      method: 'GET',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${jwtToken}`},            
  };
  fetch('/appuser/refreshtoken', requestOptions).then(response => response.json() as UserResponse).then(result => {
      if((!result.Message && !!result.Token && result.Token.length > 10)) {
          //console.log('Token refreshed.');
          jwtToken = result.Token;          
      } else {
        jwtToken = '';        
        clearInterval(tokenIntervalRef);
      }
  });        
}, 45000);

}

let notificationIntervalRef: ReturnType<typeof setInterval>;
self.addEventListener('message', (event: MessageEvent) => { 
  const msgData = event.data as MsgData;   
  refreshToken(msgData.jwtToken);

  if(!!notificationIntervalRef) {
    clearInterval(notificationIntervalRef);
  }
  notificationIntervalRef = setInterval(() => {
    if(!jwtToken) {
      clearInterval(notificationIntervalRef);
    }
    const requestOptions = {
      method: 'GET',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${msgData.jwtToken}`},            
  };
    self.fetch(msgData.newNotificationUrl, requestOptions).then(result => result.json()).then(resulJson => {
      if(!!resulJson && resulJson?.length > 0) {
        self.postMessage(resulJson);
        //Push Heading/Message
      }
    });
  }, 60000);
});
