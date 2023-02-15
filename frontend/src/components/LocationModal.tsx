import { useRecoilState } from "recoil";
import GlobalState from "../GlobalState";
import { useEffect,useState } from "react";
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
    const [globalLocationModalState, setGlobalLocationModalState] = useRecoilState(GlobalState.locationModalState);
    const [globalJwtTokenState, setGlobalJwtTokenState] = useRecoilState(GlobalState.jwtTokenState);
    
    useEffect(() => {
        if (!open) {
          setOptions([]);
        }
      }, [open]);
    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        console.log("Submit: ",event);
    }

    const handleClose = (event: React.FormEvent) => {
        setGlobalLocationModalState(false);
    } 

    const handleChange = async (event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        event.preventDefault();
        if(!event?.currentTarget?.value) {
            setOptions([]);
            return;
        }
        const requestOptions = {
            method: 'GET',
            headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${globalJwtTokenState}` }            
        };
        const response = await fetch(`/appuser/location?location=${event.currentTarget.value}`, requestOptions);
        //console.log(event.currentTarget.value);
        console.log(response);
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
            getOptionLabel={option => option.Label}
            options={options}
            renderInput={(params) => <TextField {...params} label="Locations" onChange={ev => handleChange(ev)} />}
        ></Autocomplete>
    </Box>
    </DialogContent>
    </Dialog>);
}
export default LocationModal;