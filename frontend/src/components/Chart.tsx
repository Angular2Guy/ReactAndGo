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
import {
  ScatterChart,
  Scatter,
  XAxis,
  YAxis,
  ZAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import { useEffect, useState, SyntheticEvent } from 'react';

export interface TimeSlot {
    x: string;
    e5: number;
    e10: number;
    diesel: number;
}

export interface ChartProps {
    timeSlots: TimeSlot[];
}

const values = [{diesel: 1700, e10: 1750, e5: 1800, x: '19:00'} as TimeSlot, {diesel: 1710, e10: 1760, e5: 1810, x: '19:30'} as TimeSlot]

export default function Chart(props: ChartProps) {
    const [avgTimeSlots, setAvgTimeSlots] = useState([] as TimeSlot[])
    if(props.timeSlots.length !== avgTimeSlots.length) {
        setAvgTimeSlots(props.timeSlots);
    }    
    //console.log(avgTimeSlots);
    //console.log(avgTimeSlots.filter(value => value.e5 > 10).map(value => value.e5));
    console.log(props.timeSlots);
        return ( props.timeSlots.length > 0 ? 
          <ResponsiveContainer width="100%" height={300}>
            <ScatterChart
              margin={{
                top: 20,
                right: 20,
                bottom: 20,
                left: 20,
              }}
            >
              <CartesianGrid />
              <XAxis type="category" dataKey="x" name="time" />
              <YAxis type="number" dataKey="e5" name="price" />
              <ZAxis type="number" range={[100]} />
              <Tooltip cursor={{ strokeDasharray: '3 3' }} />
              <Legend />
              <Scatter name="E5" data={avgTimeSlots.filter(value => value.e5 > 10).map(value => value.e5)} fill="#8884d8" line shape="cross" />
              <Scatter name="E10" data={avgTimeSlots.filter(value => value.e10 > 10).map(value => value.e10)} fill="#82ca9d" line shape="diamond" />
              <Scatter name="Diesel" data={avgTimeSlots.filter(value => value.diesel > 10).map(value => value.diesel)} fill="#ff8042" line shape="triangle" />
            </ScatterChart>
          </ResponsiveContainer>
          : <div></div>
        );
}
