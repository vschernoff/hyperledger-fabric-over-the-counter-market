import * as apiService from './api.service';
import {configService} from './config.service';

export const dealService = {
  getAll,
  getByPeriod,
  getForCreatorByPeriod
};

const ACTIONS = {
  getAll: 'queryDeals',
  getByPeriod: 'queryDealsByPeriod',
  getForCreatorByPeriod: 'queryDealsForCreatorByPeriod'
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
  return await apiService.query(
    channels[CHANNELS.common],
    chaincodes[CHAINCODES.reference],
    ACTIONS.getAll);
}

async function getByPeriod(period) {
  const channels = await configService.getChannels();
  const chaincodes = await configService.getChaincodes();
  return await apiService.query(
    channels[CHANNELS.common],
    chaincodes[CHAINCODES.reference],
    ACTIONS.getByPeriod,
    period);
}

async function getForCreatorByPeriod(period) {
  const channels = await configService.getChannels();
  const chaincodes = await configService.getChaincodes();
  return await apiService.query(
    channels[CHANNELS.common],
    chaincodes[CHAINCODES.reference],
    ACTIONS.getForCreatorByPeriod,
    period);
}
