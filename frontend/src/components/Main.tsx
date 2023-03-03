import {Tabs,Tab, Box} from '@mui/material';
import {useState, SyntheticEvent} from 'react';
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

    return (<Box sx={{ width: '100%' }}>
        <Tabs value={value} onChange={handleTabChange} >
            <Tab label="Item One"/>
            <Tab label="Item Two"/>
        </Tabs>
        <TabPanel value={value} index={0}>
        Item One
      </TabPanel>
      <TabPanel value={value} index={1}>
        Item Two
      </TabPanel>
    </Box>);
}