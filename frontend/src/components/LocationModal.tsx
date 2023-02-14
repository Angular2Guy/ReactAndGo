import { useRecoilState } from "recoil";
import GlobalState from "../GlobalState";
import {Box,TextField,Tab,Tabs,Button,Dialog,DialogContent, Autocomplete} from '@mui/material';

const LocationModal = () => {
    const [globalLocationModalState, setLocationModalState] = useRecoilState(GlobalState.locationModalState);

    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        console.log("Submit: ",event);
    }

    const handleClose = (event: React.FormEvent) => {
        setLocationModalState(false);
    }

    return (<Dialog open={globalLocationModalState} className="backDrop">
        <DialogContent>
         <Box
      component="form"     
      noValidate
      autoComplete="off"
      onSubmit={handleSubmit}>
        <div onClick={handleClose}>Hallo</div>
    </Box>
    </DialogContent>
    </Dialog>);
}
export default LocationModal;