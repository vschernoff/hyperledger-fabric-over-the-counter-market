import {dealConstants} from '../_constants';
import {dealService} from '../_services';
import {alertActions} from './';

export const dealActions = {
  add,
  edit,
  getAll,
  history
};

function add(product) {
  return dispatch => {
    dispatch(request());
    dealService.add(product)
      .then(
        product => {
          dispatch(success());
          dispatch(alertActions.success('Product was added'));
        },
        error => {
          dispatch(failure(error.toString()));
          dispatch(alertActions.error(error.toString()));
        }
      );
  };

  function request() {
    return {type: dealConstants.ADD_REQUEST}
  }

  function success() {
    return {type: dealConstants.ADD_SUCCESS}
  }

  function failure(error) {
    return {type: dealConstants.ADD_FAILURE, error}
  }
}

function edit(product) {
  return dispatch => {
    dispatch(request());
    dealService.edit(product)
      .then(
        product => {
          dispatch(success());
          dispatch(alertActions.success('Product was updated'));
        },
        error => {
          dispatch(failure(error.toString()));
          dispatch(alertActions.error(error.toString()));
        }
      );
  };

  function request() {
    return {type: dealConstants.EDIT_REQUEST}
  }

  function success() {
    return {type: dealConstants.EDIT_SUCCESS}
  }

  function failure(error) {
    return {type: dealConstants.EDIT_FAILURE, error}
  }
}

function getAll(shadowMode) {
  return dispatch => {
    dispatch(request(shadowMode));

    dealService.getAll()
      .then(
        products => {
          dispatch(success(products));
        },
        error => dispatch(failure(error.toString()))
      );
  };

  function request(shadowMode) {
    return {type: dealConstants.GETALL_REQUEST, shadowMode}
  }

  function success(deals) {
    return {type: dealConstants.GETALL_SUCCESS, deals}
  }

  function failure(error) {
    return {type: dealConstants.GETALL_FAILURE, error}
  }
}

function history(product) {
  return dispatch => {
    dispatch(request());

    dealService.history(product)
      .then(
        history => {
          dispatch(success(product, history));
        },
        error => dispatch(failure(error.toString()))
      );
  };

  function request() {
    return {type: dealConstants.HISTORY_REQUEST};
  }

  function success(product, history) {
    return {type: dealConstants.HISTORY_SUCCESS, product, history};
  }

  function failure(error) {
    return {type: dealConstants.HISTORY_FAILURE, error};
  }
}