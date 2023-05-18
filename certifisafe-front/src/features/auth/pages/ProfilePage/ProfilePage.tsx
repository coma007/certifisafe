import Menu from 'components/navigation/Menu/Menu'
import PageTitle from 'components/view/PageTitle/PageTitle'
import ProfilePageCSS from './ProfilePage.module.scss'
import React from 'react'

import Profile from "assets/menu/profile.png"
import InputField from 'components/forms/InputField/InputField'
import Button from 'components/forms/Button/Button'

const ProfilePage = () => {
    return (
        <div className={`page pageWithCols pageWithMenu`}>
            <Menu />
            <div>
                <PageTitle title="Profile" description="View or edit your personal profile." />
                <div className={ProfilePageCSS.section}>
                    <div className={ProfilePageCSS.subsection}>
                        <div className={ProfilePageCSS.mainInfo}>
                            <img src={Profile} />
                            <div className={ProfilePageCSS.input}>
                                <InputField usage='First Name' className={ProfilePageCSS.input} />
                                <InputField usage='Last Name' className={ProfilePageCSS.input} />
                            </div>
                        </div>
                        <InputField usage='Email' className={ProfilePageCSS.input} />
                        <InputField usage='Phone' className={ProfilePageCSS.input} />
                        <Button text='SAVE' onClick={undefined} />
                    </div>
                    <div className={ProfilePageCSS.subsection}>
                        <div className={ProfilePageCSS.warning}>
                            It is recommended to change passwords often, for your security. <br />
                            To ensure the best security, you should not use one of your previous passwords.
                        </div>
                        <InputField usage='Old Password' className={ProfilePageCSS.input} />
                        <InputField usage='New Password' className={ProfilePageCSS.input} />
                        <InputField usage='Confirm Password' className={ProfilePageCSS.input} />
                        <Button text='CHANGE' onClick={undefined} />
                    </div>
                </div>
            </div>
        </div >
    )
}

export default ProfilePage