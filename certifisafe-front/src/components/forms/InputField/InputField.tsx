import React from 'react'
import InputFieldCSS from "./InputField.module.scss"

const InputField = ({ usage, className }: { usage: string, className: string }) => {
    return (
        <span>
            <input className={`${InputFieldCSS.input}, ${className}`} placeholder={usage}></input>
        </span>
    )
}

export default InputField