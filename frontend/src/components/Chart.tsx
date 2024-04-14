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

export enum FuelType {
    e5, e10, diesel
}

export interface TimeSlot {
    x: Date;
    y: number;
    fuelType: FuelType;
}

export interface ChartProps {
    e5: TimeSlot[];
    e10: TimeSlot[];
    diesel: TimeSlot[];
}

export default function Chart(props: ChartProps) {
    
        return ( (props?.diesel?.length > 0 || props?.e10?.length > 0 || props?.e5?.length > 0) ? 
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
              <XAxis type="number" dataKey="x" name="time" />
              <YAxis type="number" dataKey="y" name="price" />
              <ZAxis type="number" range={[100]} />
              <Tooltip cursor={{ strokeDasharray: '3 3' }} />
              <Legend />
              <Scatter name="E5" data={props.e5} fill="#8884d8" line shape="cross" />
              <Scatter name="E10" data={props.e10} fill="#82ca9d" line shape="diamond" />
              <Scatter name="Diesel" data={props.diesel} fill="#ff8042" line shape="triangle" />
            </ScatterChart>
          </ResponsiveContainer>
          : <div></div>
        );
}