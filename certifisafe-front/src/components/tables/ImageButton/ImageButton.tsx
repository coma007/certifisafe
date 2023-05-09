import Tooltip from "components/view/Tooltip/Tooltip";
import ImageButtonCSS from "./ImageButton.module.scss"
import TooltipCSS from "components/view/Tooltip/Tooltip.module.scss"

const ImageButton = (props: { path: string, tooltipText: string }) => {
    return (
        <a className={`${ImageButtonCSS.button} ${TooltipCSS.bottomTooltip}`}>
            <img src={props.path}></img>
            <Tooltip tooltipText={props.tooltipText} />
        </a>
    )
}

export default ImageButton;