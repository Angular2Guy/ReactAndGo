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
import { useMemo } from 'react';
import styles from "./datatable.module.scss";

export interface TableDataRow {
    location: string;
    e5: number;
    e10: number;
    diesel: number;
    date: Date;
    longitude: number;
    latitude: number;
    e5Class: string;
    e10Class: string;
    dieselClass: string;
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

  useMemo(() => {   
    if(props?.rows?.length < 4) {
      return;
    }  
    const e5Arr = [...props.rows].filter(row => row.e5 > 10);
    e5Arr.sort((a,b) => a.e5 - b.e5);
    const e10Arr = [...props.rows].filter(row => row.e10 > 10);    
    e10Arr.sort((a,b) => a.e10 - b.e10);
    const dieselArr = [...props.rows].filter(row => row.diesel > 10);
    dieselArr.sort((a,b) => a.diesel - b.diesel);
    if(e5Arr?.length >= 3) {
    e5Arr[0].e5Class = 'best-price';
    e5Arr[1].e5Class = 'good-price';
    e5Arr[2].e5Class = 'good-price';
    }
    if(e10Arr?.length >= 3) {
    e10Arr[0].e10Class = 'best-price';
    e10Arr[1].e10Class = 'good-price';
    e10Arr[2].e10Class = 'good-price';
    }
    if(dieselArr?.length >= 3) {
    dieselArr[0].dieselClass = 'best-price';
    dieselArr[1].dieselClass = 'good-price';
    dieselArr[2].dieselClass = 'good-price';      
    }
  },[props.rows]);

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
                  <TableCell align="right" className={styles[row.e5Class]}>{row.e5}</TableCell>
                  <TableCell align="right" className={styles[row.e10Class]}>{row.e10}</TableCell>
                  <TableCell align="right" className={styles[row.dieselClass]}>{row.diesel}</TableCell>                  
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      );
}