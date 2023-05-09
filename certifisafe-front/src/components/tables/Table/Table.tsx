import { useState } from "react";
import TableRow, { TableRowData } from "../TableRow/TableRow"
import TableCSS from "./Table.module.scss"

// TODO make table resizable
const Table = (props: { headers: TableRowData[], rows: TableRowData[][] }) => {
    const [columnWidths, setColumnWidths] = useState<number[]>(props.headers.map(header => header.widthPercentage));

    const rows = props.rows.map((row) =>
        row.map((cell, columnIndex) => ({
            content: cell.content,
            widthPercentage: columnWidths[columnIndex],
        }))
    );

    return (
        <div>
            <TableRow className={TableCSS.header} key="headers" data={props.headers} />
            {rows.map((row, index) => (
                <TableRow className="" key={index} data={row} />
            ))}
        </div>
    )
}

export default Table;