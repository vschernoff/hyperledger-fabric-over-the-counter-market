export const socketService = {
  subscribe,
  unsubscribe
};
let socket;

function subscribe(callback) {
  const reference_CHAIN_CHAINCODE = 'reference';

  socket = new WebSocket(`ws://${window.location.host}/api/notifications`);

  const subscribeMessage = chaincode => JSON.stringify({
    event: 'block',
    channel: 'common',
    cc_event: 'cc_event',
    chaincode
  });

  socket.onopen = function () {
    socket.send(subscribeMessage(reference_CHAIN_CHAINCODE));
  };

  socket.onmessage = function (event) {
    let eventData = JSON.parse(event.data);
    console.log(eventData.event_type, eventData.event_type === "block");
    if (eventData.event_type === "block") {
      callback();
    }
  }
}

function unsubscribe() {
  socket.close();
}
