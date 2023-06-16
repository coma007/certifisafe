import Button from 'components/forms/Button/Button'
import React from 'react'
import Modal from "react-modal"
import ModalWindowCSS from "./ModalWindow.module.scss"

const ModalWindow = (props: { height: string, isOpen: any, closeWithdrawalModal: any, title: string, buttonText: string, children?: React.ReactNode, formId?: string}) => {
    Modal.setAppElement('#root')
    return (
        <Modal style={{
            content: {
                width: '45%',
                height: `${props.height}`,
                top: '50%',
                left: '50%',
                transform: 'translate(-50%, -50%)',
                borderRadius: '10px',
            },
            overlay: {
                backgroundColor: 'rgba(0, 0, 0, 0.2)',
            }
        }}
            isOpen={props.isOpen}
            onRequestClose={props.closeWithdrawalModal}>
            <div className={ModalWindowCSS.window}>
                <h2>{props.title}</h2>
                {props.children}
                <div className={ModalWindowCSS.buttons}>
                    <button className={ModalWindowCSS.accentButton}
                        type={props.formId ? "submit" : undefined} form={props.formId}>{props.buttonText}</button>
                    <button className={ModalWindowCSS.cancelButton} onClick={props.closeWithdrawalModal}>Discard</button>
                </div>
            </div>
        </Modal>
    )
}

export default ModalWindow