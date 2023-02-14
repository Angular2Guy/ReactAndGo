import {atom} from 'recoil';

const GlobalState = {
    userNameState: atom({
        key: 'userNameState',
        default: '',
    }),
    jwtTokenState: atom({
        key: 'jwtTokenState',
        default: '',
    }),
    locationModalState: atom({
        key: 'locationModalState',
        default: false
    })
}
export default GlobalState;