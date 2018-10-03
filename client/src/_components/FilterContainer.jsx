import React from 'react';
import {connect} from 'react-redux';
import moment from 'moment';
import DatePicker from 'react-datepicker';
import CheckBox from './secondary/CheckBox';
import 'react-datepicker/dist/react-datepicker.css';

const components = {
  DatePicker,
  CheckBox
};

class FilterContainer extends React.Component {
  constructor(props) {
    super(props);
    const {propertyFilterContainer, handleSubmit} = this.props;

    const propertyFilterContainerDefault = [
      {
        type: "DatePicker",
        label: "From",
        state: {nameProp: "selected", name: "fromDate", value: moment()},
        properties: {
          isClearable: "true"
        }
      },
      {
        type: "DatePicker",
        label: "To",
        state: {nameProp: "selected", name: "toDate", value: moment()},
        properties: {
          isClearable: "true"
        }
      },
      {
        type: "CheckBox",
        label: "",
        state: {nameProp: "checked", name: "myDeals", value: false},
        properties: {
          label: "Only my Deals"
        }
      }
    ];

    this.instrumentsNames = (propertyFilterContainer || propertyFilterContainerDefault);

    let states = {};

    this.instruments = this.instrumentsNames.map(
      element => {
        if (typeof components[element.type] !== "undefined") {
          states[element.state.name] = element.state.value;
          return components[element.type]
        } else {
          throw new Error(`There is no such Component ${element.type}`);
        }
      }
    );
    this.state = states;

    this.handleSubmit = handleSubmit.bind(this);
  }

  handleChange(state, value) {
    value = (value && typeof value !== "undefined") ? (typeof value.target !== "undefined" ? (value.target.type === 'checkbox' ? value.target.checked : value) : value) : value;

    this.setState({
      [state]: value
    });
  }

  prepareSubmit(event) {
    let parameters = {};
    this.instrumentsNames.map(
      element => {
        return parameters[element.state.name] = moment.isMoment(this.state[element.state.name]) ? this.state[element.state.name].format('X') : this.state[element.state.name];
      }
    );
    this.handleSubmit(event, parameters);
  }

  render() {
    return (
      <form onSubmit={this.prepareSubmit.bind(this)}>
        {this.instruments.map((Instrument, index) => {
          let proporties = this.instrumentsNames[index].properties;
          proporties[this.instrumentsNames[index].state.nameProp] = this.state[this.instrumentsNames[index].state.name];

          return (
            <div className="form-group row">
              <label htmlFor={"input" + index}
                     className="col-sm-2 col-form-label">{this.instrumentsNames[index].label}</label>
              <div className="col-sm-10">
                <Instrument
                  id={"input" + index}
                  onChange={this.handleChange.bind(this, this.instrumentsNames[index].state.name)}
                  {...proporties}
                />
              </div>
            </div>
          );
        })}
        <div className="form-group row">
          <div className="col-sm-10">
            <button type="submit" className="btn btn-primary">Apply</button>
          </div>
        </div>
      </form>
    );
  }
}

function mapStateToProps(state) {
  const {deals} = state;
  return {
    deals
  };
}

const connected = connect(mapStateToProps)(FilterContainer);
export {connected as FilterContainer};
