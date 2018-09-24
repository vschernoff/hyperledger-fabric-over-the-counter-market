import {authService} from '../_services/auth.service';
import {getToken} from './user-store';

const forwardRequestByError = {
  401: {
    retry: 2,
    timeout: 1000,
    fnc: authService.obtainToken.bind(authService)
  },
  500: {
    retry: 3,
    timeout: 3000,
    fnc: Promise.resolve.bind(Promise)
  }
};

function _parseMessage(input = '') {
  const [, detailedMsg] = input.replace(')', '').split('message: ');
  return detailedMsg;
}

export function sendRequest(url, options = {}, retry) {
  options.headers = updateHeaders(options.headers);
  return fetch(url, options)
    .then(handleResponse(url, options, retry));
}

function retry(asyncFn, args, timeout, attempt, resolveFn, rejectFn, e) {
  if (attempt < 0) {
    return rejectFn(e);
  }
  return asyncFn.apply(this, args)
    .then(resolveFn)
    .catch(e => {
      attempt--;
      setTimeout(() => {
        retry(asyncFn, args, timeout, attempt, resolveFn, rejectFn, e);
      }, timeout);
    });
}

const retryPromise = (url, options, timeout, attempt) => {
  return new Promise((resolve, reject) => {
    retry(sendRequest, [url, options, false], timeout, attempt, resolve, reject)
  });
};

function handleResponse(url, options, retry = true) {
  return (response) => {
    return response.json().then(data => {
      if (!response.ok) {
        const retryParams = forwardRequestByError[response.status];
        if (retryParams && retry) {
          return retryParams.fnc().then(() => {
            return retryPromise(
              url, options,
              retryParams.timeout,
              retryParams.retry);
          });
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
