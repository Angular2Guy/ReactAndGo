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
import {TableContainer, Paper, Table, TableHead, TableRow, TableCell, TableBody} from '@mui/material';
import { nanoid } from 'nanoid';

export interface TableDataRow {
    location: string;
    e5: number;
    e10: number;
    diesel: number;
    date: Date;
    longitude: number;
    latitude: number;
}

interface DataTableProps {
    location: string;
    time: string;
    e5: string;
    e10: string;
    diesel: string;
    rows: TableDataRow[];
}

export default function DataTable(props: DataTableProps) {    
    return (
        <TableContainer component={Paper}>
          <Table sx={{ minWidth: '100%' }}>
            <TableHead>
              <TableRow>
                <TableCell>{props.location}</TableCell>
                <TableCell>{props.time}</TableCell>
                <TableCell align="right">{props.e5}</TableCell>
                <TableCell align="right">{props.e10}</TableCell>
                <TableCell align="right">{props.diesel}</TableCell>                
              </TableRow>
            </TableHead>
            <TableBody>
              {props.rows.map((row) => (
                <TableRow
                  key={nanoid()}
                  sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                >
                  <TableCell component="th" scope="row">
                    {row.location}
                  </TableCell>
                  <TableCell>{row.date.toISOString().split('T')[0]+' '+row.date.toTimeString().split(' ')[0]}</TableCell>                  
                  <TableCell align="right">{row.e5}</TableCell>
                  <TableCell align="right">{row.e10}</TableCell>
                  <TableCell align="right">{row.diesel}</TableCell>                  
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      );
}