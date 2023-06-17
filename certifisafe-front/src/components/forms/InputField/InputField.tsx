import InputFieldCSS from "./InputField.module.scss"

const InputField = ({ usage, className, value, onChange, disabled }: { usage: string, className: string, value?: string, onChange?: any, disabled?: boolean }) => {
    return (
        <span>
            {disabled === undefined ?
            (<input className={`${InputFieldCSS.input}, ${className}`} placeholder={usage} value={value} onChange={onChange} type={usage.toLowerCase().includes("password") ? "password" : "text"} />)
            : (<input disabled className={`${InputFieldCSS.input}, ${className}`} placeholder={usage} value={value} onChange={onChange} type={usage.toLowerCase().includes("password") ? "password" : "text"} />)
}
        </span>
    )
}

export default InputField