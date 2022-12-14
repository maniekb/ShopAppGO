import React, { useEffect, useState } from 'react'
import Axios from 'axios'
import { Row, Col } from 'antd';
import ProductImage from './Sections/ProductImage';
import ProductInfo from './Sections/ProductInfo';
import { addToCart } from '../../../_actions/user_actions';
import { useDispatch, useSelector } from 'react-redux';

function DetailProductPage(props) {
    const user = useSelector(state => state.user)
    const dispatch = useDispatch();
    const productId = props.match.params.productId
    const [Product, setProduct] = useState([])

    useEffect(() => {
        Axios.get(`/api/product?id=${productId}&type=single`)
            .then(response => {
                setProduct(response.data.product)
            })

    }, [])

    const addToCartHandler = (productId) => {
        dispatch(addToCart(productId))
            .then((response) => {
                if (response.payload.success) {
                    alert("Product successfully added to cart!")
                }
            })
    }

    return (
        <div className="postPage" style={{ width: '100%', padding: '3rem 4rem' }}>

            <div style={{ display: 'flex', justifyContent: 'center' }}>
                <h1>{Product.title}</h1>
            </div>

            <br />

            <Row gutter={[16, 16]} >
                <Col lg={12} xs={24}>
                    <ProductImage detail={Product} />
                </Col>
                <Col lg={12} xs={24}>
                    <ProductInfo
                        isAuth={user.userData && user.userData.isAuth}
                        addToCart={addToCartHandler}
                        detail={Product} />
                </Col>
            </Row>
        </div>
    )
}

export default DetailProductPage