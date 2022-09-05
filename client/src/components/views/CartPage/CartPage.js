import React, { useEffect, useState } from 'react'
import { useDispatch } from 'react-redux';
import {
    getCartItems,
    removeCartItem,
    onSuccessBuy
} from '../../../_actions/user_actions';
import UserCardBlock from './Sections/UserCardBlock';
import { Result, Empty } from 'antd';
import Axios from 'axios';
import Paypal from '../../utils/Paypal';

function CartPage(props) {
    const dispatch = useDispatch();
    const [Total, setTotal] = useState(0)
    const [ShowTotal, setShowTotal] = useState(false)
    const [ShowSuccess, setShowSuccess] = useState(false)

    useEffect(() => {

        let cartItems = [];
        if (props.user.userData && props.user.userData.cartItems) {
            if (props.user.userData.cartItems.length > 0) {
                calculateTotal(props.user.userData.cartItems)

                props.user.userData.cartItems.forEach(item => {
                    cartItems.push(item.productId)
                });

                // dispatch(getCartItems(cartItems, props.user.userData.cartItems))
                //     .then((response) => {
                //         if (response.payload.length > 0) {
                //             calculateTotal(response.payload)
                //         }
                //     })
            }
        }

    }, [props.user.userData])

    const calculateTotal = (cartItems) => {
        let total = 0;
        
        cartItems.map(item => {
            total += parseInt(item.price, 10) * item.quantity
        });

        setTotal(total)
        setShowTotal(true)
    }

    const removeFromCart = (productId) => {

        dispatch(removeCartItem(productId))
            .then((response) => {
                console.log(response.payload)
                if (response.payload.cart.items.length <= 0) {
                    setShowTotal(false)
                } else {
                    calculateTotal(response.payload.cart.items)
                }
            })
    }

    const transactionSuccess = (data) => {
        dispatch(onSuccessBuy({
            cartDetail: props.user.cartDetail,
            paymentData: data
        }))
            .then(response => {
                if (response.payload.success) {
                    setShowSuccess(true)
                    setShowTotal(false)
                }
            })
    }

    const transactionError = () => {
        console.log('Paypal error')
    }

    const transactionCanceled = () => {
        console.log('Transaction canceled')
    }

    return (
        <div style={{ width: '85%', margin: '3rem auto' }}>
            <h1>My Cart</h1>
            <div>

                <UserCardBlock
                    products={props.user.userData}
                    removeItem={removeFromCart}
                />


                {ShowTotal ?
                    <div style={{ marginTop: '3rem' }}>
                        <h2>Total amount: ${Total} </h2>
                    </div>
                    :
                    ShowSuccess ?
                        <Result
                            status="success"
                            title="Successfully Purchased Items"
                        /> :
                        <div style={{
                            width: '100%', display: 'flex', flexDirection: 'column',
                            justifyContent: 'center'
                        }}>
                            <br />
                            <Empty description={false} />
                            <p>No Items In the Cart</p>

                        </div>
                }
            </div>

            {ShowTotal &&
                <Paypal
                    toPay={Total}
                    onSuccess={transactionSuccess}
                    transactionError={transactionError}
                    transactionCanceled={transactionCanceled}
                />
            }
            
        </div>
    )
}

export default CartPage