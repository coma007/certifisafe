const ErrorMsg = (props: { val: any }) => {
   
    return (
        <>
            { props.val ? <div>{props.val}</div> : null}
        </>
    )
}

export default ErrorMsg