import TooltipCSS from "./Tooltip.module.scss"

const Tooltip = ({ tooltipText }: { tooltipText: string }) => {
    return (
        <span className={TooltipCSS.tooltipText}>{tooltipText}</span>
    )
}

export default Tooltip