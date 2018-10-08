import React from 'react';
import ReactTable from 'react-table';
import {connect} from 'react-redux';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faHandshake, faPen, faTimes } from '@fortawesome/free-solid-svg-icons';

import {orderActions} from '../_actions';
import {modalIds} from '../_constants';
import {Modal} from '../_components';
import {formatter} from '../_helpers';
import {commonConstants} from '../_constants/common.constants';

class OrdersTable extends React.Component {
  constructor() {
    super();

    this.refreshData = this.refreshData.bind(this);
    this.makeDeal = this.makeDeal.bind(this);
    this.cancelBid = this.cancelBid.bind(this);
  }

  componentDidMount() {
    this.refreshData();
  }

  refreshData() {
    this.props.dispatch(orderActions.getAll());
  }

  render() {
    const {orders, user} = this.props;

    if(!orders || !orders.items) {
      return null;
    }
    //TODO use apporpriate API
    orders.items = orders.items.filter(i => i.value.status === 0);

    const columns = [{
      Header: 'Lender',
      id: 'lender',
      accessor: rec => rec.value.type === 0 ? formatter.org(rec.value.creator) : ''
    }, {
      Header: 'Amount Lend, ' + commonConstants.CURRENCY_SIGN,
      id: 'amount.lend',
      className: 'text-danger',
      accessor: rec => rec.value.type === 0 ? formatter.number(rec.value.amount) : ''
    }, {
      Header: 'Rate, ' + commonConstants.RATE_SIGN,
      id: 'value.rate',
      accessor: rec => formatter.rate(rec.value.rate),
      sortMethod: (a, b) => {
        return a && b && parseFloat(a) > parseFloat(b) ? 1 : -1;
      }
    }, {
      Header: 'Amount Borrow, ' + commonConstants.CURRENCY_SIGN,
      id: 'amount.borrow',
      className: 'text-success',
      accessor: rec => rec.value.type === 1 ? formatter.number(rec.value.amount) : ''
    },{
      Header: 'Borrower',
      id: 'borrower',
      accessor: rec => rec.value.type === 1 ? formatter.org(rec.value.creator) : ''
    },  {
      id: 'actions',
      Header: 'Actions',
      accessor: 'key.id',
      filterable: false,
      Cell: row => {
        const record = row.original;
        return (
          <div>
            {(record.value.status === 0 && record.value.creator !== user.org) &&
              (<button className="btn btn-sm btn-success" title="Make deal"
                       onClick={()=>{this.makeDeal(row.original)}}>
                <FontAwesomeIcon fixedWidth icon={faHandshake}/>
              </button>)
            }
            {(record.value.status === 0 && record.value.creator === user.org) &&
            (<button className="btn btn-sm btn-primary"  title="Edit"
                     onClick={Modal.open.bind(this, modalIds.editOrder, record)}>
              <FontAwesomeIcon fixedWidth icon={faPen}/>
            </button>)
            }
            {(record.value.status === 0 && record.value.creator === user.org) &&
              (<button className="btn btn-sm btn-danger"  title="Cancel"
                       onClick={()=>{this.cancelBid(row.original)}}>
                <FontAwesomeIcon fixedWidth icon={faTimes}/>
              </button>)
            }
          </div>
        )
      }
    }];

    return (
      <ReactTable
        columns={columns}
        data={orders.items || []}
        className="-striped -highlight"
        defaultPageSize={10}
        minWidth={60}
        filterable={true}
        getTrProps={(state, rowInfo) => {
          return {
            style: {
              'font-weight': rowInfo && rowInfo.original.value.creator !== user.org ? "normal" : "bold"
            }
          };
        }}
        defaultSorted={[
          {
            id: "value.rate",
            desc: true
          }
        ]}/>
    );
  }

  makeDeal(record) {
    this.props.dispatch(orderActions.accept(record));
  }

  cancelBid(record) {
    this.props.dispatch(orderActions.cancel(record));
  }

  editBid(record) {
    this.props.dispatch(orderActions.edit(record));
  }
}

function mapStateToProps(state) {
  const {orders, authentication} = state;
  const {user} = authentication;
  return {
    orders,
    user
  };
}

const connected = connect(mapStateToProps)(OrdersTable);
export {connected as OrdersTable};
