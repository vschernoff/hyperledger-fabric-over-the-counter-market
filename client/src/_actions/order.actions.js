import {orderConstants} from '../_constants';
import {orderService} from '../_services';
import {alertActions} from './';

export const orderActions = {
  getAll,
  add,
  edit,
  accept,
  cancel,
  history
};

function getAll(shadowMode) {
  return dispatch => {
    dispatch(request(shadowMode));
    orderService.getAll()
      .then(
         requests => {
          dispatch(success(requests));
        },
        error => {
          dispatch(failure(error.toString()));
          dispatch(alertActions.error(error.toString()));
        }
      );
  };

  function request(shadowMode) {
    return {type: orderConstants.GET_ALL_REQUEST, shadowMode}
  }

  function success(orders) {
    return {type: orderConstants.GET_ALL_SUCCESS, orders}
  }

  function failure(error) {
    return {type: orderConstants.GET_ALL_FAILURE, error}
  }
}

function add(bid, comment) {
  return dispatch => {
    dispatch(request());
    orderService.add(bid, comment)
      .then(
        _ => {
          dispatch(success());
          dispatch(alertActions.success('Order was added'));
        },
        error => {
          dispatch(failure(error.toString()));
          dispatch(alertActions.error(error.toString()));
        }
      );
  };

  function request() {
    return {type: orderConstants.ADD_REQUEST}
  }

  function success() {
    return {type: orderConstants.ADD_SUCCESS}
  }

  function failure(error) {
    return {type: orderConstants.ADD_FAILURE, error}
  }
}

function edit(bid) {
  return dispatch => {
    dispatch(request());
    orderService.edit(bid)
      .then(
        _ => {
          dispatch(success());
          dispatch(alertActions.success('Order was updated'));
        },
        error => {
          dispatch(failure(error.toString()));
          dispatch(alertActions.error(error.toString()));
        }
      );
  };

  function request() {
    return {type: orderConstants.EDIT_REQUEST}
  }

  function success() {
    return {type: orderConstants.EDIT_SUCCESS}
  }

  function failure(error) {
    return {type: orderConstants.EDIT_FAILURE, error}
  }
}

function accept(req) {
  return dispatch => {
    dispatch(request());
    orderService.accept(req)
      .then(
        _ => {
          dispatch(success());
          dispatch(alertActions.success('Order was accepted'));
        },
        error => {
          dispatch(failure(error.toString()));
          dispatch(alertActions.error(error.toString()));
        }
      );
  };

  function request() {
    return {type: orderConstants.ACCEPT_REQUEST}
  }

  function success() {
    return {type: orderConstants.ACCEPT_SUCCESS}
  }

  function failure(error) {
    return {type: orderConstants.ACCEPT_FAILURE, error}
  }
}

function cancel(req) {
  return dispatch => {
    dispatch(request());
    orderService.cancel(req)
      .then(
        _ => {
          dispatch(success());
          dispatch(alertActions.success('Order was cancelled'));
        },
        error => {
          dispatch(failure(error.toString()));
          dispatch(alertActions.error(error.toString()));
        }
      );
  };

  function request() {
    return {type: orderConstants.CANCEL_REQUEST}
  }

  function success() {
    return {type: orderConstants.CANCEL_SUCCESS}
  }

  function failure(error) {
    return {type: orderConstants.CANCEL_FAILURE, error}
  }
}

function history(req) {
  return dispatch => {
    dispatch(request());

    orderService.history(req)
      .then(
        history => {
          dispatch(success(req, history));
        },
        error => dispatch(failure(error.toString()))
      );
  };

  function request() {
    return {type: orderConstants.HISTORY_REQUEST};
  }

  function success(req, history) {
    return {type: orderConstants.HISTORY_SUCCESS, req, history};
  }

  function failure(error) {
    return {type: orderConstants.HISTORY_FAILURE, error};
  }
}