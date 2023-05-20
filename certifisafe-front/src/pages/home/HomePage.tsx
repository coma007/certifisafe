import Menu from 'components/navigation/Menu/Menu'
import Card from 'components/view/Card/Card'
import HomePageCSS from './HomePage.module.scss'

import Certificate from "assets/menu/certificate.png"
import Verify from "assets/menu/verify.png"
import Create from "assets/menu/create.png"
import Request from "assets/menu/request.png"

const HomePage = () => {
    return (
        <div className={`page pageWithCols ${HomePageCSS.cols}`}>
            <Menu />
            <div className={HomePageCSS.content}>
                <div className={HomePageCSS.welcome}>
                    <h1> Welcome ! </h1>
                    <h2> Ready to track <br /> your certificates ? </h2>
                    <p> With our secure platform, <br />
                        you can have peace of mind knowing that your certificates are protected
                        and always at your fingertips.
                    </p>
                </div>
            </div>
            <div className={HomePageCSS.content}>
                <div className={HomePageCSS.grid}>
                    <Card link="/certificates">
                        <h3>Preview your certificates</h3>
                        <img src={Certificate} />
                        <p>Take a look of all of your certificates, including withdrawn and expired ones.</p>
                    </Card>
                    <Card link="/verify">
                        <h3>Verify certificate</h3>
                        <img src={Verify} />
                        <p>Verify any certificate by unique identificator or copy of the certificate.</p>
                    </Card>
                    <Card>
                        <h3>Create new certificate</h3>
                        <img src={Create} />
                        <p>Create request for intermediate or end certificate.</p>
                    </Card>
                    <Card link="/requests">
                        <h3>Preview your requests</h3>
                        <img src={Request} />
                        <p>Take a look of all certificate creation requests made by you and made for you.</p>
                    </Card>
                </div>
            </div>
        </div>
    )
}

export default HomePage