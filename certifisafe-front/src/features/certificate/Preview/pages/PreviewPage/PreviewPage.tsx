import Menu from "components/navigation/Menu/Menu"
import PageTitle from "components/view/PageTitle/PageTitle"
import Table from "components/view/Table/Table"
import { TableRowData } from "components/view/TableRow/TableRow"
import { formatDate } from "utils/DateUtils"
import PreviewPageCSS from "./PreviewPage.module.scss"
import Download from "assets/actions/download.png"
import Withdraw from "assets/actions/withdraw.png"
import ImageButton from "components/tables/ImageButton/ImageButton"
import Tooltip from "components/view/Tooltip/Tooltip"

const PreviewPage = () => {

    const header: TableRowData = {
        content: "aaa",
        widthPercentage: 20
    }

    const headers: TableRowData[] = [
        { content: "Name", widthPercentage: 30 },
        { content: "Date", widthPercentage: 10 },
        { content: "Subject", widthPercentage: 15 },
        { content: "Issuer", widthPercentage: 15 },
        { content: "Type", widthPercentage: 10 },
        { content: "Status", widthPercentage: 10 },
        { content: "", widthPercentage: 5 },
        { content: "", widthPercentage: 5 }]

    const row: TableRowData[] = [{ content: "My certificate 1", widthPercentage: 0 },
    { content: formatDate(new Date(Date.now())), widthPercentage: 0 },
    { content: "UNS", widthPercentage: 0 },
    { content: "Google Inc.", widthPercentage: 0 },
    { content: "root", widthPercentage: 0 },
    { content: <i>ACTIVE</i>, widthPercentage: 0 },
    { content: <ImageButton path={Download} tooltipText="Download" />, widthPercentage: 0 },
    { content: <ImageButton path={Withdraw} tooltipText="Withdraw" />, widthPercentage: 0 }]

    const rows: TableRowData[][] = [row, row, row, row, row];

    return (
        <div className={`page pageWithCols ${PreviewPageCSS.cols}`}>
            <Menu />
            <div>
                <PageTitle title="Certificates overview" description="Take a detailed view of your certificates." />
                <div className={PreviewPageCSS.table} >
                    <Table headers={headers} rows={rows} />
                </div>
            </div>
        </div>
    )
}

export default PreviewPage