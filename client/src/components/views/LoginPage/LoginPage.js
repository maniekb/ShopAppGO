import React, { useState } from "react";
import { withRouter, useLocation } from "react-router-dom";
import { loginUser } from "../../../_actions/user_actions";
import { Formik } from 'formik';
import * as Yup from 'yup';
import { Form, Icon, Input, Button, Checkbox, Typography } from 'antd';
import { useDispatch } from "react-redux";

import { ReactComponent as GoogleLogo } from '../../../assets/google.svg';
import { ReactComponent as GitHubLogo } from '../../../assets/github.svg';
import { ReactComponent as FacebookLogo } from '../../../assets/facebook.svg';
import { getGoogleUrl } from '../../utils/GoogleAuth';
import { getGitHubUrl } from '../../utils/GitHubAuth'; 
import { getFacebookUrl } from '../../utils/FacebookAuth';


const { Title } = Typography;

function LoginPage(props) {
  const location = useLocation();
  let from = ((location.state)?.from?.pathname) || '/';
  const dispatch = useDispatch();
  const rememberMeChecked = localStorage.getItem("rememberMe") ? true : false;

  const [formErrorMessage, setFormErrorMessage] = useState('')
  const [rememberMe, setRememberMe] = useState(rememberMeChecked)

  const handleRememberMe = () => {
    setRememberMe(!rememberMe)
  };

  const initialEmail = localStorage.getItem("rememberMe") ? localStorage.getItem("rememberMe") : '';

  return (
      <Formik
        initialValues={{
          email: initialEmail,
          password: '',
        }}
        validationSchema={Yup.object().shape({
          email: Yup.string()
            .email('Email is invalid')
            .required('Email is required'),
          password: Yup.string()
            .min(6, 'Password must be at least 6 characters')
            .required('Password is required'),
        })}
        onSubmit={(values, { setSubmitting }) => {
          setTimeout(() => {
            let dataToSubmit = {
              email: values.email,
              password: values.password
            };

            dispatch(loginUser(dataToSubmit))
              .then(response => {
                if (response.payload.success) {
                  window.localStorage.setItem('userId', response.payload.userId);
                  if (rememberMe === true) {
                    window.localStorage.setItem('rememberMe', values.id);
                  } else {
                    localStorage.removeItem('rememberMe');
                  }
                  props.history.push("/");
                } else {
                  setFormErrorMessage('Check out your Account or Password again')
                }
              })
              .catch(err => {
                setFormErrorMessage('Check out your Account or Password again')
                setTimeout(() => {
                  setFormErrorMessage("")
                }, 3000);
              });
            setSubmitting(false);
          }, 500);
        }}
      >
        {props => {
          const {
            values,
            touched,
            errors,
            dirty,
            isSubmitting,
            handleChange,
            handleBlur,
            handleSubmit,
            handleReset,
          } = props;
          return (
            <div className="app">

              <Title level={2}>Log In</Title>
              <form onSubmit={handleSubmit} style={{ width: '350px' }}>

                <Form.Item required>
                  <Input
                    id="email"
                    prefix={<Icon type="user" style={{ color: 'rgba(0,0,0,.25)' }} />}
                    placeholder="Enter your email"
                    type="email"
                    value={values.email}
                    onChange={handleChange}
                    onBlur={handleBlur}
                    className={
                      errors.email && touched.email ? 'text-input error' : 'text-input'
                    }
                  />
                  {errors.email && touched.email && (
                    <div className="input-feedback">{errors.email}</div>
                  )}
                </Form.Item>

                <Form.Item required>
                  <Input
                    id="password"
                    prefix={<Icon type="lock" style={{ color: 'rgba(0,0,0,.25)' }} />}
                    placeholder="Enter your password"
                    type="password"
                    value={values.password}
                    onChange={handleChange}
                    onBlur={handleBlur}
                    className={
                      errors.password && touched.password ? 'text-input error' : 'text-input'
                    }
                  />
                  {errors.password && touched.password && (
                    <div className="input-feedback">{errors.password}</div>
                  )}
                </Form.Item>

                {formErrorMessage && (
                  <label ><p style={{ color: '#ff0000bf', fontSize: '0.7rem', border: '1px solid', padding: '1rem', borderRadius: '10px' }}>{formErrorMessage}</p></label>
                )}

                <Form.Item>
                  <Checkbox id="rememberMe" onChange={handleRememberMe} checked={rememberMe} >Remember me</Checkbox>
                  <div>
                    <Button type="primary" htmlType="submit" className="login-form-button" style={{ minWidth: '100%' }} disabled={isSubmitting} onSubmit={handleSubmit}>
                      Log in
                  </Button>
                  </div>
                  Or <a href="/register">register now!</a>
                </Form.Item>


                <Form.Item>
                  Log in with another provider:
                  <div>
                    <Button
                      href={getGoogleUrl(from)}
                      sx={{
                        backgroundColor: '#f5f6f7',
                        borderRadius: 1,
                        py: '0.6rem',
                        columnGap: '1rem',
                        textDecoration: 'none',
                        color: '#393e45',
                        cursor: 'pointer',
                        fontWeight: 500,
                        '&:hover': {
                          backgroundColor: '#fff',
                          boxShadow: '0 1px 13px 0 rgb(0 0 0 / 15%)',
                        },
                      }}
                      display='flex' 
                      justifyContent='center'
                      alignItems='center'
                    >
                      <GoogleLogo style={{ height: '2rem' }} />
                       Google
                    </Button>
                  </div>
                  <div>
                    <Button
                      href={getGitHubUrl(from)}
                      sx={{
                        backgroundColor: '#f5f6f7',
                        borderRadius: 1,
                        py: '0.6rem',
                        mt: '1.5rem',
                        columnGap: '1rem',
                        textDecoration: 'none',
                        color: '#393e45',
                        cursor: 'pointer',
                        fontWeight: 500,
                        '&:hover': {
                          backgroundColor: '#fff',
                          boxShadow: '0 1px 13px 0 rgb(0 0 0 / 15%)',
                        },
                      }}
                      display='flex'
                      justifyContent='center'
                      alignItems='center'
                      >
                      <GitHubLogo style={{ height: '2rem' }} />
                       GitHub
                    </Button>    
                  </div>
                  <div>
                    <Button
                      href={getFacebookUrl(from)}
                      sx={{
                        backgroundColor: '#f5f6f7',
                        borderRadius: 1,
                        py: '0.6rem',
                        mt: '1.5rem',
                        columnGap: '1rem',
                        textDecoration: 'none',
                        color: '#393e45',
                        cursor: 'pointer',
                        fontWeight: 500,
                        '&:hover': {
                          backgroundColor: '#fff',
                          boxShadow: '0 1px 13px 0 rgb(0 0 0 / 15%)',
                        },
                      }}
                      display='flex'
                      justifyContent='center'
                      alignItems='center'
                      >
                      <FacebookLogo style={{ height: '2rem' }} />
                      Facebook
                    </Button>    
                  </div>
                </Form.Item>
              </form>
            </div>
          );
        }}
      </Formik>
  );
};

export default withRouter(LoginPage);


