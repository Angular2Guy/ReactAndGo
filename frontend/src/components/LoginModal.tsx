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
import { useSetRecoilState, useRecoilState } from "recoil";
import styles from './modal.module.scss';
import GlobalState from "../GlobalState";
import { UserDataState } from "../GlobalState";
import { useState, ChangeEventHandler, FormEvent, BaseSyntheticEvent } from "react";
import { Box, TextField, Button, Tab, Tabs, Dialog, DialogContent } from '@mui/material';
//import { Token } from "@mui/icons-material";

export interface UserRequest {
  Username: string;
  Password: string;
  Latitude?: number;
  Longitude?: number;
  SearchRadius?: number;
  PostCode?: number;
  TargetDiesel?: string;
  TargetE10?: string;
  TargetE5?: string;
}

export interface UserResponse {
  Token?: string;
  Message?: string;
  PostCode?: number;
  Uuid?: string;
  Longitude?: number;
  Latitude?: number;
  SearchRadius?: number;
  TargetDiesel?: number;
  TargetE5?: number;
  TargetE10?: number;
}

interface MsgData {
  jwtToken?: string;
  newNotificationUrl?: string;
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
  let controller: AbortController | null = null;
  const setGlobalUserName = useSetRecoilState(GlobalState.userNameState);
  const setGlobalUuid = useSetRecoilState(GlobalState.userUuidState);
  const setGlobalJwtToken = useSetRecoilState(GlobalState.jwtTokenState);
  const setGlobalUserDataState = useSetRecoilState(GlobalState.userDataState);
  const [globalWebWorkerRefState, setGlobalWebWorkerRefState] = useRecoilState(GlobalState.webWorkerRefState);
  const [userName, setUserName] = useState('');
  const [password1, setPassword1] = useState('');
  const [password2, setPassword2] = useState('');
  const [responseMsg, setResponseMsg] = useState('');
  const [globalLoginModal, setGlobalLoginModal] = useRecoilState(GlobalState.loginModalState);
  const [activeTab, setActiveTab] = useState(0);

  const handleChangeUsername: ChangeEventHandler<HTMLInputElement> = (event) => {
    setUserName(event.currentTarget.value as string);
  };
  const handleChangePassword1: ChangeEventHandler<HTMLInputElement> = (event) => {
    setPassword1(event.currentTarget.value as string);
  };
  const handleChangePassword2: ChangeEventHandler<HTMLInputElement> = (event) => {
    setPassword2(event.currentTarget.value as string);
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    Notification.requestPermission();
    if(!!controller) {
      controller.abort();
    }
    const requestOptions = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ Username: userName, Password: password1 } as UserRequest)
    };
    setResponseMsg('');
    controller = new AbortController();
    const httpResponse = activeTab === 0 ? await fetch('/appuser/login', requestOptions) : await fetch('/appuser/signin', requestOptions);
    const userResponse = await httpResponse.json() as UserResponse;
    controller = null;
    //console.log(userResponse);
    if (!userResponse?.Message && !!userResponse?.Token && userResponse.Token?.length > 10 && !!userResponse?.Uuid && userResponse.Uuid?.length > 10) {
      setGlobalUserName(userName);
      setGlobalJwtToken(userResponse.Token);      
      setGlobalUuid(userResponse.Uuid);
      setGlobalUserDataState({
        Latitude: userResponse.Latitude, Longitude: userResponse.Longitude, SearchRadius: userResponse.SearchRadius, PostCode: userResponse.PostCode,
        TargetDiesel: userResponse.TargetDiesel, TargetE10: userResponse.TargetE10, TargetE5: userResponse.TargetE5
      } as UserDataState);
      setGlobalLoginModal(false);
      initWebWorker(userResponse);
      setUserName('');
      setPassword1('');
      setPassword2('');
    } else if (!!userResponse?.Message) {
      setResponseMsg(userResponse.Message);
    }
  }

  const initWebWorker = async (userResponse: UserResponse) => {
    let result = null;
    if (!globalWebWorkerRefState) {
      const worker = new Worker(new URL('../webpush/dedicated-worker.js', import.meta.url));
      if (!!worker) {
        worker.addEventListener('message', (event: MessageEvent) => {
          //console.log(event.data);
          if (!!event?.data?.Token && event?.data.Token?.length > 10) {
            setGlobalJwtToken(event.data.Token);
            GlobalState.jwtToken = event.data.Token;
          }
        });
        worker.postMessage({ jwtToken: userResponse.Token, newNotificationUrl: `/usernotification/new/${userResponse.Uuid}` } as MsgData);
        setGlobalWebWorkerRefState(worker);
        result = worker;
      }
    } else {
      globalWebWorkerRefState.postMessage({ jwtToken: userResponse.Token, newNotificationUrl: `/usernotification/new/${userResponse.Uuid}` } as MsgData);
      result = globalWebWorkerRefState;
    }
    return result;
  };

  const handleCancel = (event: FormEvent) => {
    event.preventDefault();
    setUserName('');
    setPassword1('');
    setPassword2('');
    setResponseMsg('');
  };
  const handleClose = () => {
    //setOpen(false);
  };
  const handleTabChange = (event: BaseSyntheticEvent, newValue: number) => {
    setActiveTab(newValue);
    setResponseMsg('');
  };
  const a11yProps = (index: number) => {
    return {
      id: `simple-tab-${index}`,
      'aria-controls': `simple-tabpanel-${index}`,
    };
  }
  return (<Dialog open={globalLoginModal} onClose={handleClose} className="backDrop">
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
            {[responseMsg].filter(value => !!value).map((value, index) =>
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
            {[responseMsg].filter(value => !!value).map((value, index) =>
              <span key={index}>Message: {value}</span>
            )}
          </div>
        </Box>
      </TabPanel>
    </DialogContent>
  </Dialog>);
}

export default LoginModal;