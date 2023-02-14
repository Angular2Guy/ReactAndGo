import styles from "./header.module.scss";
import {Button} from '@mui/material';
import { useRecoilState } from "recoil";
import GlobalState from "../GlobalState";

const Header = () => {
    const [globalLocationModalState, setLocationModalState] = useRecoilState(GlobalState.locationModalState);
    const logout = (event: React.FormEvent) => {
        console.log("Logout ",event);
    }
    const location = (event: React.FormEvent) => {
        console.log("Location ",event);
        setLocationModalState(true)
    }
    return <div className={styles.headerBase}>
        <span>Cheap Gas</span>
        <Button variant="outlined" onClick={logout} className={styles.headerButton}>Logout</Button>
        <Button variant="outlined" onClick={location} className={styles.headerButton}>Location</Button>
    </div>
}

export default Header;