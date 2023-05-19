import React from 'react'

import ModalWindow from 'components/view/Modal/ModalWindow'
import ModalWindowCSS from 'components/view/Modal/ModalWindow.module.scss'
import { useState } from 'react'
import InputField from 'components/forms/InputField/InputField'

const RequestCreatePage = (props: { createIsOpen: boolean, closeCreateModal: any }) => {
    return (
        <ModalWindow
            height="75%"
            isOpen={props.createIsOpen}
            closeWithdrawalModal={props.closeCreateModal}
            title="Create new certificate"
            buttonText="REQUEST" >
            <p>Make a request for creating new certificate. </p>
            <InputField usage="Signing certificate serial number" className={ModalWindowCSS.input} />
            <InputField usage="Certificate name" className={ModalWindowCSS.input} />
            <div>
                <input
                    type="radio"
                    name="radioGroup"
                    id="foryou"
                // checked={selectedOption === "foryou"}
                // onChange={handleOptionChange}
                />
                <label htmlFor="foryou">Root</label>
                <input
                    type="radio"
                    name="radioGroup"
                    id="fromyou"
                // checked={selectedOption === "fromyou"}
                // onChange={handleOptionChange}
                />
                <label htmlFor="fromyou">Intermediate</label>
                <input
                    type="radio"
                    name="radioGroup"
                    id="fromyou"
                // checked={selectedOption === "fromyou"}
                // onChange={handleOptionChange}
                />
                <label htmlFor="fromyou">End</label>
            </div>
            <p className={ModalWindowCSS.warning}>Please note that this certificate will be active if the request is accepted by the signing certificate owner. If so, certificate will be valid from that day until the expiration of the signing certificate.</p>
        </ModalWindow>
    )
}

export default RequestCreatePage