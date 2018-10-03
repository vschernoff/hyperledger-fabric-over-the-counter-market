import React from 'react';
import ReactTable from 'react-table';
import {connect} from 'react-redux';

import {dealActions} from '../_actions';
import {formatter} from '../_helpers';
import {commonConstants} from '../_constants';


class DealsTable extends React.Component {
  constructor() {
    super();

    this.refreshData = this.refreshData.bind(this);
  }

  componentDidMount() {
    this.refreshData();
  }

  refreshData() {
    this.props.dispatch(dealActions.getAll());
  }

  render() {
    const {deals, columns} = this.props;

    if (!deals || !deals.items) {
      return null;
    }

    const columnsDefault = [{
      Header: 'Lender',
      id: 'value.lender',
      accessor: rec => formatter.org(rec.value.lender)
    }, {
      Header: 'Borrower',
      id: 'value.borrower',
      accessor: rec => formatter.org(rec.value.borrower)
    }, {
      Header: 'Amount, ' + commonConstants.CURRENCY_SIGN,
      id: 'value.amount',
      accessor: rec => formatter.number(rec.value.amount)
    }, {
      Header: 'Rate, ' + commonConstants.RATE_SIGN,
      id: 'value.rate',
      accessor: rec => formatter.rate(rec.value.rate)
    }, {
      Header: 'Date',
      id: 'value.timestamp',
      accessor: rec => formatter.time(new Date(rec.value.timestamp * 1000))
    }];

    return (
        <ReactTable
          columns={columns || columnsDefault}
          data={deals.items || []}
          className="-striped -highlight"
          defaultPageSize={10}
          minWidth={60}
          filterable={true}
          defaultSorted={[
            {
              id: "value.timestamp",
              desc: true
            }
          ]}
        />
    );
  }
}

function mapStateToProps(state) {
  const {deals} = state;
  return {
    deals
  };
}

const connected = connect(mapStateToProps)(DealsTable);
export {connected as DealsTable};
