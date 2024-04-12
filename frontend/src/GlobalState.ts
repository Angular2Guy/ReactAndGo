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
import {atom} from 'recoil';

export interface UserDataState {
    Longitude: number;
	Latitude: number;
	SearchRadius: number;
    PostCode: number;
	TargetDiesel: number;
	TargetE5: number;
	TargetE10: number;
}

const GlobalState = {
    jwtToken: '',
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
        default: {Latitude: 0.0, Longitude: 0.0, SearchRadius: 0, PostCode: 0, TargetDiesel: 0.0, TargetE10: 0.0, TargetE5: 0.0} as UserDataState,
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
    loginModalState: atom({
        key: 'loginModalState',
        default: true
    }),
    webWorkerRefState: atom<null|Worker>({
        key: 'webWorkerRefState',    
        default: null
    }),    
}
export default GlobalState;