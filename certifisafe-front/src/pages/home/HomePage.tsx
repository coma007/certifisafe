import Menu from 'components/navigation/Menu/Menu'
import HomePageCSS from './HomePage.module.scss'

const HomePage = () => {
    return (
        <div className={`page pageWithCols ${HomePageCSS.cols}`}>
            <Menu />
            <div className={HomePageCSS.content}>
                <div>
                    <h1> Welcome ! </h1>
                    <h2> Ready to track <br /> your certificates ? </h2>
                    <br /> <br />
                    <p> With our secure platform, <br />
                        you can have peace of mind knowing that your certificates are protected
                        and always at your fingertips.
                    </p>
                </div>
            </div>
            <div className={HomePageCSS.content}>Aaaaaa</div>
        </div>
    )
}

export default HomePage