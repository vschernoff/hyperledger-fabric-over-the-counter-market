import {set, get, clear} from '../_helpers';
import {login} from './api.service';

export const authService = {
  obtainToken
};
let fetching = false;
let promise;
function obtainToken() {
  const user = get();
  if (!user) {
    fetching = false;
    clear();
    window.location.reload(true);
    return Promise.resolve();
  }
  if (fetching) {
    return promise;
  }
  fetching = true;
  promise = login(user, false)
    .then(res => {
      fetching = false;
      if (res.token) {
        user.token = res.token;
        set(user);
      }

      return user;
    })
    .catch(e => {
      fetching = false;
      return Promise.reject(e);
    });

  return promise;
}
