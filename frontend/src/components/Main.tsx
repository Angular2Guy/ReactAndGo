import { Tabs, Tab, Box } from '@mui/material';
import { useEffect, useState, SyntheticEvent } from 'react';
import DataTable, { TableDataRow } from './DataTable';
import { useRecoilValue } from 'recoil';
import GlobalState from '../GlobalState';
import styles from './main.module.scss';

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
  let controller: AbortController | null = null;
  const [value, setValue] = useState(0);
  const [first, setFirst] = useState(true);
  const [rows, setRows] = useState([] as TableDataRow[]);
  const globalJwtTokenState = useRecoilValue(GlobalState.jwtTokenState);
  const globalUserUuidState = useRecoilValue(GlobalState.userUuidState);
  const globalUserDataState = useRecoilValue(GlobalState.userDataState);  


  const handleTabChange = (event: SyntheticEvent, newValue: number) => {
    setValue(newValue);
    if (globalJwtTokenState?.length < 10 || globalUserUuidState?.length < 10) {
      return;
    }    
    //console.log(newValue);        
    if(!!controller) {
      controller.abort();
    }
    controller = new AbortController();
    const requestOptions1 = {
      method: 'GET',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${globalJwtTokenState}` },
      signal: controller.signal
    }
    const requestOptions2 = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${globalJwtTokenState}` },
      body: JSON.stringify({ Longitude: globalUserDataState.Longitude, Latitude: globalUserDataState.Latitude, Radius: globalUserDataState.SearchRadius }),
      signal: controller.signal
    }
    if (newValue === 0) {      
      fetch('/gasstation/search/location', requestOptions2).then(myResult => myResult.json() as Promise<GasStation[]>).then(myJson => setRows(myJson.filter(value => value?.GasPrices?.length > 0).map(value => ({
        location: value.Place + ' ' + value.Brand + ' ' + value.Street + ' ' + value.HouseNumber, e5: value.GasPrices[0].E5,
        e10: value.GasPrices[0].E10, diesel: value.GasPrices[0].Diesel, date: new Date(Date.parse(value.GasPrices[0].Date))
      } as TableDataRow)))).then(() => controller = null);
    } else {     
      fetch(`/usernotification/current/${globalUserUuidState}`, requestOptions1).then(myResult => myResult.json() as Promise<Notification[]>).then(myJson => {
        //console.log(myJson);
        const result = myJson.map(value => {
          //console.log(JSON.parse(value?.DataJson));
          return (JSON.parse(value?.DataJson) as MyDataJson[])?.map(value2 => {
            //console.log(value2);
            return {
              location: value2.Place + ' ' + value2.Brand + ' ' + value2.Street + ' ' + value2.HouseNumber,
              e5: value2.E5, e10: value2.E10, diesel: value2.Diesel, date: new Date(Date.parse(value2.Timestamp))
            } as TableDataRow;
          });
        })?.flat();
        setRows(result);
      }).then(() => controller = null);
    };
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
      <Tab label="Last Price changes" />
    </Tabs>
    <TabPanel value={value} index={0}>
      <DataTable diesel='Diesel' e10='E10' e5='E5' location='Location' time='Time' rows={rows}></DataTable>
    </TabPanel>
    <TabPanel value={value} index={1}>
      <DataTable diesel='Diesel' e10='E10' e5='E5' location='Location' time='Time' rows={rows}></DataTable>
    </TabPanel>
  </Box>);
}