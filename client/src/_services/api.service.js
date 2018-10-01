import {sendRequest} from '../_helpers/request';
import {configService} from './config.service';

function _getPeer() {
  return configService.getPeers()[0];
}

export function query(channel, chaincode, fcn, args) {
  const requestOptions = {
    method: 'GET'
  };

  const {org} = configService.get();
  const params = {
    peer: `${org}/${_getPeer()}`,
    fcn,
    args
  };

  const url = new URL(`${window.location.origin}/channels/${channel}/chaincodes/${chaincode}`);
  Object.keys(params).forEach(key => url.searchParams.append(key, params[key]));

  return sendRequest(`${url.pathname}${url.search}`, requestOptions);
}

export function invoke(channel, chaincode, functionName, args) {
  const {org} = configService.get();
  args = args.map(arg => arg + '');
  const requestOptions = {
    method: 'POST',
    body: JSON.stringify({
      peers: [`${org}/${_getPeer()}`],
      fcn: functionName,
      args
    })
  };

  return sendRequest(`/channels/${channel}/chaincodes/${chaincode}`, requestOptions);
}

export function login(user, retry = true) {
  const requestOptions = {
    method: 'POST',
    body: JSON.stringify({username: user.name, orgName: user.org})
  };

  return sendRequest(`/users`, requestOptions, retry);
}

export function getChannels() {
  return sendRequest(`/channels?peer=${_getPeer()}`);
}

export function getChaincodes() {
  return sendRequest(`/chaincodes?peer=${_getPeer()}`);
}

export function extendConfig() {
  return configService.update();
}

export function config() {
  const requestOptions = {
    method: 'GET'
  };

  return sendRequest(`/config`, requestOptions);
}
