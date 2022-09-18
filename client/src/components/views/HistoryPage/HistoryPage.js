import React, { useState, useEffect } from 'react';
import { useDispatch } from 'react-redux';
import {
    getHistory
} from '../../../_actions/user_actions';

function HistoryPage(props) {
    const dispatch = useDispatch();

    const [history, setHistory] = useState([]);

    useEffect(() => {
        dispatch(getHistory())
            .then(response => {
                setHistory(response.payload.history)
            })
      }, []);

    return (
        <div style={{ width: '80%', margin: '3rem auto' }}>
            <div style={{ textAlign: 'center' }}>
                <h1>History</h1>
            </div>
            <br />

            <table>
                <thead>
                    <tr>
                        <th>Payment Id</th>
                        <th>Price</th>
                        <th>Quantity</th>
                        <th>Date of Purchase</th>
                    </tr>
                </thead>

                <tbody>

                    {history && history.map(item => (
                            <tr key={item.id}>
                                <td>{item.paymentId}</td>
                                <td>{item.price}</td>
                                <td>{item.quantity}</td>
                                <td>{item.dateOfPurchase}</td>
                            </tr>
                        ))}


                </tbody>
            </table>
        </div>
    )
}

export default HistoryPage