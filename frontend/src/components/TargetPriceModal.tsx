import { useRecoilState, useRecoilValue } from "recoil";
import GlobalState from "../GlobalState";
import { Box, TextField, Button, Dialog, DialogContent } from '@mui/material';
import { useState, useMemo, ChangeEventHandler, FormEvent } from "react";
import { UserRequest, UserResponse } from "./LoginModal";

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
        if(!!controller) {
            controller.abort();
        }
        const myDiesel = updatePrice(targetDiesel);
        const myE5 = updatePrice(targetE5);
        const myE10 = updatePrice(targetE10);
        const requestOptions = {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${globalJwtTokenState}` },
            body: JSON.stringify({ Username: globalUserNameState, Password: '', TargetDiesel: myDiesel, TargetE10: myE10, TargetE5: myE5 } as UserRequest)
        };
        controller = new AbortController();
        const response = await fetch('/appuser/targetprices', requestOptions);
        const result = await response.json() as UserResponse;
        controller = null;
        setGlobalUserDataState({
            Latitude: globalUserDataState.Latitude, Longitude: globalUserDataState.Longitude, SearchRadius: globalUserDataState.SearchRadius,
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