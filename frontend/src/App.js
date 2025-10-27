import React, { useState, useEffect } from 'react';
import './App.css';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

function App() {
  const [extensions, setExtensions] = useState([]);
  const [calls, setCalls] = useState([]);
  const [fromExtension, setFromExtension] = useState('1000');
  const [toNumber, setToNumber] = useState('');
  const [status, setStatus] = useState('');

  useEffect(() => {
    fetchExtensions();
    const interval = setInterval(fetchExtensions, 5000);
    return () => clearInterval(interval);
  }, []);

  const fetchExtensions = async () => {
    try {
      const response = await fetch(`${API_URL}/api/extensions`);
      const data = await response.json();
      setExtensions(data);
    } catch (error) {
      console.error('Failed to fetch extensions:', error);
    }
  };

  const makeCall = async (e) => {
    e.preventDefault();
    setStatus('Initiating call...');

    try {
      const response = await fetch(`${API_URL}/api/call/make`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          from: fromExtension,
          to: toNumber,
        }),
      });

      const data = await response.json();
      setCalls([...calls, data]);
      setStatus(`Call initiated: ${data.id}`);
      setToNumber('');

      // Monitor call status
      monitorCall(data.id);
    } catch (error) {
      setStatus(`Error: ${error.message}`);
      console.error('Failed to make call:', error);
    }
  };

  const monitorCall = async (callId) => {
    const checkStatus = async () => {
      try {
        const response = await fetch(`${API_URL}/api/call/status?call_id=${callId}`);
        const data = await response.json();
        
        setCalls(prevCalls =>
          prevCalls.map(call =>
            call.id === callId ? data : call
          )
        );

        if (data.status !== 'ended') {
          setTimeout(checkStatus, 2000);
        }
      } catch (error) {
        console.error('Failed to check call status:', error);
      }
    };

    setTimeout(checkStatus, 2000);
  };

  const hangupCall = async (callId) => {
    try {
      await fetch(`${API_URL}/api/call/hangup`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          call_id: callId,
        }),
      });

      setStatus(`Call ${callId} ended`);
    } catch (error) {
      setStatus(`Error: ${error.message}`);
      console.error('Failed to hangup call:', error);
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>CatiSip - SIP Call Manager</h1>
        <p className="subtitle">FreeSWITCH Integration</p>
      </header>

      <div className="container">
        <div className="panel">
          <h2>Available Extensions</h2>
          <div className="extensions-list">
            {extensions.map((ext, idx) => (
              <div key={idx} className="extension-item">
                <span className="extension-number">{ext.extension}</span>
                <span className={`status ${ext.status}`}>{ext.status}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="panel">
          <h2>Make a Call</h2>
          <form onSubmit={makeCall} className="call-form">
            <div className="form-group">
              <label>From Extension:</label>
              <select
                value={fromExtension}
                onChange={(e) => setFromExtension(e.target.value)}
                className="form-control"
              >
                {extensions.map((ext, idx) => (
                  <option key={idx} value={ext.extension}>
                    {ext.extension}
                  </option>
                ))}
              </select>
            </div>

            <div className="form-group">
              <label>To Number:</label>
              <input
                type="text"
                value={toNumber}
                onChange={(e) => setToNumber(e.target.value)}
                placeholder="Enter phone number or extension"
                className="form-control"
                required
              />
            </div>

            <button type="submit" className="btn btn-primary">
              Make Call
            </button>
          </form>

          {status && <div className="status-message">{status}</div>}
        </div>

        <div className="panel">
          <h2>Active Calls</h2>
          <div className="calls-list">
            {calls.length === 0 ? (
              <p className="no-calls">No active calls</p>
            ) : (
              calls.map((call) => (
                <div key={call.id} className="call-item">
                  <div className="call-info">
                    <div className="call-direction">
                      {call.from} → {call.to}
                    </div>
                    <div className="call-meta">
                      <span className={`call-status status-${call.status}`}>
                        {call.status}
                      </span>
                      <span className="call-id">{call.id}</span>
                    </div>
                  </div>
                  {call.status !== 'ended' && (
                    <button
                      onClick={() => hangupCall(call.id)}
                      className="btn btn-danger"
                    >
                      Hangup
                    </button>
                  )}
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
