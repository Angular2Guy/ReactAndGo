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