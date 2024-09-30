import reactLogo from '../assets/logo_soundboard.svg'

function Login() {
    return (
        <div className="login-page">
            <h1>Discord Web Soundboard</h1>
            {/* <img className="logo" src={reactLogo} alt="logo"/> */}
            <div className="login-box">
                <h2>Login</h2>
                <a href="http://localhost:8080/auth" className="btn">Authenticate with Discord</a>
            </div>
        </div>
    )
}

export default Login
