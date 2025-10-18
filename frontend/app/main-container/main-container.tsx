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
import styles from './main-container.module.css';
import Header from '../header/header';
//import LocationModal from './components/LocationModal';
//import TargetPriceModal from './components/TargetPriceModal';
import Main from '../main/main';
import { useAtom } from "jotai";
import GlobalState from '../GlobalState';
import { useEffect } from 'react';
import { useNavigate } from 'react-router';

export interface TodoItem1 {
  name: string;
  id: string;
}

function MainContainer() {
  const [globalJwtTokenState, setGlobalJwtTokenState] = useAtom(GlobalState.jwtTokenState);
  const [globalUserUuidState, setGlobalUserUuidState] = useAtom(GlobalState.userUuidState);
  const navigate = useNavigate();
  
  useEffect(() => {
    if((!globalJwtTokenState || !globalUserUuidState || globalJwtTokenState.length < 10 || globalUserUuidState.length < 10)) {
      navigate('/');
    }
  });

  return (
    <div className="App">      
      <Header/>      
      {/* <LocationModal/> */}
      {/* <TargetPriceModal/> */}
      <Main/>      
    </div>
  );
}

export default MainContainer;