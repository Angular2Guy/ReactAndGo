import { useRecoilState,useRecoilValue } from "recoil";
import GlobalState from "../GlobalState";
import {UserDataState} from "../GlobalState";
import { useMemo,useEffect,useState,FormEvent,ChangeEvent,SyntheticEvent } from "react";
import {Box,TextField,Button,Dialog,DialogContent, Autocomplete} from '@mui/material';
import {UserRequest, UserResponse} from "./LoginModal";
import styles from './modal.module.scss';

interface PostCodeLocation {
    Message: string;
	Longitude:  number;
	Latitude:  number;
	Label:      string;
	PostCode:   number
	SquareKM:   number;
    Population: number;
}

const LocationModal = () => {
    const [open, setOpen] = useState(false);
    const [searchRadius, setSearchRadius] = useState(0);
    const [longitude, setLongitude] = useState(0);
    const [latitude, setLatitude] = useState(0);
    const [options, setOptions] = useState([] as PostCodeLocation[]);       
    const [globalLocationModalState, setGlobalLocationModalState] = useRecoilState(GlobalState.locationModalState);
    const globalJwtTokenState = useRecoilValue(GlobalState.jwtTokenState);
    const [globalUserDataState, setGlobalUserDataState] = useRecoilState(GlobalState.userDataState);
    const globalUserNameState = useRecoilValue(GlobalState.userNameState);
    
    useEffect(() => {
        if (!open) {
          setOptions([]);
        }            
      }, [open]);      

    const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        //console.log("Submit: ",event);
        const requestOptions = {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${globalJwtTokenState}`},
            body: JSON.stringify({Username: globalUserNameState, Password: '', Latitude: latitude, Longitude: longitude, SearchRadius: searchRadius} as UserRequest)             
        };
        const response = await fetch('/appuser/locationradius', requestOptions);
        const userResponse = response.json() as UserResponse;
        setGlobalUserDataState({Latitude: userResponse.Latitude, Longitude: userResponse.Longitude, SearchRadius: userResponse.SearchRadius, 
            TargetDiesel: globalUserDataState.TargetDiesel, TargetE10: globalUserDataState.TargetE10, TargetE5: globalUserDataState.TargetE5} as UserDataState);
        setGlobalLocationModalState(false);
    }
/*
    const handleClose = (event: React.FormEvent) => {
        setGlobalLocationModalState(false);        
    } 
*/
    const handleChange = async (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
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
        const locations = await response.json() as PostCodeLocation[];        
        setOptions(!locations ? [] : locations);
        //console.log(locations);
    }

    const handleSearchRadiusChange = (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        //console.log(event?.currentTarget?.value);
        const mySearchRadius = parseFloat(event?.currentTarget?.value);
        setSearchRadius(Number.isNaN(mySearchRadius) ? searchRadius : mySearchRadius);
    }

    const handleOptionChange = (event: SyntheticEvent<Element, Event>, value: string) =>{               
        const filteredOptions = options.filter(option => option.Label === value);
        //console.log(filteredOptions);
        if(filteredOptions.length > 0) {
            setLongitude(filteredOptions[0].Longitude);
            setLatitude(filteredOptions[0].Latitude);
        }
    }

    const handleCancel = (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
        setSearchRadius(0);
        setLongitude(0);
        setLatitude(0);
        setGlobalLocationModalState(false);
    }

    const handleGetLocation = (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {        
        if(!!navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(result => {
                if(!!result?.coords?.longitude && !!result?.coords?.latitude) {
                setLongitude(result.coords.longitude);
                setLatitude(result.coords.latitude);
                }
            });
        }
    }

    let dialogOpen = useMemo(() => {        
        //console.log(globalUserDataState.Longitude+' '+globalUserDataState.Latitude);        
        setLongitude(globalUserDataState.Longitude);
        setLatitude(globalUserDataState.Latitude);              
        setSearchRadius(globalUserDataState.SearchRadius);           
        return globalLocationModalState;
    }, [globalLocationModalState, globalUserDataState.Longitude, globalUserDataState.Latitude, globalUserDataState.SearchRadius]);    

    return (<Dialog open={dialogOpen} className="backDrop">
        <DialogContent>
         <Box
      component="form"     
      noValidate
      autoComplete="off"
      onSubmit={handleSubmit}>        
        <Autocomplete
            open={open}
            onOpen={() => {
                setOpen(true);
              }}
            onClose={() => {
                setOpen(false);
              }}
            style={{ width: 300 }}              
            onInputChange={handleOptionChange}         
            getOptionLabel={option => option.Label}
            options={options}
            renderInput={(params) => <TextField {...params} label="Locations" onChange={handleChange} />}
        ></Autocomplete>
        <div>
            <h3>Longitude: {longitude}</h3>
            <h3>Latitude: {latitude}</h3>            
        </div>
         <TextField
            autoFocus
            margin="dense"
            value={searchRadius} 
            onChange={handleSearchRadiusChange}            
            label="Search Radius"
            type="string"
            fullWidth
            variant="standard"/>      
          <div>
            <Button type="submit">Ok</Button>
            <Button onClick={handleCancel}>Cancel</Button>              
            <Button className={styles.toright} onClick={handleGetLocation}>Get Location</Button>  
          </div>
    </Box>
    </DialogContent>
    </Dialog>);
}
export default LocationModal;