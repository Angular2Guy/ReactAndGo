import {atom} from 'recoil';

const GlobalState = {
    userNameState: atom({
        key: 'userNameState',
        default: '',
    }),
    jwtTokenState: atom({
        key: 'jwtTokenState',
        default: '',
    })
}
export default GlobalState;