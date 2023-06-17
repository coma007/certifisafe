import PageTitleCSS from "./PageTitle.module.scss"

const PageTitle = (props: { title: string, description: string }) => {
    return (
        <div className={PageTitleCSS.titleBox}>
            <h1 className={PageTitleCSS.title}>
                {props.title}
            </h1>
            <div className={PageTitleCSS.description}>
                {props.description}
            </div>
        </div>
    )
}

export default PageTitle