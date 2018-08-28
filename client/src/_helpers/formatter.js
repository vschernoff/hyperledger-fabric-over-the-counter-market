import {orgConstants} from '../_constants';

export const formatter = {
  rate: val => Number(val).toFixed(3),
  number: val => Number(val).toLocaleString(),
  date: date => date.toLocaleDateString(),
  time: date => date.toLocaleTimeString(),
  org: orgKey => orgConstants[orgKey]
};