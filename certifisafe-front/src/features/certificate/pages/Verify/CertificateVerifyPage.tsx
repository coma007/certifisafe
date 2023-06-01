import Button from 'components/forms/Button/Button'
import InputField from 'components/forms/InputField/InputField'
import Menu from 'components/navigation/Menu/Menu'
import PageTitle from 'components/view/PageTitle/PageTitle'
import CertificateVerifyPageCSS from './CertificateVerifyPage.module.scss'
import Card from 'components/view/Card/Card'
import { useState } from 'react'

import Upload from "assets/actions/upload.png"
import Valid from "assets/actions/valid.png"
import Unvalid from "assets/actions/unvalid.png"

const CertificateVerifyPage = () => {
    const [isValid, setValid] = useState<boolean | undefined>(undefined);


    return (
        <div className={`page pageWithCols pageWithMenu`}>
            <Menu />
            <div>
                <PageTitle title="Verify certificate" description="Check if any certificate is valid." />
                <div className={`${CertificateVerifyPageCSS.block} pageWithCols`}>
                    <div>
                        <div className={CertificateVerifyPageCSS.section}>
                            <b>Verify certificate by ID</b>
                            <br />
                            <br />
                            <InputField usage={'Enter Certificate ID'} className={CertificateVerifyPageCSS.textInput} />
                            <br />
                            <small>Certificate ID can be found on the bottom of every certificate. </small>
                        </div>
                        <small>or</small>
                        <div className={CertificateVerifyPageCSS.section}>
                            <b>Verify certificate by its copy</b>
                            <br />
                            <br />
                            <label htmlFor="file-upload" className={CertificateVerifyPageCSS.fileUpload}>
                                Upload a copy<img src={Upload} />
                            </label>
                            <input id="file-upload" type="file" />
                        </div>
                        <Button submit={undefined} text={'VERIFY'} onClick={undefined} />
                    </div>
                    <div className={CertificateVerifyPageCSS.card}>
                        {(isValid !== undefined) &&
                            <Card>
                                {(isValid === true) &&
                                    <div>
                                        <img src={Valid} />
                                        <h3>Certificate is <br /> VALID !</h3>
                                    </div>

                                }
                                {(isValid === false) &&
                                    <div>
                                        <img src={Unvalid} />
                                        <h3>Certificate is <br /> UNVALID !</h3>
                                    </div>

                                }
                            </Card>
                        }
                    </div>
                </div>
            </div>
        </div>
    )
}

export default CertificateVerifyPage