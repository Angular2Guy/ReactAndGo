import {Tabs,Tab, Box} from '@mui/material';
import {useState, SyntheticEvent} from 'react';
import DataTable, { TableDataRow } from './DataTable';
import { useRecoilValue } from 'recoil';
import GlobalState from '../GlobalState';
import styles from './main.module.scss';

interface GasStation {  
	StationName:             string;
	Brand:                   string;
	Street:                  string;
	Place:                   string;
	HouseNumber:             string;
	PostCode:                string;
	Latitude:                number;
	Longitude:               number;
	PublicHolidayIdentifier: string;		
	OtJson:                  string;	
	FirstActive:             Date;
	GasPrices:               GasPrice[];
}

interface GasPrice {  
	E5:           number;
	E10:          number;
	Diesel:       number;
	Date:         Date;
	Changed:      number;
}

interface Notification {
  Timestamp:      Date;
	UserUuid:         string;
	Title:            string;
	Message:          string;
	DataJson:         string;
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
    const [value, setValue] = useState(0);
    const [rows, setRows] = useState([] as TableDataRow[]);
    const globalJwtTokenState = useRecoilValue(GlobalState.jwtTokenState);
    const globalUserUuidState = useRecoilValue(GlobalState.userUuidState);
    const globalUserDataState = useRecoilValue(GlobalState.userDataState);

    const handleTabChange = async (event: SyntheticEvent, newValue: number) => {
        setValue(newValue);
        console.log(newValue);        
        const requestOptions1 = {
          method: 'GET',
          headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${globalJwtTokenState}`},
        }
        const requestOptions2 = {
          method: 'POST',
          headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${globalJwtTokenState}`},
          body: JSON.stringify({Longitude: globalUserDataState.Longitude, Latitude: globalUserDataState.Latitude, Radius: globalUserDataState.SearchRadius})
        }
        if(newValue === 0) {
          const myResult = await fetch('/gasstation/search/location', requestOptions2);
          const myJson = await myResult.json() as GasStation[];  
          setRows(myJson.filter(value => value?.GasPrices?.length > 0).map(value => ({ location: value.Place + ' ' + value.Brand + ' ' + value.Street + ' ' + value.HouseNumber, e5: value.GasPrices[0].E5, 
            e10: value.GasPrices[0].E10, diesel: value.GasPrices[0].Diesel, date: value.GasPrices[0].Date } as TableDataRow)));
        } else {
          const myResult = await fetch(`/usernotification/current/${globalUserUuidState}`, requestOptions1);
          const myJson = await myResult.json() as Notification[];
          console.log(myJson);
          /*
          setRows(myJson.map(value => ({ location: value.Message, e5: value?.GasPrices?.length > 0 ? value.GasPrices[0].E5 : 0, 
            e10: value?.GasPrices?.length > 0 ? value.GasPrices[0].E10 : 0, diesel: value?.GasPrices?.length > 0 ? value.GasPrices[0].Diesel : 0, date: value?.GasPrices?.length > 0 ? value.GasPrices[0].Date : new Date() } as TableDataRow)));
            */
           setRows([]);
        }        
    };    

    return (<Box sx={{ width: '100%' }}>
        <Tabs value={value} onChange={handleTabChange} >
            <Tab label="Current Prices"/>
            <Tab label="Last Price changes"/>
        </Tabs>
        <TabPanel value={value} index={0}>
        <DataTable diesel='Diesel' e10='E10' e5='E5' location='Location' rows={rows}></DataTable>
      </TabPanel>
      <TabPanel value={value} index={1}>
      <DataTable diesel='Diesel' e10='E10' e5='E5' location='Location' rows={rows}></DataTable>
      </TabPanel>
    </Box>);
}