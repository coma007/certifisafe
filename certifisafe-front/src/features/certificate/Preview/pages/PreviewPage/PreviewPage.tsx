import Menu from "components/navigation/Menu/Menu"
import PageTitle from "components/view/PageTitle/PageTitle"
import Table from "components/view/Table/Table"
import { TableRowData } from "components/view/TableRow/TableRow"
import PreviewPageCSS from "./PreviewPage.module.scss"

const PreviewPage = () => {

    const header: TableRowData = {
        content: "aaa",
        widthPercentage: 20
    }

    const headers: TableRowData[] = [
        { content: "Name", widthPercentage: 25 },
        { content: "Date", widthPercentage: 15 },
        { content: "Subject", widthPercentage: 15 },
        { content: "Issuer", widthPercentage: 15 },
        { content: "Type", widthPercentage: 10 },
        { content: "Status", widthPercentage: 10 },
        { content: "", widthPercentage: 5 },
        { content: "", widthPercentage: 5 }]

    const row: TableRowData = {
        content: "asd",
        widthPercentage: 20
    }

    const rows: TableRowData[][] = [[row, row, row, row, row], [row, row, row, row, row]]

    return (
        <div className={`page pageWithCols ${PreviewPageCSS.cols}`}>
            <Menu />
            <div>
                <PageTitle title="Certificates preview" description="Take a detailed view of your certificates." />
                <div className={PreviewPageCSS.table} >
                    <Table headers={headers} rows={rows} />
                </div>
            </div>
        </div>
    )
}

export default PreviewPage