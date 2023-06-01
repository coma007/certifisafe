import Button from 'components/forms/Button/Button'
import InputField from 'components/forms/InputField/InputField'
import React from 'react'
import EmailFormCSS from "./EmailForm.module.scss"

const EmailForm = (props: { onClick: any }) => {
    return (
        <div>
            <InputField usage='Email' className={EmailFormCSS.input} />
            <span className="alignRight">
                <Button submit={undefined} onClick={props.onClick} text="Send code" />
            </span>
        </div>
    )
}

export default EmailForm