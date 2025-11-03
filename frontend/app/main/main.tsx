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
import { Box, Tab, Tabs } from '@mui/material';
import { useEffect, useState, type SyntheticEvent } from 'react';
import DataTable, { type TableDataRow } from '../data-table/data-table';
import GsMap from '../gs-map/gs-map';
import GlobalState from '../GlobalState';
//import styles from './main.module.css';
import Chart from '../chart/chart';
import { useNavigate } from 'react-router';
import { fetchGasStations, fetchPriceAvgs, fetchTimeSlots, fetchUserNotifications } from '../api/http-client';
import { type TimeSlot } from '../model/time-slot-response';
import { type GsValue } from '../model/gs-point';
import { type MyDataJson } from '../model/my-data-json';
import { useAtom } from "jotai";
import { useTranslation } from 'node_modules/react-i18next';

export default function Main() {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [controller, setController] = useState(null as AbortController | null);
  const [timer, setTimer] = useState(undefined as undefined | NodeJS.Timeout);
  const [value, setValue] = useState(0);
  const [first, setFirst] = useState(true);
  const [rows, setRows] = useState([] as TableDataRow[]);
  const [avgTimeSlots, setAvgTimeSlots] = useState([] as TimeSlot[])
  const [gsValues, setGsValues] = useState([] as GsValue[]);
  const [globalJwtTokenState, setGlobalJwtTokenState] = useAtom(GlobalState.jwtTokenState);
  const [globalUserUuidState, setGlobalUserUuidState] = useAtom(GlobalState.userUuidState);
  const [globalUserDataState, setGlobalUserDataState] = useAtom(GlobalState.userDataState);

  const handleTabChange = (event: SyntheticEvent, newValue: number) => {
    setValue(newValue);
    clearInterval(timer);
    getData(newValue);
    setTimer(setInterval(() => getData(newValue), 10000));
  }

  const formatPostCode = (myPlz: number) => {
    return '00000'.substring(0, 5 - myPlz?.toString()?.length > 0 ? myPlz?.toString()?.length : 0) + myPlz.toString();
  }

  const getData = (newValue: number) => {
    if (globalJwtTokenState?.length < 10 || globalUserUuidState?.length < 10) {
      navigate('/');
      return;
    }
    //console.log(newValue); 
    if (!!controller) {
      controller.abort();
    }
    setController(new AbortController());    
    const jwtToken = !!GlobalState.jwtToken ? GlobalState.jwtToken : globalJwtTokenState;
    if (newValue === 0 || newValue === 2) {
      fetchSearchLocation(jwtToken);
    } else {
      fetchLastMatches(jwtToken);
    };
  }

  const fetchSearchLocation = async (jwtToken: string) => {
    const result = await fetchGasStations(navigate, jwtToken, controller, globalUserDataState);
    const myResult = result.filter(value => value?.GasPrices?.length > 0).map(value => {
      return value;
    }).map(value => ({
      location: value.Place + ' ' + value.Brand + ' ' + value.Street + ' ' + value.HouseNumber, e5: value.GasPrices[0].E5,
      e10: value.GasPrices[0].E10, diesel: value.GasPrices[0].Diesel, date: new Date(Date.parse(value.GasPrices[0].Date)), longitude: value.Longitude, latitude: value.Latitude
    } as TableDataRow));
    const myPostcode = formatPostCode(globalUserDataState.PostCode);
    const myJson = await fetchPriceAvgs(navigate, jwtToken, controller, myPostcode);
    const rowCounty = ({
      location: myJson.County, e5: Math.round(myJson.CountyAvgE5), e10: Math.round(myJson.CountyAvgE10), diesel: Math.round(myJson.CountyAvgDiesel), date: new Date(), longitude: 0, latitude: 0
    } as TableDataRow);
    const rowState = ({
      location: myJson.State, e5: Math.round(myJson.StateAvgE5), e10: Math.round(myJson.StateAvgE10), diesel: Math.round(myJson.StateAvgDiesel), date: new Date(), longitude: 0, latitude: 0
    } as TableDataRow);
    const resultRows = [rowCounty, rowState, ...myResult]
    setRows(resultRows);
    setGsValues(myResult);
    setController(null);
  }

  const fetchLastMatches = async (jwtToken: string) => {
    const myResult = await fetchUserNotifications(navigate, jwtToken, controller, globalUserUuidState);
    //console.log(myJson);
    const result2 = myResult?.map(value => {
      //console.log(JSON.parse(value?.DataJson));
      return (JSON.parse(value?.DataJson) as MyDataJson[])?.map(value2 => {
        //console.log(value2);
        return {
          location: value2.Place + ' ' + value2.Brand + ' ' + value2.Street + ' ' + value2.HouseNumber,
          e5: value2.E5, e10: value2.E10, diesel: value2.Diesel, date: new Date(Date.parse(value2.Timestamp)), longitude: 0, latitude: 0
        } as TableDataRow;
      });
    })?.flat() || [];
    setRows(result2);
    //const result     
    const myPostcode = formatPostCode(globalUserDataState.PostCode);
    const myJson1 = await fetchTimeSlots(navigate, jwtToken, controller, myPostcode);
    const timeSlots = [] as TimeSlot[];
    timeSlots.push(...myJson1.filter(myValue => myValue.AvgDiesel > 10).map(myValue => {
      const dieselTimeSlot = { x: '00.00', diesel: 0, e10: 0, e5: 0 } as TimeSlot;
      const myDate = new Date(myValue.StartDate);
      dieselTimeSlot.x = '' + myDate.getHours() + ':' + (myDate.getMinutes().toString().length < 2 ? myDate.getMinutes().toString().length + '0' : myDate.getMinutes());
      dieselTimeSlot.diesel = myValue.AvgDiesel / 1000;
      dieselTimeSlot.e10 = myValue.AvgE10 / 1000;
      dieselTimeSlot.e5 = myValue.AvgE5 / 1000;
      return dieselTimeSlot;
    }));
    setAvgTimeSlots(timeSlots);
    //console.log(myJson1);      
    setController(null);
  }

  // eslint-disable-next-line
  useEffect(() => {
    if (globalJwtTokenState?.length > 10 && globalUserUuidState?.length > 10 && first) {
      setTimeout(() => handleTabChange({} as unknown as SyntheticEvent, value), 3000);
      setFirst(false);
    }
  });

  return (<Box sx={{ width: '100%' }}>
    <Tabs value={value} onChange={handleTabChange} centered={true}>
      <Tab label={t('table.currentPrices')} />
      <Tab label={t('table.lastPriceMatches')} />
      <Tab label={t('table.currentPriceMap')} />
    </Tabs>    
    {value === 0 &&
       <DataTable diesel={t('table.diesel')} e10={t('table.e10')} e5={t('table.e5')} location={t('table.location')} showAverages={true} time={t('table.time')} rows={rows}></DataTable>}    
    {value === 1 &&
      <Chart timeSlots={avgTimeSlots}></Chart>}      
    {value === 1 &&
      <DataTable diesel={t('table.diesel')} e10={t('table.e10')} e5={t('table.e5')} location={t('table.location')} showAverages={true} time={t('table.time')} rows={rows}></DataTable>}      
    {value === 2 &&
      <GsMap gsValues={gsValues} center={globalUserDataState}></GsMap>}      
  </Box>);
}
