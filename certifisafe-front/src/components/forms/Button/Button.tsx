import React from 'react'
import ButtonCSS from './Button.module.scss'

const Button = ({ text }: { text: string }) => {
    return (
        <button className={ButtonCSS.button}>{text}</button>
    )
}

export default Button