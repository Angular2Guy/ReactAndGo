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
import { useNavigate } from "react-router";
import { Dialog, DialogContent, Button, Tabs, Tab, Box, TextField } from '@mui/material';
import * as React from 'react';
import {useState} from 'react';
import GlobalState, { type UserDataState } from "~/GlobalState";
import styles from './login.module.css';
import { useAtom } from "jotai";
import type { FormEvent, BaseSyntheticEvent, ChangeEventHandler } from "react";
import { postLogin, postSignin } from "~/service/http-client";
import type { UserResponse } from "~/model/user";

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

export function Login() {
  const navigate = useNavigate();
  let controller: AbortController | null = null;  
  const [globalUserName,setGlobalUserName] = useAtom(GlobalState.userNameState);
  const [globalUuid,setGlobalUuid] = useAtom(GlobalState.userUuidState);
  const [globalJwtToken,setGlobalJwtToken] = useAtom(GlobalState.jwtTokenState);
  const [globalUserDataState,setGlobalUserDataState] = useAtom(GlobalState.userDataState);
  const [globalWebWorkerRefState, setGlobalWebWorkerRefState] = useAtom(GlobalState.webWorkerRefState);
  const [globalLoginModal, setGlobalLoginModal] = useAtom(GlobalState.loginModalState);
  const [userName, setUserName] = useState('');
  const [password1, setPassword1] = useState('');
  const [password2, setPassword2] = useState('');
  const [responseMsg, setResponseMsg] = useState('');
  const [activeTab, setActiveTab] = useState(0);

  const navToApp = () => {
    navigate("/app/app");
  }
  
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
    setResponseMsg('');
    controller = new AbortController();
    const userResponse = activeTab === 0 ? await postLogin(userName, password1, controller) : await postSignin(userName, password1, controller);    
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
      navigate('/app');
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
      <Button onClick={navToApp}>Test</Button>
    </DialogContent>
  </Dialog>);
}
