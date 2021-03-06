import React from 'react';
import ReactTable from 'react-table';
import {connect} from 'react-redux';

import {dealActions} from '../_actions';
import {formatter} from '../_helpers';
import {commonConstants} from '../_constants';


class SummaryPage extends React.Component {
  constructor() {
    super();

    this.refreshData = this.refreshData.bind(this);
  }

  componentDidMount() {
    this.refreshData();
  }

  componentDidUpdate(prevProps) {
    const {deals} = this.props;

    if (deals.adding === false) {
      this.refreshData();
    }
  }

  refreshData() {
    this.props.dispatch(dealActions.getAll());
  }

  render() {
    const {deals, user} = this.props;

    if (!deals || !deals.items) {
      return null;
    }

    const data = deals.items.filter(d => d.value.borrower === user.org || d.value.lender === user.org);

    const closingBalance = data.reduce((acc, rec) => {
      const val = rec.value.amount;
      return acc + (rec.value.lender === user.org ? -val : val)
    }, 0);

    const netted = formatter.number(data.reduce((acc, rec) => {
      const val = (rec.value.amount/365*rec.value.rate);
      return acc + (rec.value.lender === user.org ? val : -val)
    }, 0));

    const clearingColumns = [{
      Header: 'Cash out, ' + commonConstants.CURRENCY_SIGN,
      id: 'cashOut',
      className: 'text-danger',
      accessor: rec => rec.value.lender === user.org ? formatter.number(rec.value.amount) : '',
      Footer: (
        formatter.number(data.reduce((acc, rec) => acc + (rec.value.lender === user.org ? rec.value.amount : 0), 0))
      )
    }, {
      Header: 'Cash in, ' + commonConstants.CURRENCY_SIGN,
      id: 'cashIn',
      className: 'text-success',
      accessor: rec => rec.value.borrower === user.org ? formatter.number(rec.value.amount) : '',
      Footer: (
        formatter.number(data.reduce((acc, rec) => acc + (rec.value.borrower === user.org ? rec.value.amount : 0), 0))
      )
    }, {
      Header: 'Counterparty',
      id: 'counterParty',
      accessor: rec => formatter.org(rec.value.borrower === user.org ? rec.value.lender : rec.value.borrower)
    }, {
      Header: commonConstants.CURRENCY_SIGN,
      id: 'total',
      accessor: rec => {
        const withBraces = rec.value.borrower === user.org;
        return `${withBraces ? '(' : ''}${formatter.number(rec.value.amount)}${withBraces ? ')' : ''}`;
      },
      Footer: (
        formatter.number(closingBalance)
      )
    }];

    const interestColumns = [{
      Header: 'Receivable, ' + commonConstants.CURRENCY_SIGN,
      id: 'cashOut',
      accessor: rec => rec.value.lender === user.org ? formatter.number(rec.value.amount/365*rec.value.rate) : '',
      Footer: (
        formatter.number(data.reduce((acc, rec) => acc + (rec.value.lender === user.org ? (rec.value.amount/365*rec.value.rate) : 0), 0))
      )
    }, {
      Header: 'Payable, ' + commonConstants.CURRENCY_SIGN,
      id: 'cashIn',
      accessor: rec => rec.value.borrower === user.org ? formatter.number(rec.value.amount/365*rec.value.rate) : '',
      Footer: (
        formatter.number(data.reduce((acc, rec) => acc + (rec.value.borrower === user.org ? (rec.value.amount/365*rec.value.rate) : 0), 0))
      )
    }, {
      Header: 'Counterparty',
      id: 'counterParty',
      accessor: rec => formatter.org(rec.value.borrower === user.org ? rec.value.lender : rec.value.borrower)
    }, {
      Header: 'Value Date',
      id: 'value.timestamp',
      accessor: rec => formatter.date(new Date((rec.value.timestamp + 24 * 60 * 60) * 1000))
    }];

    return (
      <div>
        <div className="row">
          <div className="col col-form-label">
            <button className="btn" onClick={this.refreshData}>Refresh <i className="fas fa-sync"/></button>
          </div>
        </div>
        <div className="row">
          <div className="col">
            <h3>Aftermarket clearing</h3>
            {!!data.length && <ReactTable
              columns={clearingColumns}
              data={data}
              className="-striped -highlight"
              defaultPageSize={deals.items.length + 1}
              filterable={false}
              showPagination={false}
              defaultSorted={[
                {
                  id: "cashOut",
                  desc: true
                }
              ]}
            />}
            <h5 className="mt-2">Closing balance:
              <span className={closingBalance < 0 ? 'text-danger' : 'text-success'}>
                {` ${formatter.number(closingBalance)} ${commonConstants.CURRENCY_SIGN}`}
              </span>
            </h5>
          </div>

          <div className="col">
            <h3>Interest netting</h3>
            {!!data.length && <ReactTable
              columns={interestColumns}
              data={data}
              className="-striped -highlight"
              defaultPageSize={deals.items.length + 1}
              filterable={false}
              showPagination={false}
              defaultSorted={[
                {
                  id: "cashOut",
                  desc: true
                }
              ]}
            />}
            <h5 className="mt-2">Netted:
              <span className={netted < 0 ? 'text-danger' : 'text-success'}>
                {` ${formatter.number(netted)} ${commonConstants.CURRENCY_SIGN}`}
              </span>
            </h5>
          </div>
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  const {bids, deals, authentication, modals} = state;
  const {user} = authentication;
  return {
    user,
    deals,
    modals
  };
}

const connected = connect(mapStateToProps)(SummaryPage);
export {connected as SummaryPage};
