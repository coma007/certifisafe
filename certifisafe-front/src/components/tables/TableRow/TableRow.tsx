import { Resizable } from "react-resizable";
import TableRowCSS from "./TableRow.module.scss"

export interface TableRowData {
    content: any;
    widthPercentage: number;
}

const TableRow = (props: { data: TableRowData[], className: string }) => {
    return (
        <div className={`${props.className} ${TableRowCSS.row}`}>
            {props.data.map((rowData, index) => (
                <span key={index} style={{ width: `${rowData.widthPercentage}%` }}>
                    {rowData.content}
                </span>
            ))}
        </div>
    )
}

export default TableRow