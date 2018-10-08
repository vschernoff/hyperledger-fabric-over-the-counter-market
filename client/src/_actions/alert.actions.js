// @flow
import {alertConstants} from '../_constants';

const TIMEOUT = 5000;
export const alertActions = {
  success,
  error,
  clear
};

function success(message: string) {
  return (dispatch: Function) => {
    setTimeout(() => {dispatch(clear())}, TIMEOUT);
    return dispatch({type: alertConstants.SUCCESS, message});
  };
}

function error(message: string) {
  return (dispatch: Function) => {
    setTimeout(() => {dispatch(clear())}, TIMEOUT);
    return dispatch({type: alertConstants.ERROR, message});
  };
}

function clear() {
  return {type: alertConstants.CLEAR};
}