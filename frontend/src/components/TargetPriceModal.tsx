import { useRecoilState } from "recoil";
import GlobalState from "../GlobalState";
import {Box,TextField,Button,Dialog,DialogContent} from '@mui/material';

const TargetPriceModal = () => {
    const [globalTargetPriceModalState, setGlobalTargetPriceModalState] = useRecoilState(GlobalState.targetPriceModalState);

    const handleSubmit = (event: React.FormEvent) => {
        event.preventDefault();
    };

    const handleCancel = (event: React.FormEvent) => {
        setGlobalTargetPriceModalState(false);
    }

    return (<Dialog open={globalTargetPriceModalState} className="backDrop">
    <DialogContent>
     <Box
  component="form"     
  noValidate
  autoComplete="off"
  onSubmit={handleSubmit}>            
    <div>
        <h3>Target Price</h3>        
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