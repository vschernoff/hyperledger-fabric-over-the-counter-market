import {authService} from '../_services/auth.service';
import {getToken} from './user-store';

const forwardRequestByError = {
  401: {
    retry: 3,
    timeout: 1000,
    fnc: authService.obtainToken.bind(authService)
  },
  500: {
    retry: 10,
    timeout: 3000,
    fnc: Promise.resolve.bind(Promise)
  }
};

function _parseMessage(input = '') {
  const [, detailedMsg] = input.replace(')', '').split('message: ');
  return detailedMsg;
}

export function sendRequest(url, options) {
  options.headers = updateHeaders(options.headers);
  return fetch(url, options)
    .then(handleResponse(url, options));
}

const resendRequest = (url, options) => {
  options.headers = updateHeaders(options.headers);
  return fetch(url, options)
    .then((response) => {
        return response.json().then(data => (response.ok) ? Promise.resolve(data) : Promise.reject(response.ok))
      }
    );
};

const retrySend = (url, options, timeout, retry) => {
  return new Promise((resolve, reject) => {
    if (retry > 0) {
      resendRequest(url, options).then(
        resolve,
        () => {
          setTimeout(() => {
            retry--;
            return retrySend(url, options, timeout, retry)
              .then(resolve)
              .catch(reject);
          }, timeout);
        }
      );
    } else {
      reject("the steps are over");
    }
  });
};

function handleResponse(url, options) {
  return (response) => {
    return response.json().then(data => {
      if (!response.ok) {
        if (forwardRequestByError[response.status]) {
          data = forwardRequestByError[response.status].fnc().then(() => {
            return retrySend(url, options, forwardRequestByError[response.status].timeout, forwardRequestByError[response.status].retry)
              .then(
                result => {
                  return result;
                },
                error => {
                  return Promise.reject(new Error(error));
                }
              );
          });
          return data;
        } else {
          const error = (data && data.message && _parseMessage(data.message)) || response.statusText;
          return Promise.reject(new Error(error));
        }
      }
      return data;
    });
  };
}

function authHeader() {
  let token = getToken();

  if (token) {
    return {'Authorization': 'Bearer ' + token};
  } else {
    return {};
  }
}

function updateHeaders(headerObj = {}) {
  return {...headerObj, ...authHeader(), ...{'Content-Type': 'application/json'}}
}
