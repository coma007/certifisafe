import Menu from "components/navigation/Menu/Menu"
import PageTitle from "components/view/PageTitle/PageTitle"
import Table from "components/view/Table/Table"
import { TableRowData } from "components/view/TableRow/TableRow"
import { formatDate } from "utils/DateUtils"
import RequestOreviewPageCSS from "./RequestOverviewPage.module.scss"
import Accept from "assets/actions/accept.png"
import Decline from "assets/actions/decline.png"
import Remove from "assets/actions/withdraw.png"
import ImageButton from "components/tables/ImageButton/ImageButton"
import { SetStateAction, useState } from "react"

const RequestOreviewPage = () => {
    const [selectedOption, setSelectedOption] = useState("foryou");

    const handleOptionChange = (event: { target: { id: SetStateAction<string> } }) => {
        setSelectedOption(event.target.id);
    };


    const header: TableRowData = {
        content: "aaa",
        widthPercentage: 20
    }

    const headersMe: TableRowData[] = [
        { content: "Name", widthPercentage: 35 },
        { content: "Date", widthPercentage: 12 },
        { content: "Subject", widthPercentage: 25 },
        { content: "Type", widthPercentage: 18 },
        { content: "", widthPercentage: 5 },
        { content: "", widthPercentage: 5 }]

    const rowMe: TableRowData[] = [{ content: "My certificate 1", widthPercentage: 0 },
    { content: formatDate(new Date(Date.now())), widthPercentage: 0 },
    { content: "UNS", widthPercentage: 0 },
    { content: "intermediate", widthPercentage: 0 },
    { content: <ImageButton path={Accept} tooltipText="Accept" />, widthPercentage: 0 },
    { content: <ImageButton path={Decline} tooltipText="Decline" />, widthPercentage: 0 }]

    const rowsMe: TableRowData[][] = [rowMe, rowMe, rowMe, rowMe, rowMe];

    const headersMy: TableRowData[] = [
        { content: "Name", widthPercentage: 35 },
        { content: "Date", widthPercentage: 12 },
        { content: "Issuer", widthPercentage: 25 },
        { content: "Type", widthPercentage: 13 },
        { content: "Status", widthPercentage: 10 },
        { content: "", widthPercentage: 5 }]

    const rowMy: TableRowData[] = [
        { content: "My certificate 1", widthPercentage: 0 },
        { content: formatDate(new Date(Date.now())), widthPercentage: 0 },
        { content: "UNS", widthPercentage: 0 },
        { content: "intermediate", widthPercentage: 0 },
        { content: <i>PENDING</i>, widthPercentage: 0 },
        { content: <ImageButton path={Remove} tooltipText="Remove" />, widthPercentage: 0 }]

    const rowsMy: TableRowData[][] = [rowMy, rowMy, rowMy, rowMy, rowMy];


    return (
        <div className={`page pageWithCols ${RequestOreviewPageCSS.cols}`}>
            <Menu />
            <div>
                <PageTitle title="Requests overview" description="Take a detailed view of requests made for you." />


                <div className={RequestOreviewPageCSS.table} >
                    <div className={RequestOreviewPageCSS.radioContainer}>
                        <input
                            type="radio"
                            name="radioGroup"
                            id="foryou"
                            checked={selectedOption === "foryou"}
                            onChange={handleOptionChange}
                        />
                        <label htmlFor="foryou">REQUESTS FOR ME</label>
                        <input
                            type="radio"
                            name="radioGroup"
                            id="fromyou"
                            checked={selectedOption === "fromyou"}
                            onChange={handleOptionChange}
                        />
                        <label htmlFor="fromyou">MY REQUESTS</label>
                    </div>
                    {selectedOption === "foryou" ? (
                        <Table headers={headersMe} rows={rowsMe} />
                    ) : (
                        <Table headers={headersMy} rows={rowsMy} />
                    )}
                </div>
            </div>
        </div>
    )
}

export default RequestOreviewPage