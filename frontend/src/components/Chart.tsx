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
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { useEffect, useState } from 'react';
import Radio from '@mui/material/Radio';
import RadioGroup from '@mui/material/RadioGroup';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormControl from '@mui/material/FormControl';

enum FuelType {
  E5 = 'e5',
  E10 = 'e10',
  Diesel = 'diesel'
}

export interface GsPoint {
  timestamp: string;
  price: number;
}

export interface TimeSlot {
  x: string;
  e5: number;
  e10: number;
  diesel: number;
}

export interface ChartProps {
  timeSlots: TimeSlot[];
}

export default function Chart(props: ChartProps) {
  //console.log(props.timeSlots);
  const [gsValues, setGsValues] = useState([] as GsPoint[]);
  const [fuelType, setFuelType] = useState(FuelType.E5);
  const [lineColor, setLineColor] = useState('#8884d8')

  // eslint-disable-next-line
  useEffect(() => {
    if (fuelType === FuelType.E5) {
      setLineColor('#8884d8')
      const avgValue = props.timeSlots.reduce((acc, value) => value.e5 + acc, 0) / (props.timeSlots.length || 1);
      setGsValues(props.timeSlots.map(myValue => ({ timestamp: myValue.x, price: myValue.e5 - avgValue } as GsPoint)))
    } else if (fuelType === FuelType.E10) {
      setLineColor('#82ca9d')
      const avgValue = props.timeSlots.reduce((acc, value) => value.e10 + acc, 0) / (props.timeSlots.length || 1);
      setGsValues(props.timeSlots.map(myValue => ({ timestamp: myValue.x, price: myValue.e10 - avgValue } as GsPoint)))
    } else {
      setLineColor('#82caff')
      const avgValue = props.timeSlots.reduce((acc, value) => value.diesel + acc, 0) / (props.timeSlots.length || 1);
      setGsValues(props.timeSlots.map(myValue => ({ timestamp: myValue.x, price: myValue.diesel - avgValue } as GsPoint)))
    }
  });

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setFuelType(((event.target as HTMLInputElement).value) as FuelType);
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
        <Line type="monotone" dataKey="price" stroke={lineColor} />
      </LineChart>
    </ResponsiveContainer>
    <FormControl>
      <RadioGroup
        row
        aria-labelledby="demo-row-radio-buttons-group-label"
        name="row-radio-buttons-group"
        value={fuelType}
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
