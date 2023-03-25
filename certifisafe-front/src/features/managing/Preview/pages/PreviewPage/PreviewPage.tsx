import Menu from "components/navigation/Menu/Menu"
import PreviewPageCSS from "./PreviewPage.module.scss"

const PreviewPage = () => {
    return (
        <div className={`page pageWithCols ${PreviewPageCSS.cols}`}>
            <Menu />
        </div>
    )
}

export default PreviewPage