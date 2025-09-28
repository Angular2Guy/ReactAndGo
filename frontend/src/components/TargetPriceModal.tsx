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
import { useRecoilState, useRecoilValue } from "recoil";
import GlobalState from "../GlobalState";
import { Box, TextField, Button, Dialog, DialogContent } from '@mui/material';
import { useState, useMemo, ChangeEventHandler, FormEvent } from "react";
import { UserRequest } from "../service/dtos";
import { postTargetPrices } from "../service/http-client";


const TargetPriceModal = () => {
    let controller: AbortController | null = null;
    const [targetDiesel, setTargetDiesel] = useState('0');
    const [targetE5, setTargetE5] = useState('0');
    const [targetE10, setTargetE10] = useState('0');
    const [globalTargetPriceModalState, setGlobalTargetPriceModalState] = useRecoilState(GlobalState.targetPriceModalState);
    const globalJwtTokenState = useRecoilValue(GlobalState.jwtTokenState);
    const globalUserNameState = useRecoilValue(GlobalState.userNameState);
    const [globalUserDataState, setGlobalUserDataState] = useRecoilState(GlobalState.userDataState);

    let dialogOpen = useMemo(() => {
        setTargetDiesel('' + globalUserDataState.TargetDiesel);
        setTargetE10('' + globalUserDataState.TargetE10);
        setTargetE5('' + globalUserDataState.TargetE5);
        return globalTargetPriceModalState;
    }, [globalTargetPriceModalState, globalUserDataState.TargetDiesel, globalUserDataState.TargetE10, globalUserDataState.TargetE5]);

    const handleTargetDieselChange: ChangeEventHandler<HTMLInputElement> = (event) => {
        event.preventDefault();
        setTargetDiesel(event.currentTarget.value);
    }

    const handleTargetE10Change: ChangeEventHandler<HTMLInputElement> = (event) => {
        event.preventDefault();
        setTargetE10(event.currentTarget.value);
    }

    const handleTargetE5Change: ChangeEventHandler<HTMLInputElement> = (event) => {
        event.preventDefault();
        setTargetE5(event.currentTarget.value);
    }

    const updatePrice = (priceStr: string) => {
        let myPrice = priceStr.replace(/\.|,/, '');
        while (myPrice.length < 4) {
            myPrice = myPrice + '0';
        }
        return myPrice;
    }

    const handleSubmit = async (event: FormEvent) => {
        event.preventDefault();
        if (!!controller) {
            controller.abort();
        }
        const myDiesel = updatePrice(targetDiesel);
        const myE5 = updatePrice(targetE5);
        const myE10 = updatePrice(targetE10);
        controller = new AbortController();
        const requestString = JSON.stringify({ Username: globalUserNameState, Password: '', TargetDiesel: myDiesel, TargetE10: myE10, TargetE5: myE5 } as UserRequest);
        const result = await postTargetPrices(globalJwtTokenState, controller, requestString);
        controller = null;
        setGlobalUserDataState({
            Latitude: globalUserDataState.Latitude, Longitude: globalUserDataState.Longitude, SearchRadius: globalUserDataState.SearchRadius, PostCode: globalUserDataState.PostCode,
            TargetDiesel: !result.TargetDiesel ? 0 : result.TargetDiesel, TargetE10: !result.TargetE10 ? 0 : result.TargetE10, TargetE5: !result.TargetE5 ? 0 : result.TargetE5
        });
        setGlobalTargetPriceModalState(false);
        //console.log(result.TargetDiesel+' '+result.TargetE10+' '+result.TargetE5);        
    };

    const handleCancel = (event: FormEvent) => {
        event.preventDefault();
        setGlobalTargetPriceModalState(false);
        setTargetDiesel('0');
        setTargetE10('0');
        setTargetE5('0');
    }

    return (<Dialog open={dialogOpen} className="backDrop">
        <DialogContent>
            <Box
                component="form"
                noValidate
                autoComplete="off"
                onSubmit={handleSubmit}>
                <div>
                    <TextField
                        autoFocus
                        margin="dense"
                        value={targetDiesel}
                        onChange={handleTargetDieselChange}
                        label="Targetprice Diesel"
                        type="string"
                        fullWidth
                        variant="standard" />
                    <TextField
                        autoFocus
                        margin="dense"
                        value={targetE5}
                        onChange={handleTargetE5Change}
                        label="Targetprice E5"
                        type="string"
                        fullWidth
                        variant="standard" />
                    <TextField
                        autoFocus
                        margin="dense"
                        value={targetE10}
                        onChange={handleTargetE10Change}
                        label="Targetprice E10"
                        type="string"
                        fullWidth
                        variant="standard" />
                </div>
                <div>
                    <Button type="submit">Ok</Button>
                    <Button onClick={handleCancel}>Cancel</Button>
                </div>
            </Box>
        </DialogContent>
    </Dialog>);
}

export default TargetPriceModal;