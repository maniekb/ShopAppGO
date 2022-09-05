import { FACEBOOK_OAUTH_CLIENT_ID, FACEBOOK_OAUTH_CLIENT_SECRET, FACEBOOK_OAUTH_REDIRECT_URL } from '../Config.js';


export function getFacebookUrl(from) {
    const rootURl = 'https://graph.facebook.com/oauth/authorize';
  
    const options = {
      client_id: FACEBOOK_OAUTH_CLIENT_ID,
      clientSecret: FACEBOOK_OAUTH_CLIENT_SECRET,
      redirect_uri: FACEBOOK_OAUTH_REDIRECT_URL,
      scopes: ['email'],
      state: from,
    };
  
    const qs = new URLSearchParams(options);
  
    return `${rootURl}?${qs.toString()}`;
  }