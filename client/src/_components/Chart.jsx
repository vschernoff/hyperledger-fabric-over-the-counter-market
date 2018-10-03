import React from 'react';
import {connect} from 'react-redux';
import {ComposedChart, ResponsiveContainer, Line, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend} from 'recharts';

import {formatter} from '../_helpers';
import {commonConstants} from '../_constants';


const DEFAULT_TICKS = [0.125, 0.250, 0.375, 0.500, 0.625, 0.750, 0.875, 1.000];

class LineBarAreaComposedChart extends React.Component {
  render () {
    const {deals, formatXAxis, lineType} = this.props;

    if (!deals || !deals.items) {
      return null;
    }

    const data = deals.items.map(r => {
      r.value.ts = formatXAxis ? formatXAxis(new Date(r.value.timestamp * 1000)) : formatter.time(new Date(r.value.timestamp * 1000));
      return r;
    })
      .sort((a, b) => {
        return a.value.timestamp > b.value.timestamp ? 1 : -1;
      });

    return !!data.length && (
      <ResponsiveContainer width="100%" height={300}>
        <ComposedChart data={data}
                       margin={{top: 20, right: 20, bottom: 20, left: 20}}>
          <CartesianGrid stroke='#f5f5f5'/>
          <XAxis dataKey="value.ts" name="Time"/>
          <YAxis yAxisId="left" tickFormatter={formatter.rate} unit={commonConstants.RATE_SIGN} interval={0} ticks={DEFAULT_TICKS} />
          <YAxis yAxisId="right" tickFormatter={formatter.number} orientation="right" unit={commonConstants.CURRENCY_SIGN}/>
          <Tooltip />
          <Legend />
          <Bar yAxisId="right" dataKey='value.amount' name="Amount" barSize={20} fill='#413ea0' unit={commonConstants.CURRENCY_SIGN} />
          <Line yAxisId="left" dataKey='value.rate' name="Rate" type={lineType || "liner"} stroke='#ff7300' unit={commonConstants.RATE_SIGN} />
        </ComposedChart>
      </ResponsiveContainer>
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

const connected = connect(mapStateToProps)(LineBarAreaComposedChart);
export {connected as Chart};
