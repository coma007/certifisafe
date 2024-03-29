import Menu from "components/navigation/Menu/Menu"
import PageTitle from "components/view/PageTitle/PageTitle"
import Table from "components/tables/Table/Table"
import { TableRowData } from "components/tables/TableRow/TableRow"
import { formatDate } from "utils/DateUtils"
import Accept from "assets/actions/accept.png"
import Decline from "assets/actions/decline.png"
import Remove from "assets/actions/withdraw.png"
import ImageButton from "components/tables/ImageButton/ImageButton"
import { SetStateAction, useEffect, useState } from "react"
import ModalWindow from "components/view/Modal/ModalWindow"
import { Request } from "features/request/types/Request"
import { RequestService } from "features/request/service/RequestService"
import RequestOverviewPageCSS from "./RequestOverviewPage.module.scss"
import { CertificateService } from "features/certificate/services/CertificateService"

const RequestOverviewPage = () => {
    const [selectedOption, setSelectedOption] = useState("foryou");
    const [selectedRequest, setSelectedRequest] = useState<Request | undefined>(undefined);
    const [declineReason, setDeclineReason] = useState<string | undefined>(undefined);
    const handleOptionChange = (event: { target: { id: SetStateAction<string> } }) => {
        setSelectedOption(event.target.id);
    };

    const [declineIsOpen, setDeclineModalIsOpen] = useState(false);
    const openDeclineModal = (request: Request, reason: string) => {
        setDeclineReason(declineReason)
        setSelectedRequest(request)
        setDeclineModalIsOpen(true);
    };

    const closeDeclineModal = () => {
        setDeclineModalIsOpen(false);
    };

    const okDeclineModal = async () => {
        await RequestService.decline(selectedRequest!.ID, declineReason!)
        setDeclineModalIsOpen(false);
    };

    const handleReasonChange = (event: any) => {
        setDeclineReason(event.target.value);
    };


    const [rowsMe, setTableDataMe] = useState<TableRowData[][]>([]);
    const [rowsMy, setTableDataMy] = useState<TableRowData[][]>([]);

    useEffect(() => {
        (async function () {
            try {
                const fetchedRequests = await RequestService.getAllByUserSigning();
                populateMeData(fetchedRequests);
            } catch (error) {
                console.error(error);
            }
        })()
    }, []);

    useEffect(() => {
        (async function () {
            try {
                const fetchedRequests = await RequestService.getByUser();
                populateMyData(fetchedRequests);
            } catch (error) {
                console.error(error);
            }
        })()
    }, []);


    const populateMeData = (requests: Request[]) => {
        let data: TableRowData[][] = []
        if (requests !== undefined && requests !== null) {
            requests.forEach(request => {
                data.push([
                    { content: request.CertificateName, widthPercentage: 35 },
                    { content: formatDate(new Date(request.Date)), widthPercentage: 12 },
                    { content: request.Subject.FirstName, widthPercentage: 25 },
                    { content: request.CertificateType, widthPercentage: 13 },
                    { content: request.Status.toLowerCase() === "pending" ? <ImageButton path={Accept} tooltipText="Accept" onClick={async () => { await RequestService.accept(request.ID); window.location.reload() }} /> : null, widthPercentage: 10 },
                    { content: request.Status.toLowerCase() === "pending" ? <ImageButton path={Decline} tooltipText="Decline" onClick={() => openDeclineModal(request, "")} /> : null, widthPercentage: 5 },
                ]);
            });
        }
        setTableDataMe(data);
    }


    const populateMyData = (request: Request[]) => {
        let data: TableRowData[][] = []
        if (request !== undefined && request !== null) {
            request.forEach(request => {
                data.push([
                    { content: request.CertificateName, widthPercentage: 35 },
                    { content: formatDate(new Date(request.Date)), widthPercentage: 12 },
                    { content: request.Subject.FirstName, widthPercentage: 25 },
                    { content: request.CertificateType, widthPercentage: 13 },
                    { content: request.Status, widthPercentage: 10 },
                    { content: <ImageButton path={Remove} tooltipText="Remove" onClick={() => RequestService.delete(request.ID)} />, widthPercentage: 5 }
                ]);
            });
        }
        setTableDataMy(data);
    }


    const headersMe: TableRowData[] = [
        { content: "Name", widthPercentage: 35 },
        { content: "Date", widthPercentage: 12 },
        { content: "Subject", widthPercentage: 25 },
        { content: "Type", widthPercentage: 13 },
        { content: "", widthPercentage: 10 },
        { content: "", widthPercentage: 5 }]


    const headersMy: TableRowData[] = [
        { content: "Name", widthPercentage: 35 },
        { content: "Date", widthPercentage: 12 },
        { content: "Issuer", widthPercentage: 25 },
        { content: "Type", widthPercentage: 13 },
        { content: "Status", widthPercentage: 10 },
        { content: "", widthPercentage: 5 }]

    return (
        <div className={`page pageWithCols pageWithMenu`}>
            <Menu />
            <div>
                <PageTitle title="Requests overview" description="Take a detailed view of requests made for you." />


                <div className={RequestOverviewPageCSS.table} >
                    <div className={RequestOverviewPageCSS.radioContainer}>
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

                <ModalWindow
                    height="55%"
                    isOpen={declineIsOpen}
                    closeWithdrawalModal={closeDeclineModal}
                    okWithdrawalModal={okDeclineModal}
                    title="Decline request"
                    buttonText="DECLINE" >
                    <p>To decline the request, you need to provide us some more info on why you want to decline it.</p>
                    <textarea placeholder='Write your reason ...' value={declineReason} onChange={handleReasonChange}></textarea>
                </ModalWindow>
            </div>
        </div>
    )
}

export default RequestOverviewPage