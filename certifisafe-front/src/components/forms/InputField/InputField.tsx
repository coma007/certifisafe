import InputFieldCSS from "./InputField.module.scss"

const InputField = ({ usage, className, value, onChange }: { usage: string, className: string, value?: string, onChange?: any }) => {
    return (
        <span>
            <input className={`${InputFieldCSS.input}, ${className}`} placeholder={usage} value={value} onChange={onChange} type={usage.toLowerCase().includes("password") ? "password" : "text"} />
        </span>
    )
}

export default InputField