import { Link } from "react-router-dom"
import CardCSS from "./Card.module.scss"

const Card = (props: { children: React.ReactNode, link?: string }) => {
    return (
        <a href={props.link} className={CardCSS.card} > {props.children}</a >
    )
}

export default Card