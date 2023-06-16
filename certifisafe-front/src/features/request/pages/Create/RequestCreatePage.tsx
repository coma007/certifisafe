import React from 'react'

import ModalWindow from 'components/view/Modal/ModalWindow'
import ModalWindowCSS from 'components/view/Modal/ModalWindow.module.scss'
import RequestCreatePageCSS from './RequestCreatePage.module.scss'
import { useState } from 'react'
import InputField from 'components/forms/InputField/InputField'

const RequestCreatePage = (props: { createIsOpen: boolean, closeCreateModal: any }) => {
    return (
        <ModalWindow
            height="82%"
            isOpen={props.createIsOpen}
            closeWithdrawalModal={props.closeCreateModal}
            title="Create new certificate"
            buttonText="REQUEST" >
            <p>Make a request for creating new certificate. </p>
            <InputField usage="Signing certificate serial number" className={ModalWindowCSS.input} />
            <InputField usage="Certificate name" className={ModalWindowCSS.input} />
            <div>
                <p>
                    Type of certificate:
                    <div className={RequestCreatePageCSS.option}>
                        <input
                            type="radio"
                            name="radioGroup"
                            id="foryou"
                        // checked={selectedOption === "root"}
                        // onChange={handleOptionChange}
                        />
                        <label htmlFor="foryou">Root</label>
                        <span className={RequestCreatePageCSS.description}>Signed by our root certificate and can sign other certificates.</span>
                    </div>
                    <div className={RequestCreatePageCSS.option}>
                        <input
                            type="radio"
                            name="radioGroup"
                            id="fromyou"
                        // checked={selectedOption === "intermediate"}
                        // onChange={handleOptionChange}
                        />
                        <label htmlFor="fromyou">Intermediate</label>
                        <span className={RequestCreatePageCSS.description}>Must be signed by other certificate and can sign other certificates.</span>
                    </div>
                    <div className={RequestCreatePageCSS.option}>
                        <input
                            type="radio"
                            name="radioGroup"
                            id="fromyou"
                        // checked={selectedOption === "end"}
                        // onChange={handleOptionChange}
                        />
                        <label htmlFor="fromyou">End</label>
                        <span className={RequestCreatePageCSS.description}>Must be signed by other certificate and cannot sign other certificates.</span>
                    </div>
                </p>
            </div>
            <p className={ModalWindowCSS.warning}>Please note that this certificate will be active if the request is accepted by the signing certificate owner. If so, certificate will be valid from that day until the expiration of the signing certificate.</p>
        </ModalWindow >
    )
}

export default RequestCreatePage