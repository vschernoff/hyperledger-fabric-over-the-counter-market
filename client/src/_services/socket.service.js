export const socketService = {
  subscribe,
  unsubscribe
};
let socket;
let intervalID;

function subscribe(callback) {
  const reference_CHAIN_CHAINCODE = 'reference';

  socket = new WebSocket(`ws://${window.location.host}/api/notifications`);

  const subscribeMessage = chaincode => JSON.stringify({
    event: 'block',
    channel: 'common',
    cc_event: 'cc_event',
    chaincode
  });

  const heartbeat = () => {
    intervalID = setInterval(() => {
      socket.send(JSON.stringify('ping'));
    }, 10000);
  };

  socket.onopen = function () {
    socket.send(subscribeMessage(reference_CHAIN_CHAINCODE));
    heartbeat();
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
  intervalID && clearInterval(intervalID);
}
