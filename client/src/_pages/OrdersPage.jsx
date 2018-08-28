import React from 'react';
import {connect} from 'react-redux';

import {orderActions, modalActions} from '../_actions';
import {AddOrder, Modal, OrdersTable} from '../_components';
import {modalIds} from '../_constants';


class OrdersPage extends React.Component {
  constructor() {
    super();

    this.refreshData = this.refreshData.bind(this);
  }

  componentDidUpdate(prevProps) {
    const {orders, dispatch} = this.props;

    if (orders.adding === false) {
      dispatch(modalActions.hide(modalIds.addOrder));
      this.refreshData();
    }
  }

  refreshData() {
    this.props.dispatch(orderActions.getAll());
  }

  render() {
    return (
      <div>
        <button className="btn btn-primary" onClick={Modal.open.bind(this, modalIds.addOrder)}>
          Place Order
        </button>
        <hr/>
        <h3>All orders <button className="btn" onClick={this.refreshData}><i className="fas fa-sync"/></button></h3>
        <Modal modalId={modalIds.addOrder} title="Place Order">
          <AddOrder/>
        </Modal>
        <OrdersTable/>
      </div>
    );
  }
}

function mapStateToProps(state) {
  const {orders} = state;
  return {
    orders
  };
}

const connected = connect(mapStateToProps)(OrdersPage);
export {connected as OrdersPage};