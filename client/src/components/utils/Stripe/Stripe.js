import React from 'react';
import { useState } from "react";
import { loadStripe } from "@stripe/stripe-js";

import CardIcon from "../../../assets/credit-card.svg";

import "../../../assets/styles/styles.css";

let stripePromise;

const getStripe = () => { 

  if (!stripePromise) {
    stripePromise = loadStripe("pk_test_51LdgVsCd5qIKDzW7ElrO6sYzxNzod4HLhDnCtB8T2A2nCrLg3RBP3cyp4jQRrHYBW6JkGsf4SceYpcdE1fZRKOtX00bfeNgS0J");
  }

  return stripePromise;
};

const Stripe = (props) => {
  console.log(props.data.cartItems)
  const [stripeError, setStripeError] = useState(null);
  const [isLoading, setLoading] = useState(false);
  const item = {
    price: "price_1K3TfMA4B8Maa00LFZ4EFwdX",
    quantity: 1
  };


  var lineItems = props.data.cartItems.map(function(item) { return {price: "10", quantity: item.quantity}; });
  console.log(lineItems)
  const checkoutOptions = {
    lineItems: lineItems,
    mode: "payment",
    successUrl: `${window.location.origin}/success`,
    cancelUrl: `${window.location.origin}/cancel`
  };

  const redirectToCheckout = async () => {
    setLoading(true);
    console.log("redirectToCheckout");

    const stripe = await getStripe();
    const { error } = await stripe.redirectToCheckout(checkoutOptions);
    console.log("Stripe checkout error", error);

    if (error) setStripeError(error.message);
    setLoading(false);
  };

  if (stripeError) alert(stripeError);

  return (
    <div className="checkout">
      <h1>Stripe Checkout</h1>
      <h1 className="checkout-price">$19</h1>
      <button
        className="checkout-button"
        onClick={redirectToCheckout}
        disabled={isLoading}
      >
        <div className="grey-circle">
          <div className="purple-circle">
            <img className="icon" src={CardIcon} alt="credit-card-icon" />
          </div>
        </div>
        <div className="text-container">
          <p className="text">{isLoading ? "Loading..." : "Buy"}</p>
        </div>
      </button>
    </div>
  );
};

export default Stripe;
