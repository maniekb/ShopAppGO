import React from 'react'

function UserCardBlock(props) {

    //const renderCartImage = (images) => {
        // if(images.length > 0) {
        //     let image = images[0]
        //     return `http://localhost:5000/${image}`
        // }
    //}

    const renderItems = () => {
        console.log(props.products)
        return props.products && props.products.cartItems.map(product => (
            <tr key={product.productId}>
                <td>
                    <img style={{ width: '70px' }} alt="product" 
                    /* src={renderCartImage(product.images)} *//>
                </td> 
                <td>{product.quantity} EA</td>
                <td>$ {product.price} </td>
                <td><button 
                onClick={()=> {props.removeItem(product.productId)}}
                >Remove </button> </td>
            </tr>
        ))
    }


    return (
        <div>
            <table>
                <thead>
                    <tr>
                        <th>Product Image</th>
                        <th>Product Quantity</th>
                        <th>Product Price</th>
                        <th>Remove from Cart</th>
                    </tr>
                </thead>
                <tbody>
                    {renderItems()}
                </tbody>
            </table>
        </div>
    )
}

export default UserCardBlock