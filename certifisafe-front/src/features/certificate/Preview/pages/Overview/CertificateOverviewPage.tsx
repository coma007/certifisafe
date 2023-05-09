import Menu from "components/navigation/Menu/Menu"
import PageTitle from "components/view/PageTitle/PageTitle"
import Table from "components/tables/Table/Table"
import { TableRowData } from "components/tables/TableRow/TableRow"
import { formatDate } from "utils/DateUtils"
import CertificateOreviewPageCSS from "./CertificateOverviewPage.module.scss"
import Download from "assets/actions/download.png"
import Withdraw from "assets/actions/withdraw.png"
import ImageButton from "components/tables/ImageButton/ImageButton"
import { useState } from "react"
import Modal from "react-modal";
import ModalWindow from "components/view/Modal/ModalWindow"

const CertificateOreviewPage = () => {

    const [withdrawIsOpen, setWithdrawModalIsOpen] = useState(false);

    const openWithdrawModal = () => {
        setWithdrawModalIsOpen(true);
    };

    const closeWithdrawModal = () => {
        setWithdrawModalIsOpen(false);
    };

    const header: TableRowData = {
        content: "aaa",
        widthPercentage: 20
    }

    const headers: TableRowData[] = [
        { content: "Name", widthPercentage: 28 },
        { content: "Date", widthPercentage: 12 },
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
    { content: <ImageButton path={Download} tooltipText="Download" onClick={null} />, widthPercentage: 0 },
    { content: <ImageButton path={Withdraw} tooltipText="Withdraw" onClick={openWithdrawModal} />, widthPercentage: 0 }]

    const rows: TableRowData[][] = [row, row, row, row, row];

    return (
        <div className={`page pageWithCols ${CertificateOreviewPageCSS.cols}`}>
            <Menu />
            <div>
                <PageTitle title="Certificates overview" description="Take a detailed view of your certificates." />
                <div className={CertificateOreviewPageCSS.table} >
                    <Table headers={headers} rows={rows} />
                </div>
                <ModalWindow height="67%"
                    isOpen={withdrawIsOpen}
                    closeWithdrawalModal={closeWithdrawModal}
                    title="Withdraw Certificate"
                    description="To withdraw the certificate, you need to provide us some more info on why you want to withdraw it. "
                    warning="Please note that if the certificate is revoked, all in the chain below it is automatically retracted. This means that all certificates signed by this certificate, as well as certificates signed by those certificates will be automatically revoked."
                    buttonText="WITHDRAW" />
            </div>
        </div>
    )
}

export default CertificateOreviewPage