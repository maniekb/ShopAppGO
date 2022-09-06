import React from 'react';
import { PayPalButton } from "react-paypal-button-v2";


export default class Paypal extends React.Component {
    render() {
        const onSuccess = (payment) => {
            console.log("The payment was succeeded!", payment);
            this.props.onSuccess(payment);
        
        }

        const onCancel = (data) => {
            console.log('The payment was cancelled!', data);
        }

        const onError = (err) => {
            console.log("Error!", err);
        }
 
        let total = this.props.toPay;

        return (
            <PayPalButton
                amount={total}
                onError={onError}
                onSuccess={onSuccess}
                onCancel={onCancel}
                options={{
                clientId: "AXGVJI7L5hMkyrXv3YZ4RQTGWXRLRqdLLfj5brrJB_q553zlRdoy2wqlBUxGVXJsWuIpswmS2ZJign66"
                }}
            />
        );
    }
}