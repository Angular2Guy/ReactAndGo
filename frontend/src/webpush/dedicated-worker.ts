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
        if (Notification.permission === 'granted') { 
          if(resultJson?.length > 0 && resultJson[0]?.Message?.length > 1 && resultJson[0]?.Title?.length > 1) {            
            for(let value of resultJson) {
            new Notification(value?.Title, {body: value?.Message});
            }
          }                
        }
      }
    });
  }, 60000);
});

