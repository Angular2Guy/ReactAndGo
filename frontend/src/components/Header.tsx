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
import styles from "./header.module.scss";
import {Button} from '@mui/material';
import {FormEvent} from 'react';
import { useSetRecoilState,useRecoilState } from "recoil";
import GlobalState from "../GlobalState";

const Header = () => {
    const setLocationModalState = useSetRecoilState(GlobalState.locationModalState);
    const setTargetPriceModalState = useSetRecoilState(GlobalState.targetPriceModalState);    
    const setJwtTokenState = useSetRecoilState(GlobalState.jwtTokenState);   
    const setGlobalLoginModal = useSetRecoilState(GlobalState.loginModalState);
    const [globalWebWorkerRefState, setGlobalWebWorkerRefState] = useRecoilState(GlobalState.webWorkerRefState); 

    const logout = (event: FormEvent) => {
        console.log("Logout ",event);
        setJwtTokenState('');    
        globalWebWorkerRefState?.postMessage({jwtToken: '', newNotificationUrl: ''});
        setGlobalWebWorkerRefState(null);
        setGlobalLoginModal(true);
    }
    const location = (event: FormEvent) => {
        //console.log("Location ",event);
        setLocationModalState(true)
    }
    const targetPrice = (event: FormEvent) => {
        setTargetPriceModalState(true);    
    }

    return <div className={styles.headerBase}>
        <span>Cheap Gas</span>
        <Button variant="outlined" onClick={logout} className={styles.headerButton}>Logout</Button>
        <Button variant="outlined" onClick={location} className={styles.headerButton}>Location</Button>
        <Button variant="outlined" onClick={targetPrice} className={styles.headerButton}>Target Price</Button>
    </div>
}

export default Header;