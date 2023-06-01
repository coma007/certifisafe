import ButtonCSS from './Button.module.scss'

const Button = (props: { text: string, onClick: any, submit: "button" | "submit" | "reset" | undefined }) => {
    return (
        <button  type={props.submit} className={ButtonCSS.button} onClick={props.onClick}>{props.text}</button>
    )
}

export default Button