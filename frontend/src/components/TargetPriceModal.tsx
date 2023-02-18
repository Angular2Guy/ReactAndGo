import { useRecoilState,useRecoilValue } from "recoil";
import GlobalState from "../GlobalState";
import globalUserDataState from "../GlobalState";
import {Box,TextField,Button,Dialog,DialogContent} from '@mui/material';
import {useState,useMemo} from "react";

const TargetPriceModal = () => {
    const [targetDiesel, setTargetDiesel] = useState('0');
    const [targetE5, setTargetE5] = useState('0');
    const [targetE10, setTargetE10] = useState('0');
    const [globalTargetPriceModalState, setGlobalTargetPriceModalState] = useRecoilState(GlobalState.targetPriceModalState);
    const globalJwtTokenState = useRecoilValue(GlobalState.jwtTokenState);
    const [globalUserDataState, setGlobalUserDataState] = useRecoilState(GlobalState.userDataState);

    let dialogOpen = useMemo(() => {        
        setTargetDiesel(''+globalUserDataState.TargetDiesel);
        setTargetE10(''+globalUserDataState.TargetE10);      
        setTargetE5(''+globalUserDataState.TargetE5);           
        return globalTargetPriceModalState;
    }, [globalTargetPriceModalState, globalUserDataState.TargetDiesel, globalUserDataState.TargetE10, globalUserDataState.TargetE5]);    

    const handleTargetDieselChange: React.ChangeEventHandler<HTMLInputElement> = (event) => {
        event.preventDefault();
        setTargetDiesel(event.currentTarget.value);
    }

    const handleTargetE10Change: React.ChangeEventHandler<HTMLInputElement> = (event) => {
        event.preventDefault();
        setTargetE10(event.currentTarget.value);
    }

    const handleTargetE5Change: React.ChangeEventHandler<HTMLInputElement> = (event) => {
        event.preventDefault();
        setTargetE5(event.currentTarget.value);
    }

    const handleSubmit = (event: React.FormEvent) => {
        event.preventDefault();
        console.log(targetDiesel);
        console.log(targetE5);
        console.log(targetE10);
    };

    const handleCancel = (event: React.FormEvent) => {
        event.preventDefault();
        setGlobalTargetPriceModalState(false);
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
            variant="standard"/>
    <TextField
            autoFocus
            margin="dense"
            value={targetE5} 
            onChange={handleTargetE5Change}            
            label="Targetprice E5"
            type="string"
            fullWidth
            variant="standard"/>
    <TextField
            autoFocus
            margin="dense"
            value={targetE10} 
            onChange={handleTargetE10Change}            
            label="Targetprice E10"
            type="string"
            fullWidth
            variant="standard"/>        
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