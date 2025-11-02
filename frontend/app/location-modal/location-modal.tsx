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
import * as React from 'react';
import GlobalState from "../GlobalState";
import { type UserDataState} from "../GlobalState";
import { useMemo,useEffect,useState,type FormEvent, type ChangeEvent,type SyntheticEvent } from "react";
import {Box,TextField,Button,Dialog,DialogContent, Autocomplete} from '@mui/material';
import styles from './location-modal.module.css';
import { fetchLocation, postLocationRadius } from "../service/http-client";
import { type PostCodeLocation } from "../model/location";
import { type UserRequest } from "../model/user";
import { useAtom } from "jotai";
import { useTranslation } from 'node_modules/react-i18next';
import { useNavigate } from 'react-router';

const LocationModal = () => {
    const navigate = useNavigate();
    let controller: AbortController | null = null;
    const { t } = useTranslation();
    const [open, setOpen] = useState(false);
    const [searchRadius, setSearchRadius] = useState(0);
    const [longitude, setLongitude] = useState(0);
    const [latitude, setLatitude] = useState(0);
    const [postCode, setPostCode] = useState('');
    const [options, setOptions] = useState([] as PostCodeLocation[]);       
    const [globalLocationModalState, setGlobalLocationModalState] = useAtom(GlobalState.locationModalState);
    const [globalJwtTokenState, setGlobalJwtTokenState] = useAtom(GlobalState.jwtTokenState);
    const [globalUserDataState, setGlobalUserDataState] = useAtom(GlobalState.userDataState);
    const [globalUserNameState, setGlobalUserNameState] = useAtom(GlobalState.userNameState);

    useEffect(() => {
        if (!open) {
          setOptions([]);
        }            
      }, [open]);      

    const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        if(!!controller) {
            controller.abort();
        }
        controller = new AbortController();
        //console.log("Submit: ",event);
        const requestString = JSON.stringify({Username: globalUserNameState, Password: '', Latitude: latitude, Longitude: longitude, SearchRadius: searchRadius, PostCode: parseInt(postCode)} as UserRequest);
        const userResponse = await postLocationRadius(navigate, globalJwtTokenState, controller, requestString);
        controller = null;
        setGlobalUserDataState({Latitude: userResponse.Latitude, Longitude: userResponse.Longitude, SearchRadius: userResponse.SearchRadius, PostCode: postCode.toString() || 0,
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
                if(!!controller) {
            controller.abort();
        } 
        controller = new AbortController();
        const locations = await fetchLocation(navigate, globalJwtTokenState, controller, event.currentTarget.value);
        setOptions(!locations ? [] : locations);
        controller = null;
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
            setPostCode(formatPostCode(filteredOptions[0].PostCode));
        }
    }

    const formatPostCode = (myPlz: number) => {                
        return '00000'.substring(0, 5 - myPlz?.toString()?.length > 0 ? myPlz?.toString()?.length : 0) + myPlz.toString();
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
        setPostCode(formatPostCode(globalUserDataState.PostCode));          
        return globalLocationModalState;
    }, [globalLocationModalState, globalUserDataState.Longitude, globalUserDataState.Latitude, globalUserDataState.SearchRadius, globalUserDataState.PostCode]);    

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
            <h3>{t('location.longitude')}: {longitude}</h3>
            <h3>{t('location.latitude')}: {latitude}</h3>
            <h3>{t('location.postalCode')}: {postCode}</h3>
        </div>
         <TextField
            autoFocus
            margin="dense"
            value={searchRadius} 
            onChange={handleSearchRadiusChange}            
            label={t('location.searchRadius')}
            type="string"
            fullWidth
            variant="standard"/>      
          <div>
            <Button type="submit">{t('common.ok')}</Button>
            <Button onClick={handleCancel}>{t('common.cancel')}</Button>              
            <Button className={styles.toright} onClick={handleGetLocation}>{t('location.getLocation')}</Button>  
          </div>
    </Box>
    </DialogContent>
    </Dialog>);
}
export default LocationModal;
