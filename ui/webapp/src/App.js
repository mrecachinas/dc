import { useState, useEffect, useRef } from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import { Message } from "rsuite";

import StatusView from "./components/StatusView";
import TaskView from "./components/TaskView";
import NotFoundView from "./components/NotFoundView";
import HistoricalPlotView from "./components/HistoricalPlotView";
import StreamingPlotView from "./components/StreamingPlotView";
import HistoricalTable from "./components/HistoricalTable";
import ActiveTable from "./components/ActiveTable";
import AppNavbar from "./components/AppNavbar";

import "rsuite/dist/styles/rsuite-default.css";

const setTaskField = (fieldname) => (fieldvalue) =>
  fetch("/api/tasks/update", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ [fieldname]: fieldvalue }),
  });

export default function App() {
  const [websocketURL, setWebsocketURL] = useState("");

  const [updateTime, setUpdateTime] = useState(null);

  // `{status, tasks}` returned from websockets
  const [info, setInfo] = useState({
    status: { historical: [], active: [] },
    tasks: [],
  });

  // Any errors from websockets or otherwise (e.g., browser-side)
  const [message, setMessage] = useState(null);
  const webSocket = useRef(null);

  // `useEffect` operates similar to `componentDidMount`
  // and `componentDidUpdate`
  useEffect(() => {
    /* TODO: Run this once the server is written

    fetch("/api/websocket_uri")
      .then((response) => response.json())
      .then((data) => setWebsocketURL(data.websocket_uri))
      .catch((error) => console.error(error));
    
    */
    console.log("Setting websocket URL");
    setWebsocketURL("ws://localhost:8080");
  }, []);

  // After getting the `websocketURL`...
  useEffect(() => {
    console.log(`Websocket URL: ${websocketURL}`);
    try {
      webSocket.current = new WebSocket(websocketURL);
      webSocket.current.onmessage = (message) => {
        /**
         * Messages will look like:
         * {
         *   "action": "UPDATE"|"NEW",
         *   "table": "active"|"historical"|"tasking",
         *   "data": {
         *
         *   }
         * }
         */
        console.log(message);
        const jsonData = JSON.parse(message.data);
        setInfo(jsonData);
        const updateTimeString = `Last Updated: ${new Date().toLocaleString()}`;
        setMessage(updateTimeString);
      };
      webSocket.current.onopen = (event) => {
        const msg = `Connected to ${websocketURL}`;
        console.debug(msg);
        setMessage(msg);
      };
      webSocket.current.onclose = (event) => {
        const msg = `Connection to ${websocketURL} closed`;
        console.debug(msg);
        setMessage(msg);
      };
      webSocket.current.onerror = (event) => {
        setMessage(`ERROR: Connection to ${websocketURL} failed`);
      };
      return () => webSocket.current.close();
    } catch {
      setMessage("ERROR: Using fake data");
      console.log("ERROR!");
      setInfo({
        status: {
          historical: [
            { uuid: "12345", foo: "bar", bar: 123 },
            { uuid: "12346", foo: "bar2", bar: 124 },
            { uuid: "12347", foo: "bar3", bar: 120 },
          ],
          active: [
            { uuid: "12345", foo: "bar", bar: 123 },
            { uuid: "12346", foo: "bar2", bar: 124 },
            { uuid: "12347", foo: "bar3", bar: 120 },
          ],
        },
        tasks: [{ a: 1 }, { b: 2 }],
      });
    }
  }, [websocketURL]);

  const { status, tasks } = info;
  return (
    <Router>
      <div>
        <AppNavbar message={message} updateTime={updateTime} />
        <Switch>
          <Route path="/status" render={() => <StatusView status={status} />} />
          <Route
            path="/tasks"
            render={() => (
              <TaskView tasks={tasks} setFieldCallback={setTaskField} />
            )}
          />
          <Route
            path="/active/:id"
            render={(id) => <StreamingPlotView streamID={id} />}
          />
          <Route
            path="/historical/:id"
            render={(id) => <HistoricalPlotView historicalID={id} />}
          />
          <Route
            path="/active"
            render={() => <ActiveTable activeStatus={status.active} />}
          />
          <Route
            path="/historical"
            render={() => (
              <HistoricalTable historicalStatus={status.historical} />
            )}
          />
          <Route component={NotFoundView} />
        </Switch>
      </div>
    </Router>
  );
}
