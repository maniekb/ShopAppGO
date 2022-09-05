import React, { useEffect, useState } from 'react'
import { Button, Descriptions } from 'antd';
import { useSelector } from 'react-redux';

function ProductInfo(props) {
    const user = useSelector(state => state.user)
    const [Product, setProduct] = useState({})

    useEffect(() => {

        setProduct(props.detail)

    }, [props.detail])

    const addToCartHandler = () => {
        props.addToCart(props.detail.id)
    }

    return (
        <div>
            <Descriptions title="Product Info">
                <Descriptions.Item label="Price"> {Product.price}</Descriptions.Item>
                <Descriptions.Item label="Sold">{Product.sold}</Descriptions.Item>
                <Descriptions.Item label="Views"> {Product.views}</Descriptions.Item>
                <Descriptions.Item label="Description"> {Product.description}</Descriptions.Item>
            </Descriptions>

            <br />
            <br />
            <br />
            { props.isAuth && 
                <div style={{ display: 'flex', justifyContent: 'center' }}>
                <Button size="large" shape="round" type="danger"
                    onClick={addToCartHandler}
                >
                    Add to Cart
                    </Button>
                </div>
            }
        </div>
    )
}

export default ProductInfo