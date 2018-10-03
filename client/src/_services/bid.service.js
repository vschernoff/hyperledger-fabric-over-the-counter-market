import * as apiService from './api.service';
import {configService} from './config.service';

export const bidService = {
  getAll,
  add,
  edit,
  accept,
  cancel
};

const ACTIONS = {
  getAll: 'queryBids',
  add: 'placeBid',
  edit: 'editBid',
  accept: 'makeDeal',
  cancel: 'cancelBid'
};

const CHAINCODES = {
  reference: 'reference'
};

const CHANNELS = {
  common: 'common'
};

async function getAll() {
  const channels = await configService.getChannels();
  const chaincodes = await configService.getChaincodes();
  return apiService.query(
    channels[CHANNELS.common],
    chaincodes[CHAINCODES.reference],
    ACTIONS.getAll);
}

async function add(bid) {
  const channels = await configService.getChannels();
  const chaincodes = await configService.getChaincodes();
  return apiService.invoke(
    channels[CHANNELS.common],
    chaincodes[CHAINCODES.reference],
    ACTIONS.add,
    [bid.type + '', bid.amount, bid.rate]
  );
}

async function edit(bid) {
  const channels = await configService.getChannels();
  const chaincodes = await configService.getChaincodes();
  return apiService.invoke(
    channels[CHANNELS.common],
    chaincodes[CHAINCODES.reference],
    ACTIONS.edit,
    [bid.id, 0, bid.amount, bid.rate]
  );
}

async function accept(bid) {
  const channels = await configService.getChannels();
  const chaincodes = await configService.getChaincodes();
  return apiService.invoke(
    channels[CHANNELS.common],
    chaincodes[CHAINCODES.reference],
    ACTIONS.accept,
    [bid.key.id]
  );
}

async function cancel(bid) {
  const channels = await configService.getChannels();
  const chaincodes = await configService.getChaincodes();
  return apiService.invoke(
    channels[CHANNELS.common],
    chaincodes[CHAINCODES.reference],
    ACTIONS.cancel,
    [bid.key.id]
  );
}
