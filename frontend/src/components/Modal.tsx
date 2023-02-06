import { useRecoilState } from "recoil";
import styles from './modal.module.scss';
import GlobalState from "../GlobalState";
import { useState } from "react";
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogContent from '@mui/material/DialogContent';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import TextField from '@mui/material/TextField';
import Box from '@mui/material/Box';
import DialogActions from '@mui/material/DialogActions';


interface TabPanelProps {
    children?: React.ReactNode;
    index: number;
    value: number;
  }
  
  function TabPanel(props: TabPanelProps) {
    const { children, value, index, ...other } = props;
  
    return (
      <div
        role="tabpanel"
        hidden={value !== index}
        id={`simple-tabpanel-${index}`}
        aria-labelledby={`simple-tab-${index}`}
        {...other}
      >
        {value === index && (
          <>{children}</>
        )}
      </div>
    );
  }

const Modal = () => {
   const [globalUserName, setGlobalUserName] = useRecoilState(GlobalState.userNameState);
   const [userName, setUserName] = useState('');
   const [password1, setPassword1] = useState('');
   const [password2, setPassword2] = useState('');
   const [open, setOpen] = useState(true);
   const [activeTab, setActiveTab] = useState(0);

   const handleChangeUsername: React.ChangeEventHandler<HTMLInputElement> = (event) => {
      setUserName(event.currentTarget.value as string);      
  };
  const handleChangePassword1: React.ChangeEventHandler<HTMLInputElement> = (event) => {
    setPassword1(event.currentTarget.value as string);      
};
const handleChangePassword2: React.ChangeEventHandler<HTMLInputElement> = (event) => {
    setPassword2(event.currentTarget.value as string);      
};
  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
      event.preventDefault();      
      setGlobalUserName(userName);
      setUserName('');
      setOpen(false);      
  }
  const handleCancel = (event: React.FormEvent) => {
   event.preventDefault();
      setUserName('');
      setPassword1('');
      setPassword2('');
  };
  const handleClose = () => {
    //setOpen(false);
  };
  const handleTabChange = (event: React.BaseSyntheticEvent, newValue: number) => {
    setActiveTab(newValue);
  };
  const a11yProps = (index: number) => {
    return {
      id: `simple-tab-${index}`,
      'aria-controls': `simple-tabpanel-${index}`,
    };
  }
     return (<Dialog open={open} onClose={handleClose} className="backDrop">
        <DialogContent>
        <Tabs value={activeTab} onChange={handleTabChange} aria-label="basic tabs example">
          <Tab label="Login" {...a11yProps(0)} />
          <Tab label="Singin" {...a11yProps(1)} />
        </Tabs>      
      <TabPanel value={activeTab} index={0}>
      <Box
      component="form"     
      noValidate
      autoComplete="off"
      onSubmit={handleSubmit}
    >      
      <TextField
            autoFocus
            margin="dense"
            value={userName} 
            onChange={handleChangeUsername}
            id="userName"
            label="user name"
            type="string"
            fullWidth
            variant="standard"
          />     
          <TextField
            autoFocus
            margin="dense"
            value={password1} 
            onChange={handleChangePassword1}
            id="password1"
            label="password"
            type="password"
            fullWidth
            variant="standard"
          />      
            <Button type="submit">Ok</Button>
          <Button onClick={handleCancel}>Cancel</Button>
          <div>GlobalUserName: {globalUserName}</div>
      </Box>
      </TabPanel>
      <TabPanel value={activeTab} index={1}>
      <Box
      component="form"     
      noValidate
      autoComplete="off"
      onSubmit={handleSubmit}
    >      
      <TextField
            autoFocus
            margin="dense"
            value={userName} 
            onChange={handleChangeUsername}
            id="userName"
            label="user name"
            type="string"
            fullWidth
            variant="standard"
          />     
          <TextField
            autoFocus
            margin="dense"
            value={password1} 
            onChange={handleChangePassword1}
            id="password1"
            label="password"
            type="password"
            fullWidth
            variant="standard"
          />
        <TextField
            autoFocus
            margin="dense"
            value={password2} 
            onChange={handleChangePassword2}
            id="password2"
            label="password"
            type="password"
            fullWidth
            variant="standard"
          />      
            <Button type="submit">Ok</Button>
          <Button onClick={handleCancel}>Cancel</Button>
          <div>GlobalUserName: {globalUserName}</div>
      </Box>
      </TabPanel>      
        </DialogContent>
        </Dialog>);
}

export default Modal;