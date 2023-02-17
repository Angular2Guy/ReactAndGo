import { useRecoilState,useSetRecoilState } from "recoil";
import styles from './modal.module.scss';
import GlobalState from "../GlobalState";
import {UserDataState} from "../GlobalState";
import { useState } from "react";
import {Box,TextField,Button,Tab,Tabs,Dialog,DialogContent} from '@mui/material';
//import { Token } from "@mui/icons-material";

export interface UserRequest {
  Username:  string
	Password:  string
	Latitude?: number
	Longitude?: number
  SearchRadius?: number
}

export interface UserResponse {
  Token?:  string
	Message?: string
  Longitude?: number
	Latitude?: number
	SearchRadius?: number
	TargetDiesel?: number
	TargetE5?: number
	TargetE10?: number
}

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

const LoginModal = () => {
   const setGlobalUserName = useSetRecoilState(GlobalState.userNameState);
   const setGlobalJwtToken = useSetRecoilState(GlobalState.jwtTokenState);
   const [globalUserDataState, setGlobalUserDataState] = useRecoilState(GlobalState.userDataState);
   const [userName, setUserName] = useState('');
   const [password1, setPassword1] = useState('');
   const [password2, setPassword2] = useState('');
   const [responseMsg, setResponseMsg] = useState('');
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
  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
      event.preventDefault();      
      const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json'},
        body: JSON.stringify({ Username: userName, Password: password1 } as UserRequest)
    };
    setResponseMsg('');
    const httpResponse = activeTab === 0 ? await fetch('/appuser/login', requestOptions) : await fetch('/appuser/signin', requestOptions);
    const userResponse = await httpResponse.json() as UserResponse;
    console.log(userResponse);
    if(!userResponse?.Message && !!userResponse?.Token && userResponse.Token?.length > 10) {
      setGlobalUserName(userName);  
      setGlobalJwtToken(userResponse.Token);    
      setGlobalUserDataState({Latitude: userResponse.Latitude, Longitude: userResponse.Longitude, SearchRadius: userResponse.SearchRadius,
        TargetDiesel: globalUserDataState.TargetDiesel, TargetE10: globalUserDataState.TargetE10, TargetE5: globalUserDataState.TargetE5} as UserDataState);
      setUserName('');
      setOpen(false);               
    } else if(!!userResponse?.Message) {
      setResponseMsg(userResponse.Message);
    }
  }
  const handleCancel = (event: React.FormEvent) => {
   event.preventDefault();
      setUserName('');
      setPassword1('');
      setPassword2('');
      setResponseMsg('');
  };
  const handleClose = () => {
    //setOpen(false);
  };
  const handleTabChange = (event: React.BaseSyntheticEvent, newValue: number) => {
    setActiveTab(newValue);
    setResponseMsg('');
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
          <div>
            <Button type="submit">Ok</Button>
          <Button onClick={handleCancel}>Cancel</Button>   
          </div>                
          <div className={styles.responseMsg}>
          {[responseMsg].filter(value => !!value).map((value,index) => 
              <span key={index}>Message: {value}</span>
              )}          
          </div> 
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
          <div>
            <Button type="submit">Ok</Button>
          <Button onClick={handleCancel}>Cancel</Button>          
          </div>  
          <div className={styles.responseMsg}>
          {[responseMsg].filter(value => !!value).map((value,index) => 
              <span key={index}>Message: {value}</span>
              )}
          </div> 
      </Box>
      </TabPanel>      
        </DialogContent>
        </Dialog>);
}

export default LoginModal;