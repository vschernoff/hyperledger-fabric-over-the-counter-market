import React from 'react';
import {connect} from 'react-redux';
import Switch from 'rc-switch';
import Creatable from 'react-select/lib/Creatable';
import {components} from 'react-select';

import {orderActions} from '../_actions';
import {commonConstants} from '../_constants/common.constants';

const rateOptions = [
  {value: '1.000', label: '1.000'},
  {value: '0.875', label: '0.875'},
  {value: '0.750', label: '0.750'},
  {value: '0.625', label: '0.625'},
  {value: '0.500', label: '0.500'},
  {value: '0.375', label: '0.375'},
  {value: '0.250', label: '0.250'},
  {value: '0.125', label: '0.125'}
];

class AddOrder extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      bid: {
        amount: undefined,
        rate: undefined,
        type: 1,
        created: false
      },
      submitted: false
    };
    this._fill();

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.changeType = this.changeType.bind(this);
    this.changeRate = this.changeRate.bind(this);

    //Modal integration
    this.props.setSubmitFn && this.props.setSubmitFn(this.handleSubmit);
  }

  componentDidMount() {
    this.amountInput.focus();
  }

  _fill() {
    if(this.props.initData && this.props.initData.key) {
      this.state.bid.id = this.props.initData.key.id;
      this.state.bid.amount = this.props.initData.value.amount;
      this.state.bid.rate = this.props.initData.value.rate;
      this.state.bid.created = true;
    }
  }

  changeType(newV) {
    this.setState({
      bid: {
        ...this.state.bid,
        type: newV === true ? 0 : 1
      }
    });
  }

  changeRate(newV) {
    this.setState({
      bid: {
        ...this.state.bid,
        rate: (newV && newV.value) || ''
      }
    });
  }

  handleChange(event) {
    const {name, value} = event.target;
    const {bid} = this.state;
    this.setState({
      bid: {
        ...bid,
        [name]: value
      },
      submitted: false
    });
  }

  handleSubmit(event) {
    event.preventDefault();

    this.setState({submitted: true});
    const {bid} = this.state;
    if (bid.amount && bid.rate) {
      this.props.dispatch(orderActions[bid.id ? 'edit' : 'add'](bid));
    }
  }


  render() {
    const {bid, submitted} = this.state;
    return (
      <form name="form" onSubmit={this.handleSubmit}>
        <div className='form-group'>
          <label htmlFor="amount">Amount</label>

          <div className='input-group'>
            <div className="input-group-prepend">
              <span className="input-group-text">{commonConstants.CURRENCY_SIGN}</span>
            </div>
            <input type="text" className={"form-control" + (submitted && !bid.amount ? ' is-invalid' : '')}
                      name="amount" value={bid.amount}
                      ref={(input) => { this.amountInput = input; }}
                      onChange={this.handleChange}/>
            {submitted && !bid.amount && <div className="text-danger invalid-feedback">Amount is required</div>}
          </div>
        </div>
        <Creatable
          name="rate"
          isClearable
          onChange={this.changeRate}
          placeholder="Select or type any value..."
          options={rateOptions}
          value={bid.rate ? {value: bid.rate, label: bid.rate} : undefined}
          components={{ Control: (props) => (
            <div className='form-group'>
              <label htmlFor="rate">Rate</label>
              <div className="input-group">
                <div className="input-group-prepend">
                  <span className="input-group-text">{commonConstants.RATE_SIGN}</span>
                </div>
                <components.Control
                  className={"form-control" + (submitted && !bid.rate ? ' is-invalid' : '')} {...props} />
                {submitted && !bid.rate && <div className="text-danger invalid-feedback">Rate is required</div>}
              </div>
            </div>
          ) }}
        />

        {!bid.id && <div>
          <Switch
            name="type"
            onChange={this.changeType}
            checkedChildren={'Lend'}
            unCheckedChildren={'Borrow'}
          />
        </div>}
      </form>
    );
  }
}

function mapStateToProps(state) {
  const {bid} = state;

  return {
    bid
  }
}

const connected = connect(mapStateToProps)(AddOrder);
export {connected as AddOrder};