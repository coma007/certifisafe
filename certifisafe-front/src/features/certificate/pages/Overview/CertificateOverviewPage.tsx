import Menu from "components/navigation/Menu/Menu"
import PageTitle from "components/view/PageTitle/PageTitle"
import Table from "components/tables/Table/Table"
import { TableRowData } from "components/tables/TableRow/TableRow"
import { formatDate } from "utils/DateUtils"
import CertificateOreviewPageCSS from "./CertificateOverviewPage.module.scss"
import Download from "assets/actions/download.png"
import Withdraw from "assets/actions/withdraw.png"
import ImageButton from "components/tables/ImageButton/ImageButton"
import { useEffect, useState } from "react"
import ModalWindow from "components/view/Modal/ModalWindow"
import ModalWindowCSS from "components/view/Modal/ModalWindow.module.scss"
import { CertificateService } from "features/certificate/services/CertificateService"
import { Certificate } from "features/certificate/types/Certificate"

const CertificateOreviewPage = () => {

    const [withdrawIsOpen, setWithdrawModalIsOpen] = useState(false);
    const [selectedCertificate, setSelectedCertificate] = useState<Certificate|undefined>(undefined);
    const [tableData, setTableData] = useState<TableRowData[][]>([]);

    useEffect(() => {
        (async function () {
            try {
                const fetchedCertificates = await CertificateService.getAll();
                populateData(fetchedCertificates);
            } catch (error) {
                console.error(error);
            }
        })()
    }, []);

    const openWithdrawModal = (certificate : Certificate) => {
        setSelectedCertificate(certificate)
        setWithdrawModalIsOpen(true);
    };

    const closeWithdrawModal = () => {
        setWithdrawModalIsOpen(false);
    };

    const okWithdrawModal = () => {
        CertificateService.withdraw(selectedCertificate!.Serial)
        setWithdrawModalIsOpen(false);
    };

    const headers: TableRowData[] = [
        { content: "Name", widthPercentage: 28 },
        { content: "Date", widthPercentage: 12 },
        { content: "Subject", widthPercentage: 15 },
        { content: "Issuer", widthPercentage: 15 },
        { content: "Type", widthPercentage: 10 },
        { content: "Status", widthPercentage: 10 },
        { content: "", widthPercentage: 5 },
        { content: "", widthPercentage: 5 }
    ]


    const populateData = (certificates: Certificate[]) => {
        let data: TableRowData[][] = []
        if (certificates !== undefined) {
            certificates.forEach(certificate => {
                data.push([
                    { content: certificate.Name, widthPercentage: 28 },
                    { content: formatDate(new Date(certificate.ValidFrom)), widthPercentage: 12 },
                    { content: certificate.Subject.FirstName, widthPercentage: 15 },
                    { content: certificate.Issuer.FirstName, widthPercentage: 15 },
                    { content: certificate.Type.toLowerCase(), widthPercentage: 10 },
                    { content: <i> {certificate.Status}</i>, widthPercentage: 10 },
                    {
                        content: <ImageButton path={Download} tooltipText="Download" onClick={() => {
                            CertificateService.download(certificate.Serial);
                        }} />, widthPercentage: 5
                    },
                    { content: <ImageButton path={Withdraw} tooltipText="Withdraw" onClick={() => openWithdrawModal(certificate)} />, widthPercentage: 5 }]
                );
            });
        }
        setTableData(data);
    }


    return (
        <div className={`page pageWithCols pageWithMenu`}>
            <Menu />
            <div>
                <PageTitle title="Certificates overview" description="Take a detailed view of your certificates." />
                <div className={CertificateOreviewPageCSS.table} >
                    <Table headers={headers} rows={tableData} />
                </div>
                <ModalWindow
                    height="75%"
                    isOpen={withdrawIsOpen}
                    closeWithdrawalModal={closeWithdrawModal}
                    okWithdrawalModal={okWithdrawModal}
                    title="Withdraw certificate"
                    buttonText="WITHDRAW" >
                    <p>To withdraw the certificate, you need to provide us some more info on why you want to withdraw it. </p>
                    <textarea placeholder='Write your reason ...'></textarea>
                    <p className={ModalWindowCSS.warning}>Please note that if the certificate is revoked, all in the chain below it is automatically retracted. This means that all certificates signed by this certificate, as well as certificates signed by those certificates will be automatically revoked.</p>
                </ModalWindow>
            </div>
        </div>
    )
}

export default CertificateOreviewPage