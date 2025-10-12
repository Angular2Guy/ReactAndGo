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
import Button from '@mui/material/Button';
import * as React from 'react';
import GlobalState from "~/GlobalState";
import styles from './login.module.css';
import { useAtom } from "jotai";

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
  const [userName, setUserName] = React.useState('');
  const [password1, setPassword1] = React.useState('');
  const [password2, setPassword2] = React.useState('');
  const [responseMsg, setResponseMsg] = React.useState('');
  const [activeTab, setActiveTab] = React.useState(0);

  const navToApp = () => {
    navigate("/app/app");
  }

  return (
    <div>
      <h1>Welcome to the App!</h1>
      <p>This is a simple React Router application.</p>      
      <Button variant="contained" color="primary" onClick={navToApp}>
        Login
      </Button>
    </div>
  );
}
