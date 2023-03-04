import {TableContainer, Paper, Table, TableHead, TableRow, TableCell, TableBody} from '@mui/material';

export interface TableDataRow {
    location: string;
    e5: number;
    e10: number;
    diesel: number;
}

interface DataTableProps {
    location: string;
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
                <TableCell align="right">{props.e5}</TableCell>
                <TableCell align="right">{props.e10}</TableCell>
                <TableCell align="right">{props.diesel}</TableCell>                
              </TableRow>
            </TableHead>
            <TableBody>
              {props.rows.map((row,index) => (
                <TableRow
                  key={index}
                  sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                >
                  <TableCell component="th" scope="row">
                    {row.location}
                  </TableCell>
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