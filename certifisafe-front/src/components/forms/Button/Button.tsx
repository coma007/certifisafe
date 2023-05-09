import ButtonCSS from './Button.module.scss'

const Button = (props: { text: string, onClick: any }) => {
    return (
        <button className={ButtonCSS.button} onClick={props.onClick}>{props.text}</button>
    )
}

export default Button