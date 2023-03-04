import {Tabs,Tab, Box} from '@mui/material';
import {useState, SyntheticEvent} from 'react';
import DataTable, { TableDataRow } from './DataTable';
import styles from './main.module.scss';

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
    const handleTabChange = (event: SyntheticEvent, newValue: number) => {
        setValue(newValue);
        console.log(newValue);
    };
    const rows1 = [{
      diesel: 1.75,
      e10: 1.80,
      e5: 1.85,
      location: 'Bremen'
    } as TableDataRow];
    const rows2 = [{
      diesel: 1.70,
      e10: 1.75,
      e5: 1.80,
      location: 'Berlin'
    } as TableDataRow];

    return (<Box sx={{ width: '100%' }}>
        <Tabs value={value} onChange={handleTabChange} >
            <Tab label="Current Prices"/>
            <Tab label="Last Price changes"/>
        </Tabs>
        <TabPanel value={value} index={0}>
        <DataTable diesel='Diesel' e10='E10' e5='E5' location='Location' rows={rows1}></DataTable>
      </TabPanel>
      <TabPanel value={value} index={1}>
      <DataTable diesel='Diesel' e10='E10' e5='E5' location='Location' rows={rows2}></DataTable>
      </TabPanel>
    </Box>);
}