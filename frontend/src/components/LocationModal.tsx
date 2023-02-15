import { useRecoilState } from "recoil";
import GlobalState from "../GlobalState";
import { useState } from "react";
import {Box,TextField,Button,Dialog,DialogContent, Autocomplete} from '@mui/material';

interface PostCodeLocation {
	Longitude:  number;
	Latitude:  number;
	Label:      string;
	PostCode:   number
	SquareKM:   number;
    Population: number;
}

const LocationModal = () => {
    const [open, setOpen] = useState(false);
    const [options, setOptions] = useState([] as PostCodeLocation[]);
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
        <Autocomplete
            open={open}
            onOpen={() => {
                setOpen(true);
              }}
            onClose={() => {
                setOpen(false);
              }}
            style={{ width: 300 }}
              options={options}
              renderInput={(params) => <TextField {...params} label="Locations" />}
        ></Autocomplete>
    </Box>
    </DialogContent>
    </Dialog>);
}
export default LocationModal;