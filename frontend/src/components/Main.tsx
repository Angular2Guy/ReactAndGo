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
import { Tabs, Tab, Box } from '@mui/material';
import { useEffect, useState, SyntheticEvent } from 'react';
import DataTable, { TableDataRow } from './DataTable';
import GsMap, { GsValue } from './GsMap';
import { useRecoilRefresher_UNSTABLE, useRecoilValue } from 'recoil';
import GlobalState from '../GlobalState';
import styles from './main.module.scss';

interface GasPriceAvgs {
	Postcode:        string
	County:          string
	State:           string
	CountyAvgDiesel: number
	CountyAvgE10:    number
	CountyAvgE5:     number
	StateAvgDiesel:  number
	StateAvgE10:     number
	StateAvgE5:      number
}

interface GasStation {  
  StationName: string;
  Brand: string;
  Street: string;
  Place: string;
  HouseNumber: string;
  PostCode: string;
  Latitude: number;
  Longitude: number;
  PublicHolidayIdentifier: string;
  OtJson: string;
  FirstActive: Date;
  GasPrices: GasPrice[];
}

interface GasPrice {
  E5: number;
  E10: number;
  Diesel: number;
  Date: string;
  Changed: number;
}

interface Notification {
  Timestamp: Date;
  UserUuid: string;
  Title: string;
  Message: string;
  DataJson: string;
}

interface MyDataJson {
  StationName: string;
  Brand: string;
  Street: string;
  Place: string;
  HouseNumber: string;
  PostCode: string;
  Latitude: number;
  Longitude: number;
  E5: number;
  E10: number;
  Diesel: number;
  Timestamp: string;
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
        <Box sx={{ p: 3 }} className={styles.myText}>
          {children}
        </Box>
      )}
    </div>
  );
}

export default function Main() {  
  const [controller, setController] = useState(null as AbortController | null);
  const [timer, setTimer] = useState(undefined as undefined | NodeJS.Timer);
  const [value, setValue] = useState(0);
  const [first, setFirst] = useState(true);
  const [rows, setRows] = useState([] as TableDataRow[]);
//  const [avgTimeSlots, setAvgTimeSlots] = useState([])
  const [gsValues, setGsValues] = useState([] as GsValue[]);
  const globalJwtTokenState = useRecoilValue(GlobalState.jwtTokenState);
  const globalUserUuidState = useRecoilValue(GlobalState.userUuidState);
  const globalUserDataState = useRecoilValue(GlobalState.userDataState);    
  const refreshJwtTokenState = useRecoilRefresher_UNSTABLE(GlobalState.jwtTokenState);  


  const handleTabChange = (event: SyntheticEvent, newValue: number) => {
    setValue(newValue);
    clearInterval(timer);
    getData(newValue);
    setTimer(setInterval(() => getData(newValue), 10000));
  }

  const getData = (newValue: number) => {
    if (globalJwtTokenState?.length < 10 || globalUserUuidState?.length < 10) {
      return;
    }    
    //console.log(newValue); 
    if(!!controller) {
      controller.abort();
    }
    setController(new AbortController()); 
    refreshJwtTokenState();  
    // jwtToken = globalJwtTokenState; //When recoil makes refresh work.
    const jwtToken = !!GlobalState.jwtToken ? GlobalState.jwtToken : globalJwtTokenState;  
    if (newValue === 0 || newValue === 2) {           
      fetchSearchLocation(jwtToken);
    } else {     
      fetchLastMatches(jwtToken);
    };
  }

  const fetchSearchLocation = (jwtToken: string) => {    
    const requestOptions2 = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${jwtToken}` },
      body: JSON.stringify({ Longitude: globalUserDataState.Longitude, Latitude: globalUserDataState.Latitude, Radius: globalUserDataState.SearchRadius }),
      signal: controller?.signal
    }  
    let postcode = ''
    fetch('/gasstation/search/location', requestOptions2).then(myResult => myResult.json() as Promise<GasStation[]>).then(myJson => {
      const myResult = myJson.filter(value => value?.GasPrices?.length > 0).map(value => {
        postcode = value.PostCode;
        return value;
      }).map(value => ({        
        location: value.Place + ' ' + value.Brand + ' ' + value.Street + ' ' + value.HouseNumber, e5: value.GasPrices[0].E5,
        e10: value.GasPrices[0].E10, diesel: value.GasPrices[0].Diesel, date: new Date(Date.parse(value.GasPrices[0].Date)), longitude: value.Longitude, latitude: value.Latitude
      } as TableDataRow));      
      const requestOptions3 = {
        method: 'GET',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${jwtToken}` },        
        signal: controller?.signal
      }  
      fetch(`/gasprice/avgs/${postcode}`, requestOptions3).then(myResult => myResult.json() as Promise<GasPriceAvgs>).then(myJson => {
        const rowCounty = ({        
          location: myJson.County, e5: Math.round(myJson.CountyAvgE5), e10: Math.round(myJson.CountyAvgE10), diesel: Math.round(myJson.CountyAvgDiesel), date: new Date(), longitude: 0, latitude: 0
        } as TableDataRow);
        const rowState = ({        
          location: myJson.State, e5: Math.round(myJson.StateAvgE5), e10: Math.round(myJson.StateAvgE10), diesel: Math.round(myJson.StateAvgDiesel), date: new Date(), longitude: 0, latitude: 0
        } as TableDataRow);
        const resultRows = [rowCounty, rowState, ...myResult]
        setRows(resultRows);
      });         
      setGsValues(myResult);
    }).then(() => setController(null));
  }

  const fetchLastMatches = (jwtToken: string) => {    
    const requestOptions1 = {
      method: 'GET',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${jwtToken}` },
      signal: controller?.signal
    }
    fetch(`/usernotification/current/${globalUserUuidState}`, requestOptions1).then(myResult => myResult.json() as Promise<Notification[]>).then(myJson => {
      //console.log(myJson);
      const result = myJson.map(value => {
        //console.log(JSON.parse(value?.DataJson));
        return (JSON.parse(value?.DataJson) as MyDataJson[])?.map(value2 => {
          //console.log(value2);
          return {
            location: value2.Place + ' ' + value2.Brand + ' ' + value2.Street + ' ' + value2.HouseNumber,
            e5: value2.E5, e10: value2.E10, diesel: value2.Diesel, date: new Date(Date.parse(value2.Timestamp)), longitude: 0, latitude: 0
          } as TableDataRow;
        });
      })?.flat();
      setRows(result);
      //const result 
    }).then(() => setController(null));
  }

  // eslint-disable-next-line
  useEffect(() => {
    if (globalJwtTokenState?.length > 10 && globalUserUuidState?.length > 10 && first) {
      setTimeout(() => handleTabChange({} as unknown as SyntheticEvent, value), 3000);
      setFirst(false);
    }
  });

  return (<Box sx={{ width: '100%' }}>
    <Tabs value={value} onChange={handleTabChange} >
      <Tab label="Current Prices" />
      <Tab label="Last Price matches" />
      <Tab label="Current Prices Map" />
    </Tabs>
    <TabPanel value={value} index={0}>
      <DataTable diesel='Diesel' e10='E10' e5='E5' location='Location' showAverages={true} time='Time' rows={rows}></DataTable>
    </TabPanel>
    <TabPanel value={value} index={1}>
      <DataTable diesel='Diesel' e10='E10' e5='E5' location='Location' showAverages={true} time='Time' rows={rows}></DataTable>
    </TabPanel>
    <TabPanel value={value} index={2}>
      <GsMap gsValues={gsValues} center={globalUserDataState}></GsMap>      
    </TabPanel>
  </Box>);
}