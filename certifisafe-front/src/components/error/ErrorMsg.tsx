import ErrorMessageCSS from "./ErrorMsg.module.scss"

const ErrorMsg = (props: { val: any, customClass?: any }) => {

    return (
        <>
            {props.val ? <div className={`${ErrorMessageCSS.error} ${props.customClass}`}>{props.val}</div> : null}
        </>
    )
}

export default ErrorMsg