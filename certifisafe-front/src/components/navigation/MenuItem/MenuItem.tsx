import Tooltip from '../Tooltip/Tooltip'
import MenuItemCSS from './MenuItem.module.scss'
import TooltipCSS from '../Tooltip/Tooltip.module.scss'

const MenuItem = ({ image, className, tooltipText, tooltip }:
    { image: string, className: string, tooltipText: string, tooltip: boolean }) => {
    return (
        <a className={`${tooltip ? TooltipCSS.tooltip : TooltipCSS.noTooltip}`}>
            <img className={`${MenuItemCSS.item} ${className}`} src={image} />
            <Tooltip tooltipText={tooltipText} />
        </a>
    )
}

export default MenuItem