import { useEffect, useState } from 'react'
import './App.css'
import Login from './components/Login'
import Dashboard from './components/Dashboard'

function App() {
    const [isAuthenticated, setIsAuthenticated] = useState(false)

    useEffect(() => {
        // Call /auth/isloggedin to check if user is logged in
        // Then setup isauthenticated data
        fetch("http://localhost:8080/auth/isauthenticated", {
			method: "GET",
			credentials: "include",
		})
		.then(response => {
            if(response.ok) {
                setIsAuthenticated(true)
            }
            else{
                throw new Error('failed to authenticate user');
            }
		})
		.catch(error => {
			console.log('Error')
			console.log(error)
		});
      }, []);

    return (
        <>
        {
            isAuthenticated ? (
                <Dashboard />
            ) : (
                <Login />
            )
        }
        </>
    )
}

export default App
