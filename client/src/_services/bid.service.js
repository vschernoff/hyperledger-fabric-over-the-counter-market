import * as apiService from './api.service';
import {configService} from './config.service';

export const bidService = {
  getAll,
  add,
  edit,
  accept,
  cancel
};

// function _selectChannelFromDeal(deal) {
//   return _selectChannel(configService.get().org, bid.value.creator);
// }
// function _selectChannelFromBid(bid) {
//   const {requestSender, requestReceiver} = bid.key;
//   return _selectChannel(requestSender, requestReceiver);
// }
// function _selectChannel(org1, org2) {
//   const [first, second] = [org1, org2].sort();
//   return apiService.channels[`${first}${second}`];
// }
function _getChannels() {
  const userOrg = configService.get().org;
  return Object.entries(apiService.channels)
    .filter(([k, v]) => {
      return k.includes(userOrg) && v.includes('-');
    })
    .map(([k, v]) => v);
}

function getAll() {
  return apiService.query(
    apiService.channels.common,
    apiService.contracts.reference,
    'queryBids',
    `[]`);
}

function add(bid) {
  return apiService.invoke(
    apiService.channels.common,
    apiService.contracts.reference,
    'placeBid',
    [bid.type + '', bid.amount, bid.rate]
  );
}

function edit(bid) {
  return apiService.invoke(
    apiService.channels.common,
    apiService.contracts.reference,
    'editBid',
    [bid.id, 0, bid.amount, bid.rate]
  );
}

function accept(bid) {
  return apiService.invoke(
    apiService.channels.common,
    apiService.contracts.reference,
    'makeDeal',
    [bid.key.id]
  );
}

function cancel(bid) {
  return apiService.invoke(
    apiService.channels.common,
    apiService.contracts.reference,
    'cancelBid',
    [bid.key.id]
  );
}
