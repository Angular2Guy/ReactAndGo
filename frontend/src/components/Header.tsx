import styles from "./header.module.scss";
import {Button} from '@mui/material';
import { useSetRecoilState } from "recoil";
import GlobalState from "../GlobalState";

const Header = () => {
    const setLocationModalState = useSetRecoilState(GlobalState.locationModalState);
    const setTargetPriceModalState = useSetRecoilState(GlobalState.targetPriceModalState);    

    const logout = (event: React.FormEvent) => {
        console.log("Logout ",event);
    }
    const location = (event: React.FormEvent) => {
        console.log("Location ",event);
        setLocationModalState(true)
    }
    const targetPrice = (event: React.FormEvent) => {
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