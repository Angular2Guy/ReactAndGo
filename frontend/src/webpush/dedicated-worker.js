let jwtToken = '';
let tokenIntervalRef;
const refreshToken = (myToken) => {
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
            fetch('/appuser/refreshtoken', requestOptions).then(response => response.json()).then(result => {
                if ((!result.Message && !!result.Token && result.Token.length > 10)) {
                    //console.log('Token refreshed.');
                    jwtToken = result.Token;
                    /* eslint-disable-next-line no-restricted-globals */
                    self.postMessage(result);
                }
                else {
                    jwtToken = '';
                    clearInterval(tokenIntervalRef);
                }
            });
        }, 45000);
    }
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
        self.fetch(msgData.newNotificationUrl, requestOptions).then(result => result.json()).then(resultJson => {
            var _a, _b, _c, _d;
            if (!!resultJson && (resultJson === null || resultJson === void 0 ? void 0 : resultJson.length) > 0) {
                /* eslint-disable-next-line no-restricted-globals */
                self.postMessage(resultJson);
                //Notification
                //console.log(Notification.permission);
                if (Notification.permission === 'granted') {
                    if ((resultJson === null || resultJson === void 0 ? void 0 : resultJson.length) > 0 && ((_b = (_a = resultJson[0]) === null || _a === void 0 ? void 0 : _a.Message) === null || _b === void 0 ? void 0 : _b.length) > 1 && ((_d = (_c = resultJson[0]) === null || _c === void 0 ? void 0 : _c.Title) === null || _d === void 0 ? void 0 : _d.length) > 1) {
                        for (let value of resultJson) {
                            new Notification(value === null || value === void 0 ? void 0 : value.Title, { body: value === null || value === void 0 ? void 0 : value.Message });
                        }
                    }
                }
            }
        });
    }, 60000);
});
export {};