import React from 'react';
import {connect} from 'react-redux';

import {Chart, DealsTable, FilterContainer} from '../_components';
import {dealActions} from '../_actions';
import {formatter} from '../_helpers';
import {commonConstants} from '../_constants';
import moment from 'moment';

const parametersMap = {
  from: "fromDate",
  to: "toDate",
  creator: "myDeals"
};

const columns = [{
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
  Header: 'Date & Time',
  id: 'value.timestamp',
  accessor: rec => formatter.datetime(new Date(rec.value.timestamp * 1000))
}];

const propertyFilterContainer = [
  {
    type: "DatePicker",
    label: "From",
    state: {nameProp: "selected", name: parametersMap.from, value: moment()},
    properties: {
      isClearable: "true"
    }
  },
  {
    type: "DatePicker",
    label: "To",
    state: {nameProp: "selected", name: parametersMap.to, value: moment()},
    properties: {
      isClearable: "true"
    }
  },
  {
    type: "CheckBox",
    label: "",
    state: {nameProp: "checked", name: parametersMap.creator, value: false},
    properties: {
      label: "Only my Deals"
    }
  }
];

class DealsPage extends React.Component {

  componentDidMount() {
    this.timerId = setInterval(this.refreshData.bind(this, true), 50000);
  }

  componentWillUnmount() {
    clearInterval(this.timerId);
  }

  refreshData(shadowMode) {
    this.props.dispatch(dealActions.getAll(shadowMode));
  }

  handleSubmit(event, parameters) {
    event.preventDefault();

    Object.keys(parameters).map(function(key) {
      parameters[key] = parameters[key] === null ? '0'  : parameters[key];
    });

    parameters[parametersMap.creator] ?
      this.props.dispatch(dealActions.getForCreatorByPeriod([parameters[parametersMap.from], parameters[parametersMap.to]]))
      : this.props.dispatch(dealActions.getByPeriod([parameters[parametersMap.from], parameters[parametersMap.to]]));
  }

  render() {
    return (
      <div>
        <div className="row">
          <div className="col">
            <FilterContainer propertyFilterContainer={propertyFilterContainer}
                             handleSubmit={this.handleSubmit.bind(this)}/>
          </div>
        </div>
        <div className="row">
          <div className="col">
            <DealsTable columns={columns}/>
          </div>
        </div>
        <div className="row">
          <div className="col">
            <Chart formatXAxis={formatter.datetime} lineType="liner"/>
          </div>
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  const {deals, authentication} = state;
  const {user} = authentication;
  return {
    user,
    deals
  };
}

const connected = connect(mapStateToProps)(DealsPage);
export {connected as DealsPage};