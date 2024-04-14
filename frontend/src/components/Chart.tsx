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
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

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
    return (
        <ResponsiveContainer width="100%" height={300}>
          <BarChart
            margin={{
              top: 20,
              right: 20,
              bottom: 20,
              left: 20,
            }}
            data={props.timeSlots}
          >
            <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="x" />
          <YAxis  />
          <Tooltip />
          <Legend />
          <Bar dataKey="e5" fill="#8884d8" />
          <Bar dataKey="e10" fill="#82ca9d" />
          <Bar dataKey="diesel" fill="#82caff" />
          </BarChart>
        </ResponsiveContainer>
      );
}
