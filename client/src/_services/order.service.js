// @flow
import * as apiService from './api.service';
import {configService} from './config.service';
import type {Order} from '../_types';
export const orderService = {
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

async function add(order: Order) {
  const channels = await configService.getChannels();
  const chaincodes = await configService.getChaincodes();
  return apiService.invoke(
    channels[CHANNELS.common],
    chaincodes[CHAINCODES.reference],
    ACTIONS.add,
    [order.value.type, order.value.amount, order.value.rate]
  );
}

async function edit(order: Order) {
  const channels = await configService.getChannels();
  const chaincodes = await configService.getChaincodes();
  return apiService.invoke(
    channels[CHANNELS.common],
    chaincodes[CHAINCODES.reference],
    ACTIONS.edit,
    [order.key.id, 0, order.value.amount, order.value.rate]
  );
}

async function accept(order: Order) {
  const channels = await configService.getChannels();
  const chaincodes = await configService.getChaincodes();
  return apiService.invoke(
    channels[CHANNELS.common],
    chaincodes[CHAINCODES.reference],
    ACTIONS.accept,
    [order.key.id]
  );
}

async function cancel(order: Order) {
  const channels = await configService.getChannels();
  const chaincodes = await configService.getChaincodes();
  return apiService.invoke(
    channels[CHANNELS.common],
    chaincodes[CHAINCODES.reference],
    ACTIONS.cancel,
    [order.key.id]
  );
}
