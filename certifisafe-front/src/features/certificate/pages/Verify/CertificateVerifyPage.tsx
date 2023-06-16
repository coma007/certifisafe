import Button from 'components/forms/Button/Button'
import InputField from 'components/forms/InputField/InputField'
import Menu from 'components/navigation/Menu/Menu'
import PageTitle from 'components/view/PageTitle/PageTitle'
import CertificateVerifyPageCSS from './CertificateVerifyPage.module.scss'
import Card from 'components/view/Card/Card'
import { useRef, useState } from 'react'
import * as yup from 'yup'
import { Formik, Form, Field, ErrorMessage } from 'formik';

import Upload from "assets/actions/upload.png"
import Valid from "assets/actions/valid.png"
import Unvalid from "assets/actions/unvalid.png"
import { CertificateService } from 'features/certificate/services/CertificateService'
import { AxiosError } from 'axios'

const CertificateVerifyPage = () => {
  const [isValid, setValid] = useState<boolean | undefined>(undefined);
  const [certificateId, setCertificateId] = useState("");

  let fileRef = useRef<HTMLInputElement>(null);

  const handleResponse = (response: boolean) => {
    setValid(response);
  }


  const schemaId = yup.object().shape({
    "certificate id": yup.number().required()
      .test("int", "Must be integer", val => {
        return val === undefined ? false : val % 1 == 0;
      }).typeError("certificate id must be an integer"),


  })

  const fileSchema = yup.object().shape({
    file: yup.mixed()
      .test("is-file-too-big", "File exceeds 10MB", () => {
        let valid = true;
        const files = fileRef?.current?.files;
        if (files) {
          const fileArr = Array.from(files);
          fileArr.forEach((file) => {
            const size = file.size / 1024 / 1024;
            if (size > 10) {
              valid = false;
            }
          });
        }
        return valid;
      })
      .test(
        "is-file-of-correct-type",
        "File is not of supported type",
        () => {
          let valid = true;
          const files = fileRef?.current?.files;
          if (files) {
            const fileArr = Array.from(files);
            fileArr.forEach((file) => {
              const type = file.type.split("/")[1];
              const validTypes = [
                "pkix-cert",
              ];
              if (!validTypes.includes(type)) {
                valid = false;
              }
            });
          }
          return valid;
        }
      )
  })

  return (
    <div className={`page pageWithCols pageWithMenu`}>
      <Menu />
      <div>
        <PageTitle title="Verify certificate" description="Check if any certificate is valid." />
        <div className={`${CertificateVerifyPageCSS.block} pageWithCols`}>


          <div>
            <div className={CertificateVerifyPageCSS.section}>
              <Formik
                initialValues={{
                  certificateId: '',
                }}
                validationSchema={schemaId}
                onSubmit={values => {
                  CertificateService.verifyById(parseInt(certificateId)).then((response: boolean) => {
                    handleResponse(response);
                  }).catch((error: AxiosError) => {
                    alert(error.response?.data);
                    setValid(undefined);
                  })
                }}
              >
                {({ errors, touched, setFieldValue }) => (
                  <Form>
                    <b>Verify certificate by ID</b>
                    <br />
                    <br />
                    <Field name="certificate id" component={InputField} className={CertificateVerifyPageCSS.textInput} usage="Enter Certificate ID" value={certificateId} onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                      setCertificateId(e.target.value);
                      setFieldValue("certificate id", e.target.value);
                    }} />
                    <ErrorMessage name="certificate id" />
                    <span className={CertificateVerifyPageCSS.button}>
                      <Button submit={"submit"} text={'VERIFY'} onClick={undefined} />
                    </span>
                    <br />
                    <small>Certificate ID can be found on the bottom of every certificate. </small>


                  </Form>
                )}
              </Formik>
            </div>

            <small>or</small>

            <div className={CertificateVerifyPageCSS.section}>

              <Formik
                initialValues={{
                  file: "",
                }}
                validationSchema={fileSchema}
                onSubmit={values => {
                  const files = fileRef?.current?.files;
                  if (files) {
                    const fileArr = Array.from(files);
                    fileArr.forEach((file) => {
                      CertificateService.verifyByFile(file).then((response: boolean) => {
                        handleResponse(response);
                      }).catch((error: AxiosError) => {
                        alert(error.response?.data);
                        setValid(undefined);
                      })
                    });
                  }

                }}
              >
                {({ errors, touched, setFieldValue }) => (
                  <Form>
                    <b>Verify certificate by its copy</b>
                    <br />
                    <br />
                    <span className={CertificateVerifyPageCSS.inline}>
                      <label htmlFor="file-upload" className={CertificateVerifyPageCSS.fileUpload}>
                        Upload a copy<img src={Upload} />
                      </label>
                      <span className={CertificateVerifyPageCSS.button}>
                        <Button submit={"submit"} text={'VERIFY'} onClick={undefined} />
                      </span>
                    </span>
                    <input name="file" ref={fileRef} id="file-upload" type="file" />
                    <ErrorMessage name="file" />
                  </Form>
                )}
              </Formik>
            </div>
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
    </div >
  )
}

export default CertificateVerifyPage