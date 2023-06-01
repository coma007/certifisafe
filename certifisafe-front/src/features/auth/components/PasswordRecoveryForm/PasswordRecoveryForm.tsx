import Button from 'components/forms/Button/Button'
import InputField from 'components/forms/InputField/InputField'
import React from 'react'
import VerificationInput from 'react-verification-input'
import PasswordRecoveryFormCSS from "./PasswordRecoveryForm.module.scss"
import * as yup from 'yup' 
import { Formik, Form, Field, ErrorMessage } from 'formik';
import { useState } from 'react';


const PasswordRecoveryForm = () => {
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    
    const passwordValidator =  yup.string().min(8, "password is too short")
    .matches( /[a-z]+/, "needs to contain lowercase letter")
    .matches( /[A-Z]+/, "needs to contain uppercase letter")
    .matches( /[0-9]+/, "needs to contain number letter")
    .required();
    
    const schema = yup.object().shape({
        "new password": passwordValidator,
        "confirm password": passwordValidator.oneOf([yup.ref('new password')], 'Passwords must match'),
      })


    return (
        <Formik
        initialValues={{
            "new password": "",
            "confirm password": "",
        }}
        validationSchema={schema}
        onSubmit={values => {
    
        }}
      >
        {({ errors, touched, setFieldValue }) => (
            <Form>
                <VerificationInput length={4} placeholder={""} autoFocus={true} classNames={{
                    character: "codeField",
                    characterInactive: "codeFieldInactive",
                    characterSelected: "codeFieldActive",
                }} />
                <p className={PasswordRecoveryFormCSS.label}>
                    Did not receive code ?
                    <a onClick={undefined}>
                        <b> Click to resend.</b>
                    </a>
                </p>
                <Field name="new password" component={ InputField} className={PasswordRecoveryFormCSS.input}  usage="New password" value={newPassword} onChange={(e:React.ChangeEvent<HTMLInputElement>) => {
                        setNewPassword(e.target.value);
                        setFieldValue("new password", e.target.value);
                    }}/>
                <ErrorMessage name="new password" />

                <Field name="confirm password" component={ InputField} className={PasswordRecoveryFormCSS.input} usage="Confirm password" value={confirmPassword} onChange={(e:React.ChangeEvent<HTMLInputElement>) => {
                            setConfirmPassword(e.target.value);
                            setFieldValue("confirm password", e.target.value);
                        }}/>
                <ErrorMessage name="confirm password" />
                <span className="alignRight">
                    <Button submit={undefined} onClick={undefined} text="Reset" />
                </span>
            </Form>
    )}
    </Formik >
)}

export default PasswordRecoveryForm