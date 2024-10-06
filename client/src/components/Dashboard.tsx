import { useEffect, useState } from 'react'
import reactLogo from '../assets/logo_soundboard.svg'

function Dashboard() {
    const [webSocketReady, setWebSocketReady] = useState(false);
    const [webSocket, setWebSocket] = useState(new WebSocket("http://localhost:8080/ws"));

    useEffect(() => {
        webSocket.onopen = (event) => {
            console.log('On Open')
            //setWebSocketReady(true);
        };

        webSocket.onmessage = function (event) {
            console.log('Server message : '+ JSON.parse(event.data))
        };

        webSocket.onclose = function (event) {
            console.log('On Close')
            //setWebSocketReady(false);
            // setTimeout(() => {
            //     setWebSocket(new WebSocket("http://localhost:8080/ws"));
            // }, 1000);
        };

        webSocket.onerror = function (err) {
            console.log('Socket encountered error: ', err, 'Closing socket');
            //setWebSocketReady(false);
            webSocket.close();
        };

        return () => {
            if (webSocket.readyState === 1) { // <-- This is important
                webSocket.close();
            }
        };
    }, []);

    function playSound(){
        webSocket.send("Play Sound please");
    }

    return (
        <main>
            <div className="left-panel">
                <header>
                    <h1>Discord Web Soundboard</h1>
                    {/* <div className="logo"><img src={reactLogo} /></div> */}
                </header>
            </div>
            <div id='page-container'>
                <h2>Dashboard Content</h2>
                <button onClick={playSound}>Play Sound</button>

                <a href="https://github.com/adrien-sk/Discord-Web-Soundboard" target="_blank"><span className="git-link"><i className="fab fa-github"></i></span></a>
            </div>
        </main>
    )
}

export default Dashboard
