import Menu from 'components/navigation/Menu/Menu'
import PageTitle from 'components/view/PageTitle/PageTitle'
import ProfilePageCSS from './ProfilePage.module.scss'
import React, { useState } from 'react'

import Profile from "assets/menu/profile.png"
import InputField from 'components/forms/InputField/InputField'
import Button from 'components/forms/Button/Button'
import { AuthService } from 'features/auth/services/AuthService'

const ProfilePage = () => {
    const [firstName, setFirstName] = useState('');
    const [lastName, setLastName] = useState('');
    const [email, setEmail] = useState('');
    const [phone, setPhone] = useState('');

    (async function () {
        try {
            const data = await AuthService.getUserData();
            setFirstName(data.FirstName)
            setLastName(data.LastName)
            setEmail(data.Email)
            setPhone(data.Phone)
        } catch (error: any) {
          alert(error.response.data);
      }
    })()
    return (
        <div className={`page pageWithCols pageWithMenu`}>
            <Menu />
            <div>
                <PageTitle title="Profile" description="View your personal profile." />
                <div className={ProfilePageCSS.section}>
                    <div className={ProfilePageCSS.subsection}>
                        <div className={ProfilePageCSS.mainInfo}>
                            <img src={Profile} />
                            <div className={ProfilePageCSS.input}>
                                <InputField disabled={true} value={firstName} onChange={setFirstName} usage='First Name' className={ProfilePageCSS.input} />
                                <InputField disabled={true} value={lastName} onChange={setLastName} usage='Last Name' className={ProfilePageCSS.input} />
                            </div>
                        </div>
                        <InputField disabled={true} value={email} onChange={setEmail} usage='Email' className={ProfilePageCSS.input} />
                        <InputField disabled={true} value={phone} onChange={setPhone} usage='Phone' className={ProfilePageCSS.input} /> 
                    </div>
                </div>
            </div>
        </div >
    )
}

export default ProfilePage