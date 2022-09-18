import React, { useEffect, useState } from 'react'
import { useDispatch } from 'react-redux';
import {
    removeCartItem,
    onSuccessBuy
} from '../../../_actions/user_actions';
import UserCardBlock from './Sections/UserCardBlock';
import { Result, Empty } from 'antd';
import Paypal from '../../utils/Paypal';

function CartPage(props) {
    const dispatch = useDispatch();
    const [Total, setTotal] = useState(0)
    const [ShowTotal, setShowTotal] = useState(false)
    const [ShowSuccess, setShowSuccess] = useState(false)

    useEffect(() => {
        if (props.user.userData && props.user.userData.cartItems) {
            if (props.user.userData.cartItems.length > 0) {
                calculateTotal(props.user.userData.cartItems)
            }
        }

    }, [props.user.userData])

    const calculateTotal = (cartItems) => {
        let total = 0;
        
        cartItems.map(item => {
            total += parseInt(item.price, 10)
        });

        setTotal(total)
        setShowTotal(true)
    }

    function removeElementWithId(arr, productId) {
        const objWithIdIndex = arr.findIndex((obj) => obj.productId === productId);
        arr.splice(objWithIdIndex, 1);
      
        return arr;
      }

    const removeFromCart = (productId) => {

        dispatch(removeCartItem(productId))
            .then((response) => {
                removeElementWithId(props.user.userData.cartItems, productId)
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

    const showSummary = () => {
        if(ShowTotal) {
            return (
                <div style={{ marginTop: '3rem' }}>
                    <h2>Total amount: ${Total} </h2>
                </div>
            )
        }

        if(ShowSuccess) {
            return (
                <Result
                    status="success"
                    title="Successfully Purchased Items"
                />
            )
        } else {
            return (
                <div style={{
                    width: '100%', display: 'flex', flexDirection: 'column',
                    justifyContent: 'center'
                }}>
                    <br />
                    <Empty description={false} />
                    <p>No Items In the Cart</p>

                </div>
            )
        }
    }

    return (
        <div style={{ width: '85%', margin: '3rem auto' }}>
            <h1>My Cart</h1>
            <div>

                <UserCardBlock
                    products={props.user.userData}
                    removeItem={removeFromCart}
                />


                {showSummary()}
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