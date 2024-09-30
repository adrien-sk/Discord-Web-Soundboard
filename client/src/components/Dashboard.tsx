import reactLogo from '../assets/logo_soundboard.svg'

function Dashboard() {
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

                <a href="https://github.com/adrien-sk/Discord-Web-Soundboard" target="_blank"><span className="git-link"><i className="fab fa-github"></i></span></a>
            </div>
        </main>
    )
}

export default Dashboard
