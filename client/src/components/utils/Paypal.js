import React from 'react';
import PaypalExpressBtn from 'react-paypal-express-checkout';
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

        let env = 'sandbox';
        let currency = 'USD'; 
        let total = this.props.toPay;

        const client = {
            client_id: 'AXGVJI7L5hMkyrXv3YZ4RQTGWXRLRqdLLfj5brrJB_q553zlRdoy2wqlBUxGVXJsWuIpswmS2ZJign66',
            production: 'YOUR-PRODUCTION-APP-ID',
        }

        return (
            // <PaypalExpressBtn
            //     env={env}
            //     client={client}
            //     currency={currency}
            //     total={total}
            //     onError={onError}
            //     onSuccess={onSuccess}
            //     onCancel={onCancel}
            //     style={{ 
            //         size:'large',
            //         color:'blue',
            //         shape: 'rect',
            //         label: 'checkout'
            //     }}
            //      />

            <PayPalButton
                amount={total}
                // shippingPreference="NO_SHIPPING" // default is "GET_FROM_FILE"
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