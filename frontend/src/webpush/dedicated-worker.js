let jwtToken = '';
let tokenIntervalRef;
const refreshToken = (myToken) => {
    if (!!tokenIntervalRef) {
        clearInterval(tokenIntervalRef);
    }
    jwtToken = myToken;
    tokenIntervalRef = setInterval(() => {
        const requestOptions = {
            method: 'GET',
            headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${jwtToken}` },
        };
        fetch('/appuser/refreshtoken', requestOptions).then(response => response.json()).then(result => {
            if ((!result.Message && !!result.Token && result.Token.length > 10)) {
                //console.log('Token refreshed.');
                jwtToken = result.Token;
            }
            else {
                jwtToken = '';
                clearInterval(tokenIntervalRef);
            }
        });
    }, 45000);
};
let notificationIntervalRef;
/* eslint-disable-next-line no-restricted-globals */
self.addEventListener('message', (event) => {
    const msgData = event.data;
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
        self.fetch(msgData.newNotificationUrl, requestOptions).then(result => result.json()).then(resulJson => {
            if (!!resulJson && (resulJson === null || resulJson === void 0 ? void 0 : resulJson.length) > 0) {
                /* eslint-disable-next-line no-restricted-globals */
                self.postMessage(resulJson);
                //Push Heading/Message
            }
        });
    }, 60000);
});
export {};