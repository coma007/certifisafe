import React, { useRef } from 'react'

import ModalWindow from 'components/view/Modal/ModalWindow'
import ModalWindowCSS from 'components/view/Modal/ModalWindow.module.scss'
import RequestCreatePageCSS from './RequestCreatePage.module.scss'
import { useState } from 'react'
import InputField from 'components/forms/InputField/InputField'
import * as yup from 'yup'
import { Formik, Form, Field, ErrorMessage } from 'formik';
import ReCAPTCHA from 'react-google-recaptcha';
import ReactDOM from 'react-dom'
import { RequestService } from 'features/request/service/RequestService'
import { useNavigate } from 'react-router'
import ErrorMsg from 'components/error/ErrorMsg'

const RequestCreatePage = (props: { createIsOpen: boolean, okCreateModal : any, closeCreateModal: any }) => {

    const [signerSerial, setSignerSerial] = useState('');
    const [name, setName] = useState('');
    const [type, setType] = useState('');

    const schema = yup.object().shape({
        "signer serial": yup.number().required(),
        "name": yup.string().required(),
        "type": yup.string().required(),
    })

    const captchaRef: any = useRef(null)

    const handleSubmit = async (values: any) => {
        try {
            const token = captchaRef.current?.getValue();
            captchaRef.current.reset();
            await RequestService.createRequest({ ParentSerial: parseInt(signerSerial), CertificateType: values.type, CertificateName: name, Token: token })
            props.closeCreateModal()
        } catch (error: any) {
            alert(error.response.data);
        }
    }
    return (
        <ModalWindow
            height="82%"
            isOpen={props.createIsOpen}
            closeWithdrawalModal={props.closeCreateModal}
            okWithdrawalModal={props.okCreateModal}
            title="Create new certificate"
            buttonText="REQUEST"
            formId="request-create-form">

            <Formik
                initialValues={{
                    "signer serial": "",
                    "name": "",
                    "type": "",
                    "token": "",
                }}
                validationSchema={schema}
                validateOnChange
                onSubmit={values => handleSubmit(values)}
            >
                {({ errors, touched, setFieldValue, validateForm }) => (
                    <Form id="request-create-form">
                        <p>Make a request for creating new certificate. </p>
                        <Field name="signer serial" component={InputField} usage="Signing certificate serial number" className={ModalWindowCSS.input}
                            value={signerSerial} onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                                setSignerSerial(e.target.value);
                                setFieldValue("signer serial", e.target.value);
                            }} />
                        <ErrorMsg val={errors["signer serial"]} />

                        <Field name="name" component={InputField} usage="Certificate name" className={ModalWindowCSS.input}
                            value={name} onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                                setName(e.target.value);
                                setFieldValue("name", e.target.value);
                            }} />
                        <ErrorMsg val={errors["name"]} />

                        <div>
                            <p>
                                Type of certificate:
                                <div className={RequestCreatePageCSS.option}>
                                    <Field
                                        type="radio"
                                        name="type"
                                        id="foryou"
                                        value="root"
                                    // checked={selectedOption === "root"}
                                    // onChange={handleOptionChange}
                                    />
                                    <label htmlFor="foryou">Root</label>
                                    <span className={RequestCreatePageCSS.description}>Signed by our root certificate and can sign other certificates.</span>
                                </div>
                                <div className={RequestCreatePageCSS.option}>
                                    <Field
                                        type="radio"
                                        name="type"
                                        id="fromyou"
                                        value="intermediate"
                                    // checked={selectedOption === "intermediate"}
                                    // onChange={handleOptionChange}
                                    />
                                    <label htmlFor="fromyou">Intermediate</label>
                                    <span className={RequestCreatePageCSS.description}>Must be signed by other certificate and can sign other certificates.</span>
                                </div>
                                <div className={RequestCreatePageCSS.option}>
                                    <Field
                                        type="radio"
                                        name="type"
                                        id="fromyou"
                                        value="end"
                                    // checked={selectedOption === "end"}
                                    // onChange={handleOptionChange}
                                    />
                                    <label htmlFor="fromyou">End</label>
                                    <span className={RequestCreatePageCSS.description}>Must be signed by other certificate and cannot sign other certificates.</span>
                                </div>
                            </p>
                            <ErrorMsg val={errors["type"]} />
                        </div>
                        <p className={ModalWindowCSS.warning}>Please note that this certificate will be active if the request is accepted by the signing certificate owner. If so, certificate will be valid from that day until the expiration of the signing certificate.</p>
                        <ReCAPTCHA className='recaptcha' sitekey={process.env.RECAPTCHA_SITE_KEY as string} ref={captchaRef} />
                    </Form>
                )}
            </Formik >
        </ModalWindow >
    )
}

export default RequestCreatePage
