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
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { useEffect, useState } from 'react';
import Radio from '@mui/material/Radio';
import RadioGroup from '@mui/material/RadioGroup';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormControl from '@mui/material/FormControl';
import GlobalState, { FuelType } from '../GlobalState';
import { type TimeSlot } from '../model/time-slot-response';
import { type GsPoint } from '../model/gs-point';
import { useAtom } from "jotai";

export interface ChartProps {
  timeSlots: TimeSlot[];
}

export default function Chart(props: ChartProps) {
  //console.log(props.timeSlots);  
  const [gsValues, setGsValues] = useState([] as GsPoint[]);  
  const [fuelTypeState, setfuelTypeState] = useAtom(GlobalState.fuelTypeState);
  const [lineColor, setLineColor] = useState('#8884d8');
  const [avgValue, setAvgValue] = useState(0.0);
  const [timeSlots, setTimeSlots] = useState([] as TimeSlot[]);

  // eslint-disable-next-line
  useEffect(() => {
    if(props.timeSlots.length > 0 && timeSlots !== props.timeSlots) {    
      updateChart();            
      setTimeSlots(props.timeSlots);
    }    
  });

  function updateChart() {
    const avg = props.timeSlots.reduce((acc, value) => value[fuelTypeState] + acc, 0) / (props.timeSlots.length || 1);
    if (fuelTypeState === FuelType.E5) {
      setLineColor('#8884d8');
      setAvgValue(avg);
      setGsValues(props.timeSlots.map(myValue => ({ timestamp: myValue.x, price: myValue.e5 - avg } as GsPoint)))
    } else if (fuelTypeState === FuelType.E10) {
      setLineColor('#82ca9d');
      setAvgValue(avg);
      setGsValues(props.timeSlots.map(myValue => ({ timestamp: myValue.x, price: myValue.e10 - avg } as GsPoint)))
    } else {
      setLineColor('#82caff');
      setAvgValue(avg);
      setGsValues(props.timeSlots.map(myValue => ({ timestamp: myValue.x, price: myValue.diesel - avg } as GsPoint)))
    }
  }

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {    
    setfuelTypeState(((event.target as HTMLInputElement).value) as FuelType);
    setTimeSlots([]);
    updateChart();    
  };
  return (<div>
    <ResponsiveContainer width="100%" height={300}>
      <LineChart data={gsValues}
        margin={{ top: 20, right: 20, left: 20, bottom: 20 }}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="timestamp" />
        <YAxis />
        <Tooltip />
        <Legend />
        <Line name={'AvgPrice: '+(Math.round(avgValue*1000)/1000)} type="monotone" dataKey="price" stroke={lineColor} />
      </LineChart>
    </ResponsiveContainer>
    <FormControl>
      <RadioGroup
        row
        aria-labelledby="demo-row-radio-buttons-group-label"
        name="row-radio-buttons-group"
        value={fuelTypeState}
        onChange={handleChange}
      >
        <FormControlLabel value={FuelType.E5} control={<Radio />} label="E5" />
        <FormControlLabel value={FuelType.E10} control={<Radio />} label="E10" />
        <FormControlLabel value={FuelType.Diesel} control={<Radio />} label="Diesel" />
      </RadioGroup>
    </FormControl>
  </div>
  );
}
