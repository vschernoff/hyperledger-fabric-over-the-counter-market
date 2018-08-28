import React from 'react';
import ReactTable from 'react-table';
import {connect} from 'react-redux';

import {dealActions} from '../_actions';
import {DealsTable} from '../_components';


class DealsPage extends React.Component {
  constructor() {
    super();

    this.refreshData = this.refreshData.bind(this);
  }

  refreshData() {
    this.props.dispatch(dealActions.getAll());
  }

  render() {
    return (
      <div>
        <hr/>
        <h3>All deals <button className="btn" onClick={this.refreshData}><i className="fas fa-sync"/></button></h3>
        <DealsTable/>
      </div>
    );
  }
}

function mapStateToProps(state) {
}

const connected = connect(mapStateToProps)(DealsPage);
export {connected as DealsPage};