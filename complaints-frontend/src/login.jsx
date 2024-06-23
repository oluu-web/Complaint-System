import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loginError, setLoginError] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true); // Set loading state to true

    const response = await fetch('http://localhost:4000/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });

    const data = await response.json();
    setIsLoading(false); // Set loading state to false after response

    if (data.response.ok) {
      console.log("Login successful", data);
      setLoginError(false);
      sessionStorage.setItem("token", data.response.token);
      sessionStorage.setItem("userID", data.response.user_id);
      if (data.response.role === "S") {
        navigate("/student-dashboard");
      } else if (data.response.role === "L") {
        navigate('/lecturer-dashboard');
      } else if (data.response.role === "H") {
        navigate('/hod-dashboard');
      } else if (data.response.role === "B") {
        navigate('/senate-dashboard');
      }
    } else {
      console.log("Login failed", data);
      setLoginError(true);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8" style={{
      "backgroundImage" : `url("https://independent.ng/wp-content/uploads/Drone-view-of-Covenanat-University-Ota-Ogun-State-727x430.jpg")`,
      "backgroundSize": 'cover',
      }}>
      <div className="max-w-md w-full space-y-8 bg-white">
        <div className='flex justify-center'>
          <img src='https://cuportal.covenantuniversity.edu.ng/assets/img/CU_LOGO.jpg' alt= 'CU Logo' height={"100px"} width={"100px"} className='mt-5'></img>
        </div>
        <div>
          <h3 className="mt-6 text-center text-2xl font-bold text-gray-900">Sign in to your account</h3>
        </div>
        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <input type="hidden" name="remember" value="true" />
          <div className="rounded-md shadow-sm -space-y-px">
            <div>
              <label htmlFor="username" className="sr-only">Username</label>
              <input
                id="username"
                name="username"
                type="text"
                required
                className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                placeholder="Username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
              />
            </div>
            <br />
            <div>
              <label htmlFor="password" className="sr-only">Password</label>
              <input
                id="password"
                name="password"
                type="password"
                autoComplete="current-password"
                required
                className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
              {loginError && <p className="text-red-500">Invalid username or password</p>}
            </div>
          </div>
          <div>
            <button
              type="submit"
              className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              disabled={isLoading} // Disable button when loading
            >
              {isLoading ? (
                <svg className="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              ) : (
                <>
                  <span className="absolute left-0 inset-y-0 flex items-center pl-3">
                    <svg className="h-5 w-5 text-indigo-500 group-hover:text-indigo-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                      <path fillRule="evenodd" d="M5 8V6a5 5 0 0 1 10 0v2h2a1 1 0 0 1 1 1v8a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1v-8a1 1 0 0 1 1-1h2zm2-2V6a3 3 0 0 1 6 0v2h-2V6a1 1 0 0 0-2 0v2H7V6a1 1 0 0 0-2 0v2H3V6a3 3 0 0 1 6 0zM4 16h12V9H4v7z" clipRule="evenodd" />
                    </svg>
                  </span>
                  Sign in
                </>
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default Login;
