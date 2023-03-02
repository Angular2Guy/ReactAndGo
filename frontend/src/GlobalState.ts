import {atom} from 'recoil';

export interface UserDataState {
    Longitude: number
	Latitude: number
	SearchRadius: number
	TargetDiesel: number
	TargetE5: number
	TargetE10: number
}

const GlobalState = {
    userNameState: atom({
        key: 'userNameState',
        default: '',
    }),
    userUuidState: atom({
        key: 'userUuidState',
        default: '',
    }),
    userDataState: atom({
        key: 'userDataState',
        default: {Latitude: 0.0, Longitude: 0.0, SearchRadius: 0, TargetDiesel: 0.0, TargetE10: 0.0, TargetE5: 0.0} as UserDataState,
    }),
    jwtTokenState: atom({
        key: 'jwtTokenState',
        default: '',
    }),
    locationModalState: atom({
        key: 'locationModalState',
        default: false
    }),
    targetPriceModalState: atom({
        key: 'targetPriceModalState',
        default: false
    }),
    webWorkerRefState: atom<null|Worker>({
        key: 'webWorkerRefState',    
        default: null
    })
}
export default GlobalState;